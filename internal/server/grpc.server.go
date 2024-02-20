package server

import (
	"context"
	"log"
	"net"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	proto "github.com/vangxitrum/payment-host/internal/proto/payment_host"
	"github.com/vangxitrum/payment-host/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func MustMakeGrpcPaymentHostServerAndRun(
	listenAddr string,
	entityService services.EntityService,
) {
	grpcSourceControlServer := newPaymentHostServer(entityService)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	option := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(option...)
	healthServer := health.NewServer()
	healthServer.SetServingStatus(proto.PaymentHostService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	proto.RegisterPaymentHostServiceServer(grpcServer, grpcSourceControlServer)
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

type PaymentHostServer struct {
	entityService services.EntityService
	proto.UnimplementedPaymentHostServiceServer
}

func newPaymentHostServer(entityService services.EntityService) *PaymentHostServer {
	return &PaymentHostServer{
		entityService: entityService,
	}
}

func (s *PaymentHostServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if req.Name == "" {
		return nil, status.Newf(codes.InvalidArgument, "name is required").Err()
	}

	entity, err := s.entityService.Register(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &proto.RegisterResponse{
		WalletAddress: entity.WalletAddress,
	}, nil
}

func (s *PaymentHostServer) Deposit(ctx context.Context, req *proto.WithdrawRequest) (*proto.WithdrawResponse, error) {
	if req.ReceiverWalletAddress == "" {
		return nil, status.Newf(codes.InvalidArgument, "wallet address is required").Err()
	}

	if req.Amount <= 0 {
		return nil, status.Newf(codes.InvalidArgument, "amount must be greater than 0").Err()
	}

	amount := decimal.NewFromInt(req.Amount)
	receiverAddr := common.HexToAddress(req.ReceiverWalletAddress)
	txHash, err := s.entityService.Withdraw(ctx, req.EntityName, amount, receiverAddr)
	if err != nil {
		return nil, err
	}

	return &proto.WithdrawResponse{
		TransactionHash: txHash,
	}, nil
}
