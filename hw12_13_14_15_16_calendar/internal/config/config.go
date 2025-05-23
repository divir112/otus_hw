package config

import (
	"fmt"

	"github.com/divir112/otus_hw/internal/apperror"
	"github.com/ilyakaznacheev/cleanenv"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf
	Database DBConf
}

type LoggerConf struct {
	Level string "yaml:level env-required" //nolint
}

type DBConf struct {
	Type     string "yaml:level env-required"  //nolint
	Host     string "yaml:host env-required"   //nolint
	Port     int    "yaml:port env-required"   //nolint
	Username string "yaml:username"            //nolint
	Password string "yaml:password"            //nolint
	DBName   string "yaml:dbname env-required" //nolint
}

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("[NewConfig]: %w, %w", apperror.ErrorParseConfig, err)
	}
	return cfg, nil
}
