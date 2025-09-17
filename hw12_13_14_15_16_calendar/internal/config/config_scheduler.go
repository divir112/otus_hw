package config

import (
	"fmt"

	"github.com/divir112/otus_hw/internal/apperror"
	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigScheduler struct {
	Logger   LoggerConf
	Database DBConf
	Rabbit   RabbitScheduler
}

type RabbitScheduler struct {
	Host  string "yaml:host env-required"  //nolint
	Port  string "yaml:port env-required"  //nolint
	Queue string "yaml:queue env-required" //nolint
}

func (r RabbitScheduler) GetHost() string {
	return r.Host
}
func (r RabbitScheduler) GetPort() string {
	return r.Port
}
func (r RabbitScheduler) GetQueue() string {
	return r.Queue
}

func NewConfigScheduler(path string) (*ConfigScheduler, error) {
	cfg := &ConfigScheduler{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("[config::NewConfigScheduler]: %w, %w", apperror.ErrorParseConfig, err)
	}
	return cfg, nil
}
