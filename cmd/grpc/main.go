package main

import (
	"fmt"

	"github.com/vangxitrum/payment-host/internal/server"
)

func main() {
	fmt.Printf("Payment host server is running on port %s\n", appConfig.ServerPort)

	cron.Start()
	server.MustMakeGrpcPaymentHostServerAndRun(fmt.Sprintf(":%s", appConfig.ServerPort), entityService)
}
