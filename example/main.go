package main

import (
	"context"
	"github.com/IT-Kungfu/logger"
	"github.com/IT-Kungfu/mt5client"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	log    *logger.Logger
	mt5cfg *mt5client.Config
	mt5    *mt5client.Pool
)

func init() {
	var err error
	log, err = logger.New(&logger.Config{
		LogLevel:     "debug",
		SentryDSN:    "",
		LogstashAddr: "",
		ServiceName:  "logger",
		InstanceName: "dev",
	})
	if err != nil {
		panic(err)
	}

	mt5cfg = &mt5client.Config{
		MT5Host:           "mt5hostname",
		MT5Port:           443,
		MT5Login:          "login",
		MT5Password:       "password",
		MT5APIVersion:     "2190",
		MT5APIAgent:       "mt5client",
		MT5PingTimeout:    20,
		MT5RequestTimeout: 10,
	}
}

func main() {
	services := map[string]interface{}{
		"mt5cfg": mt5cfg,
		"log":    log,
	}
	ctx := context.WithValue(context.Background(), "services", services)

	var err error
	mt5, err = mt5client.NewMT5ClientPool(ctx, 10)
	if err != nil {
		log.Fatalf("MT5 client error: %v", err)
	}

	mt5.GetDealsBatch("2170933", []string{}, []string{}, 0, time.Now().Unix())

	go response()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	mt5.Close()
}

func response() {
	defer close(mt5.Response())
	for m := range mt5.Response() {
		if m.Err != nil {
			log.Errorf("#%d MT5 client error: %v", m.ClientId, m.Err)
			continue
		}

		deals := m.Response.(*mt5client.DealsResponse).Deals
		for _, d := range deals {
			log.Debugf("%+v", d)
		}
	}
}
