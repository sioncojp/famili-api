package di

import (
	"database/sql"

	"github.com/sioncojp/famili-api/application"
	v1todos "github.com/sioncojp/famili-api/application/v1/todos"
	"github.com/sioncojp/famili-api/infrastructure/database"
	"github.com/sioncojp/famili-api/utils/config"
	"github.com/sioncojp/famili-api/utils/log"
	"github.com/sioncojp/famili-api/utils/mysql"
)

// NewApplication...Applicationを動かすための依存関係を解決する
func NewApplication(configPath string) (*application.HttpHandler, *sql.DB, error) {
	// configの読み込み
	appConfig, err := config.NewConfig(configPath)
	if err != nil {
		return nil, nil, err
	}

	// configの各フィールドのvalidateと初期化
	if err := appConfig.Validate(
		config.ValidateServerConfig,
		config.ValidateServiceConfig,
		config.ValidateMySQLConfig,
		config.ValidateLogConfig,
	); err != nil {
		return nil, nil, err
	}

	// logger初期化
	if err := log.NewLogger(&appConfig.Log); err != nil {
		return nil, nil, err
	}

	// MySQLのhandler初期化
	mysqlHandler, err := mysql.NewMySQLHandler(&appConfig.MySQL)
	if err != nil {
		return nil, nil, err
	}
	db, err := mysqlHandler.DB()

	if err != nil {
		return nil, nil, err
	}

	// service初期化
	s := &application.HttpHandler{}
	s.AppConfig = appConfig
	s.Router.V1.TodosHandler = v1todos.NewHandler(database.NewTodoRepository(mysqlHandler))

	// Router setting
	s.NewRouter()

	return s, db, nil
}
