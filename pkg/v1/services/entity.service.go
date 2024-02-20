package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	lens "github.com/strangelove-ventures/lens/client"
	httpClient "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	sdk "github.com/cosmos/cosmos-sdk/types"
	accaddress "github.com/vangxitrum/payment-host/internal/common/accaddress"
	"github.com/vangxitrum/payment-host/internal/common/blockchain"
	"github.com/vangxitrum/payment-host/internal/models"
	internal_services "github.com/vangxitrum/payment-host/internal/services"
	"github.com/vangxitrum/payment-host/pkg/v1/db"
)

type EntityService struct {
	rpcClient *httpClient.HTTP
	ethClient *ethclient.Client

	entityRepo        models.EntityRepository
	paymentMarkRepo   models.PaymentMarkRepository
	walletAddressRepo models.WalletRepository
	txRepo            models.TransactionRepository

	chainId *big.Int

	businessWalletAddr string
	passphrase         string
}

func MustNewEntityService(
	rpcUrl,
	ethUrl,
	passphrase,
	businessAddr string,

	entityRepository models.EntityRepository,
	paymentMarkRepository models.PaymentMarkRepository,
	walletAddressRepository models.WalletRepository,
	txRepo models.TransactionRepository,
) internal_services.EntityService {
	rpcClient, err := lens.NewRPCClient(rpcUrl, time.Second*5)
	if err != nil {
		panic(err)
	}

	ethClient, err := ethclient.Dial(ethUrl)
	if err != nil {
		panic(err)
	}

	chainId, err := ethClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	setupConfig()

	return &EntityService{
		rpcClient: rpcClient,
		ethClient: ethClient,
		chainId:   chainId,

		entityRepo:        entityRepository,
		paymentMarkRepo:   paymentMarkRepository,
		walletAddressRepo: walletAddressRepository,
		txRepo:            txRepo,

		businessWalletAddr: businessAddr,
		passphrase:         passphrase,
	}
}

func setupConfig() {
	sdkCfg := sdk.GetConfig()
	blockchain.SetBech32Prefixes(sdkCfg)
	blockchain.SetBip44CoinType(sdkCfg)
	blockchain.SetPowerReduction()
	sdkCfg.Seal()
}

func (s *EntityService) Register(ctx context.Context, name string) (*models.Entity, error) {
	wallet, err := models.NewWallet(s.passphrase)
	if err != nil {
		return nil, status.Newf(codes.Internal, "failed to create wallet").Err()
	}

	entity := models.NewEntity(name, wallet)
	if err := s.entityRepo.Create(ctx, entity); err != nil {
		return nil, status.Newf(codes.Internal, "failed to create entity").Err()
	}

	return entity, nil
}

func (s *EntityService) Withdraw(
	ctx context.Context,
	entityName string,
	amount decimal.Decimal,
	receiverAddr common.Address,
) (string, error) {
	entity, err := s.entityRepo.GetEntityByName(ctx, entityName)
	if err != nil {
		return "", status.Newf(codes.Internal, "failed to get entity").Err()
	}

	entityPrivateKey, err := crypto.ToECDSA(entity.Wallet.PrivateKey)
	if err != nil {
		return "", status.Newf(codes.Internal, "failed to get private key").Err()
	}

	entityWallet := common.HexToAddress(entity.WalletAddress)
	balance, err := s.ethClient.BalanceAt(ctx, entityWallet, nil)
	if err != nil {
		return "", status.Newf(codes.Internal, "failed to get balance").Err()
	}

	if balance.Cmp(amount.BigInt()) < 0 {
		return "", status.Newf(codes.Internal, "not enough balance").Err()
	}

	nonce, err := s.ethClient.PendingNonceAt(ctx, entityWallet)
	if err != nil {
		return "", status.Newf(codes.Internal, "failed to get nonce").Err()
	}

	gasLimit := uint64(21000)
	gasPrice, err := s.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return "", status.Newf(codes.Internal, "failed to get gas price").Err()
	}

	tx := types.NewTransaction(nonce, receiverAddr, amount.BigInt(), gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(s.chainId), entityPrivateKey)
	if err != nil {
		return "", status.Newf(codes.Internal, "failed to sign transaction").Err()
	}

	if err := s.ethClient.SendTransaction(ctx, signedTx); err != nil {
		return "", status.Newf(codes.Internal, "failed to send transaction").Err()
	}

	return signedTx.Hash().Hex(), nil
}

func (s EntityService) WatchTransaction(ctx context.Context) error {
	chainStatus, err := s.rpcClient.Status(ctx)
	if err != nil {
		return status.Newf(codes.Internal, "failed to get status").Err()
	}

	latestBlock := chainStatus.SyncInfo.LatestBlockHeight
	paymentMark, err := s.paymentMarkRepo.GetPaymentMarkByChainId(ctx, s.chainId.Int64())
	if err != nil && err != gorm.ErrRecordNotFound {
		return status.Newf(codes.Internal, "failed to get payment mark").Err()
	}

	if paymentMark == nil {
		paymentMark = &models.PaymentMark{
			ChainId:     s.chainId.Int64(),
			BlockNumber: latestBlock,
		}

		if err := s.paymentMarkRepo.Create(ctx, paymentMark); err != nil {
			return status.Newf(codes.Internal, "failed to create payment mark").Err()
		}
	}

	fromBlock := paymentMark.BlockNumber
	toBlock := fromBlock + 100
	if toBlock > latestBlock {
		toBlock = latestBlock
	}

	walletAddresses, err := s.walletAddressRepo.GetActiveWallets(ctx)
	if err != nil {
		return status.Newf(codes.Internal, "failed to get active wallets").Err()
	}

	if len(walletAddresses) == 0 {
		if err := s.paymentMarkRepo.UpdatePaymentMarkByChainId(ctx, s.chainId.Int64(), toBlock); err != nil {
			return status.Newf(codes.Internal, "failed to update payment mark").Err()
		}

		return nil
	}

	blockQuery := fmt.Sprintf("tx.height >= %d AND tx.height <= %d", fromBlock, toBlock)
	page := 1
	pageSize := 100
	total := 0
	for {
		resp, err := s.rpcClient.TxSearch(ctx, blockQuery, false, &page, &pageSize, "asc")
		if err != nil {
			return status.Newf(codes.Internal, "failed to get tx search").Err()
		}

		if len(resp.Txs) == 0 {
			break
		}

		for _, tx := range resp.Txs {
			if err := s.handleTransaction(ctx, tx, walletAddresses); err != nil {
				fmt.Println("Handle transaction error: ", err)
			}
		}

		total += len(resp.Txs)
		if total == resp.TotalCount {
			break
		}

	}

	if err := s.paymentMarkRepo.UpdatePaymentMarkByChainId(ctx, s.chainId.Int64(), toBlock); err != nil {
		return status.Newf(codes.Internal, "failed to update payment mark").Err()
	}

	return nil
}

func (s EntityService) handleTransaction(
	ctx context.Context,
	tx *coretypes.ResultTx,
	wallets []*models.Wallet,
) error {
	if tx == nil {
		return nil
	}

	var (
		cosmosTxHash, ethTxHash, senderAddr, receiverAddr string
		contractAddr                                      string
		amount                                            decimal.Decimal
		bigIntAmount                                      *big.Int
		txLog                                             struct {
			Address     string   `json:"address"`
			Topics      []string `json:"topics"`
			Data        []byte   `json:"data"`
			BlockNumber uint64   `json:"blockNumber"`
			LogIndex    int      `json:"logIndex"`
		}
		senderAcc    accaddress.AccAddress
		recipientAcc accaddress.AccAddress
		err          error
		denom        string
		index        int
		blockNumber  uint64
	)

	cosmosTxHash = tx.Hash.String()
	blockNumber = uint64(tx.Height)
	for _, event := range tx.TxResult.Events {
		switch event.Type {
		case "ethereum_tx":
			for _, attr := range event.Attributes {
				if string(attr.Key) == "ethereumTxHash" {
					ethTxHash = string(attr.Value)
				}
			}
		case "tx_log":
			for _, a := range event.Attributes {
				if string(a.Key) == "txLog" {
					if err := json.Unmarshal(a.Value, &txLog); err != nil {
						return status.Newf(codes.Internal, "failed to unmarshal tx log").Err()
					}
				}
			}

			if txLog.Address == "" {
				continue
			}

			contractAddr = txLog.Address
			for i := 1; i < len(txLog.Topics); i++ {
				if i == 1 {
					bytes := common.RightPadBytes(common.FromHex(txLog.Topics[i]), 32)
					senderAddr = common.BytesToAddress(bytes).String()
				} else {
					bytes := common.RightPadBytes(common.FromHex(txLog.Topics[i]), 32)
					receiverAddr = common.BytesToAddress(bytes).String()
					if !isValidAddress(receiverAddr, wallets) && receiverAddr != s.businessWalletAddr {
						receiverAddr = ""
						continue
					}
				}
			}

			bigIntAmount = new(big.Int).SetBytes(common.TrimLeftZeroes(txLog.Data))
			amount, err = decimal.NewFromString(bigIntAmount.String())
			if err != nil {
				return err
			}
		}
	}

	if txLog.Address == "" {
	Loop:
		for i, event := range tx.TxResult.Events {
			if event.Type == bank.EventTypeTransfer {
				for _, attr := range event.Attributes {
					switch string(attr.Key) {
					case "sender":
						senderAcc, err = accaddress.AccAddressFromString(string(attr.Value))
						if err != nil {
							log.Println("AccAddressFromString error ", err)
							continue
						}

						senderAddr = senderAcc.String()
					case "recipient":
						recipientAcc, err = accaddress.AccAddressFromString(string(attr.Value))
						if err != nil {
							log.Println("AccAddressFromString error ", err)
							continue
						}

						receiverAddr = recipientAcc.String()
						if !isValidAddress(receiverAddr, wallets) &&
							receiverAddr != s.businessWalletAddr {
							recipientAcc = nil
							receiverAddr = ""
							continue
						}

						index = i
					case "amount":
						amount, denom, err = models.ParseCoinAmount(string(attr.Value))
						if err != nil {
							log.Println("ParseAmount error ", err)
							continue
						}
					default:
						continue
					}
				}

				if senderAddr != "" && receiverAddr != "" {
					break Loop
				}
			}
		}

	} else {
		if senderAddr == "" || receiverAddr == "" || bigIntAmount == nil {
			return nil
		}

		if !isValidAddress(receiverAddr, wallets) && receiverAddr != s.businessWalletAddr {
			return nil
		}

		blockNumber = txLog.BlockNumber
	}

	entity, err := s.entityRepo.GetEntityByWalletAddress(ctx, receiverAddr)
	if err != nil {
		return status.Newf(codes.Internal, "failed to get entity").Err()
	}

	if contractAddr == "" {
		contractAddr = "aioz"
	}

	transaction := models.Transaction{
		Id:              uuid.New(),
		EntityId:        entity.Id,
		CosmosHash:      cosmosTxHash,
		EvmHash:         ethTxHash,
		ContractAddress: contractAddr,
		From:            senderAddr,
		To:              receiverAddr,
		BlockNumber:     blockNumber,
		Type:            models.CONTRACT_IN_TYPE,
		Index:           index,
		Denom:           denom,
		Amount:          amount,
		Status:          models.TX_STATUS_NEW,
		CreatedAt:       time.Now().UTC().Unix(),
		UpdatedAt:       time.Now().UTC().Unix(),
	}

	txExisted, err := s.txRepo.GetTransactionByHashIndexAndReceiverAddr(
		ctx,
		cosmosTxHash,
		index,
		receiverAddr,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if txExisted == nil {
		if err := s.txRepo.Create(ctx, &transaction); err != nil {
			return err
		}
	}

	return nil
}

func (s *EntityService) NewEntityServiceWithTx(tx *gorm.DB) *EntityService {
	return &EntityService{
		rpcClient: s.rpcClient,
		ethClient: s.ethClient,

		entityRepo:        db.MustNewEntityRepository(tx, false),
		paymentMarkRepo:   db.MustNewPaymentMarkRepository(tx, false),
		walletAddressRepo: db.MustNewWalletRepository(tx, false),
		txRepo:            db.MustNewTransactionRepository(tx, false),

		chainId:            s.chainId,
		businessWalletAddr: s.businessWalletAddr,
		passphrase:         s.passphrase,
	}
}

func isValidAddress(addr string, wallets []*models.Wallet) bool {
	for _, wallet := range wallets {
		if addr == wallet.Address {
			return true
		}
	}

	return false
}
