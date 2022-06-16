package main

import (
	"flag"
	"os"

	"go.uber.org/zap"

	"github.com/sioncojp/famili-api/di"
)

func main() {
	file := flag.String("c", "", "toml file")
	flag.Parse()
	os.Exit(Run(*file))
}

// Run...初期化して、appを動かす
func Run(configFilePath string) int {
	// TODO: datadog apm

	// DIコンテナ初期化用のlogger作成
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	suggar := logger.Sugar()

	app, db, err := di.NewApplication(configFilePath)
	if err != nil {
		suggar.Errorf("%+v", err)
		return 1
	}
	defer db.Close()
	app.RunServer()

	return 0
}
