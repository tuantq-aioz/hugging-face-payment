package server

import (
	"context"
	"testing"

	proto "github.com/vangxitrum/payment-host/internal/proto/payment_host"
	"google.golang.org/grpc"
)

func TestRegister(t *testing.T) {
	conn, err := grpc.Dial(":8083", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	client := proto.NewPaymentHostServiceClient(conn)

	ctx := context.Background()
	entity, err := client.Register(ctx, &proto.RegisterRequest{
		Name: "test-service",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(entity)
}
