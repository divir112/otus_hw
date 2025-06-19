//nolint:depguard
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/divir112/otus_hw/internal/app"                          //nolint:depguard
	"github.com/divir112/otus_hw/internal/config"                       //nolint:depguard
	"github.com/divir112/otus_hw/internal/logger"                       //nolint:depguard
	"github.com/divir112/otus_hw/internal/model"                        //nolint:depguard
	internalhttp "github.com/divir112/otus_hw/internal/server/http"     //nolint:depguard
	memorystorage "github.com/divir112/otus_hw/internal/storage/memory" //nolint:depguard
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
	events := make(map[int]model.Event)
	storage := memorystorage.New(events)
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)
	server.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("HelloWorld!"))
		if err != nil {
			logg.Error(fmt.Sprintf("can't response %v", err))
		}
	})

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

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
