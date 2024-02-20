package blockchain

import (
	"math/big"

	"github.com/vangxitrum/payment-host/internal/common/aiozcoin"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethermint "github.com/evmos/ethermint/types"
)

const (
	// Bech32Prefix defines the main Bech32 prefix of an account's address
	Bech32Prefix = "aioz"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

const (
	StandardDenom = "aioz"
	MilliDenom    = "milliaioz"
	MicroDenom    = "microaioz"
	NanoDenom     = "nanoaioz"
)

// SetBech32Prefixes sets the global prefixes to be used when serializing addresses and public keys to Bech32 strings.
func SetBech32Prefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(ethermint.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)                      // Shared
	config.SetFullFundraiserPath(ethermint.BIP44HDPath) // nolint: staticcheck
}

var PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(aiozcoin.BaseDenomUnit), nil))

func SetPowerReduction() {
	sdk.DefaultPowerReduction = PowerReduction
}

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if err := sdk.RegisterDenom(StandardDenom, sdk.OneDec()); err != nil {
		panic(err)
	}
	if err := sdk.RegisterDenom(MilliDenom, sdk.NewDecWithPrec(1, 3)); err != nil {
		panic(err)
	}
	if err := sdk.RegisterDenom(MicroDenom, sdk.NewDecWithPrec(1, 6)); err != nil {
		panic(err)
	}
	if err := sdk.RegisterDenom(NanoDenom, sdk.NewDecWithPrec(1, 9)); err != nil {
		panic(err)
	}
	if err := sdk.RegisterDenom(aiozcoin.DefaultDenom, sdk.NewDecWithPrec(1, aiozcoin.BaseDenomUnit)); err != nil {
		panic(err)
	}
}
