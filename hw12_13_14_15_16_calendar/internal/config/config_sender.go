package config

import (
	"fmt"

	"github.com/divir112/otus_hw/internal/apperror"
	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigSender struct {
	Logger   LoggerConf
	Database DBConf
	Rabbit   RabbitSender
}

type RabbitSender struct {
	Host  string "yaml:host env-required"  //nolint
	Port  string "yaml:port env-required"  //nolint
	Queue string "yaml:queue env-required" //nolint
}

func (r RabbitSender) GetHost() string {
	return r.Host
}
func (r RabbitSender) GetPort() string {
	return r.Port
}
func (r RabbitSender) GetQueue() string {
	return r.Queue
}

func NewConfigSender(path string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("[config::NewConfigSender]: %w, %w", apperror.ErrorParseConfig, err)
	}
	return cfg, nil
}
