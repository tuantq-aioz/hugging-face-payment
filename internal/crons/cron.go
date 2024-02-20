package crons

import (
	"context"

	"github.com/robfig/cron"

	"github.com/vangxitrum/payment-host/internal/services"
)

type Cron struct {
	cron    *cron.Cron
	service services.EntityService
}

func NewCron(service services.EntityService) *Cron {
	return &Cron{
		cron:    cron.New(),
		service: service,
	}
}

func (c *Cron) Start() {
	c.cron.AddFunc("@every 5s", func() {
		c.service.WatchTransaction(context.Background())
	})

	c.cron.Start()
}
