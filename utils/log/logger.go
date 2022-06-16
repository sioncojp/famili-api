package log

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/sioncojp/famili-api/utils/config"
)

// Log...パッケージ全体で使いやすくするために変数でセットする
var Log Logger
var ZapLogger *zap.Logger

// Logger
type Logger interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

// New...http以外で使うLoggerを初期化する
func NewLogger(c *config.LogConfig) error {
	Logger, err := initLog(c)
	if err != nil {
		return errors.Wrap(err, "logger.New")
	}
	defer Logger.Sync()

	// ベースはSugarLoggerを使う。少し低速だが扱いやすいため
	Log = Logger.Sugar()
	ZapLogger = &Logger

	return nil
}

// initLog...zapを使ったloggerを初期化する
func initLog(c *config.LogConfig) (zap.Logger, error) {
	// zap setting
	op := append([]string{"stdout"}, c.OutputPaths...)
	eop := append([]string{"stderr"}, c.ErrorOutputPaths...)

	opAddDoubleQuotes := []string{}
	eopAddDoubleQuotes := []string{}
	for _, v := range op {
		v = strconv.Quote(v)
		opAddDoubleQuotes = append(opAddDoubleQuotes, v)

	}

	for _, v := range eop {
		v = strconv.Quote(v)
		eopAddDoubleQuotes = append(eopAddDoubleQuotes, v)

	}

	rawJSON := []byte(fmt.Sprintf(`{
	 "level": "info",
     "Development": true,
	 "DisableCaller": false,
	 "encoding": "json",
	 "outputPaths": [%s],
	 "errorOutputPaths": [%s],
	 "encoderConfig": {
		"timeKey":        "ts",
		"levelKey":       "level",
		"messageKey":     "msg",
		"nameKey":        "name",
		"stacktraceKey":  "stacktrace",
		"callerKey":      "caller",
		"lineEnding":     "\n",
        "timeEncoder":     "time",
		"levelEncoder":    "lowercaseLevel",
        "durationEncoder": "stringDuration",
		"callerEncoder":   "shortCaller"
	 }
	}`,
		strings.Join(opAddDoubleQuotes, ", "),
		strings.Join(eopAddDoubleQuotes, ", "),
	))

	var cfg zap.Config
	zLogger := &zap.Logger{}

	//standard configuration
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return *zLogger, errors.Wrap(err, "Unmarshal")
	}

	// change log level
	l := zap.NewAtomicLevel().Level()
	err := l.Set(c.Level)
	if err != nil {
		return *zLogger, errors.Wrap(err, "cannot set log level")
	}
	cfg.Level.SetLevel(l)

	// change time jst
	zLogger, err = cfg.Build()
	if err != nil {
		return *zLogger, errors.Wrap(err, "cfg.Build()")
	}

	zLogger.Debug("infralogger construction succeeded")
	return *zLogger, nil
}
