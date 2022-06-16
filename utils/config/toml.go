package config

import (
	"os"

	"github.com/pkg/errors"
	toml "github.com/sioncojp/tomlssm"
)

const (
	ServerName = "famili-apis"
	ServerPort = "8080"
	MySQLPort  = "3306"
	LogLevel   = "info"
)

type ValidateFunc func(*AppConfig) error

// NewConfig...Appを立ち上げるのに必要なロジックを初期化する。tomlファイルを読む。またParameterStoreからDecodeして取得する
func NewConfig(configFilePath string) (*AppConfig, error) {
	var config AppConfig
	if _, err := toml.DecodeFile(configFilePath, &config, os.Getenv("AWS_DEFAULT_REGION")); err != nil {
		return nil, errors.Wrap(err, "unmarshal config file")
	}
	return &config, nil
}

// Validate...ConfigのValidateを行う
func (c *AppConfig) Validate(validateFuncs ...ValidateFunc) error {
	var errorCollector []error

	// 引数のエラー関数を実行し、エラーを収集
	for _, fn := range validateFuncs {
		err := fn(c)
		if err != nil {
			errorCollector = append(errorCollector, err)
		}
	}

	// エラーが1つでもあれば返す
	if len(errorCollector) > 0 {
		var result error
		for _, v := range errorCollector {
			result = errors.Wrap(result, v.Error())
		}

		return result
	}

	return nil
}

// ValidateServerConfig...Server Structのvalidate
var ValidateServerConfig ValidateFunc = func(c *AppConfig) error {
	v := c.Server
	if v.Name == "" {
		c.Server.Name = ServerName
	}

	if v.Port == "" {
		c.Server.Port = ServerPort
	}
	return nil
}

// ValidateServiceConfig...Service Structのvalidate
var ValidateServiceConfig ValidateFunc = func(c *AppConfig) error {
	v := c.Service
	if v.Env == "" {
		return errors.New("env is not set in validateService")
	}
	return nil
}

// ValidateMySQLConfig...MySQL Structのvalidate
var ValidateMySQLConfig ValidateFunc = func(c *AppConfig) error {
	v := c.MySQL
	if v.Port == "" {
		c.MySQL.Port = MySQLPort
	}

	return nil
}

// ValidateLogConfig...Log Structのvalidate
var ValidateLogConfig ValidateFunc = func(c *AppConfig) error {
	v := c.Log
	if v.Level == "" {
		c.Log.Level = LogLevel
	}
	return nil
}
