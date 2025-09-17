package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/divir112/otus_hw/internal/config"
	rabbit "github.com/divir112/otus_hw/internal/rabbit_mq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config_sender.yaml", "Path to configuration file")
}

func main() {
	senderConfig, err := config.NewConfigScheduler(configFile)
	if err != nil {
		panic(err)
	}

	rabbit_client, err := rabbit.NewConnection(senderConfig.Rabbit)
	if err != nil {
		panic(fmt.Sprintf("rabbit client: %v", err))
	}

	ctx := context.Background()
	err = rabbit_client.Consume(ctx, func(msg []byte) error {
		fmt.Println(string(msg))
		return nil
	})

	if err != nil {
		panic(err)
	}
}
