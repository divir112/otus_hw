//nolint:depguard
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/divir112/otus_hw/internal/app"    //nolint:depguard
	"github.com/divir112/otus_hw/internal/config" //nolint:depguard
	"github.com/divir112/otus_hw/internal/logger" //nolint:depguard

	//nolint:depguard
	"github.com/divir112/otus_hw/internal/server/grpc"
	internalhttp "github.com/divir112/otus_hw/internal/server/http" //nolint:depguard

	//nolint:depguard
	sqlstorage "github.com/divir112/otus_hw/internal/storage/sql"
	"github.com/jackc/pgx/v4/pgxpool"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.NewConfig(configFile)
	if err != nil {
		panic(fmt.Sprintf("Can't create config %v", err))
	}

	fmt.Println(config.Logger.Level)
	logg := logger.New(config.Logger.Level, os.Stdout)
	// events := make(map[int]model.Event)
	ctx := context.Background()
	// storage := memorystorage.New(events)
	connstring := fmt.Sprintf("postgresql://postgresql@%s:%d?dbname=%s&user=%s", config.Database.Host, config.Database.Port, config.Database.DBName, config.Database.Username)
	pool, err := pgxpool.Connect(ctx, connstring)
	if err != nil {
		panic(err)
	}
	storageSQL := sqlstorage.New(pool)
	calendar := app.New(logg, storageSQL)

	server := internalhttp.NewServer(logg, calendar, config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		serverGRPC := grpc.NewServer(calendar, logg)
		err := serverGRPC.Start()
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
