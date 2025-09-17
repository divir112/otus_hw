package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/divir112/otus_hw/internal/config"
	"github.com/divir112/otus_hw/internal/model"
	rabbit "github.com/divir112/otus_hw/internal/rabbit_mq"
	sqlstorage "github.com/divir112/otus_hw/internal/storage/sql"

	"github.com/jackc/pgx/v4/pgxpool"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config_scheduler.yaml", "Path to configuration file")
}

func main() {
	schedulerConfig, err := config.NewConfigScheduler(configFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", schedulerConfig)
	rabbit_client, err := rabbit.NewConnection(schedulerConfig.Rabbit)
	if err != nil {
		panic(fmt.Sprintf("rabbit client: %v", err))
	}

	defer rabbit_client.Close()

	ctx := context.Background()
	connString := fmt.Sprintf("postgresql://%s:%d?dbname=%s&user=%s&password=%s&sslmode=disable", schedulerConfig.Database.Host, schedulerConfig.Database.Port, schedulerConfig.Database.DBName, schedulerConfig.Database.Username, schedulerConfig.Database.Password)
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		panic(err)
	}

	storageSQL := sqlstorage.New(pool)
	err = storageSQL.DeleteOldEvents(ctx, time.Hour*24*365)
	if err != nil {
		panic(err)
	}

	events, err := storageSQL.GetReminderEvents(ctx)
	if err != nil {
		panic(err)
	}

	for _, event := range events {
		notification := model.Notification{
			ID:        event.ID,
			Title:     event.Title,
			StartTime: event.StartTime,
			UserID:    event.UserID,
		}

		msg, err := json.Marshal(notification)
		if err != nil {
			panic(err)
		}

		err = rabbit_client.Publish(ctx, msg)
		if err != nil {
			panic(err)
		}
	}

}
