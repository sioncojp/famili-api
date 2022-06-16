package config

// Config...Tomlで設定したConfigのstruct
type AppConfig struct {
	Server  ServerConfig    `toml:"server"`
	Service ServiceConfig   `toml:"service"`
	MySQL   DataStoreConfig `toml:"mysql"`
	Log     LogConfig       `toml:"log"`
}

// ServerConfig...serverを立ち上げるために使うもの
type ServerConfig struct {
	// default: famili-apis
	Name string `toml:"name"`

	// default: 8080
	Port string `toml:"port"`
}

// ServiceConfig...Service内で使うもの
type ServiceConfig struct {
	Env string `toml:"env"`
}

// LogConfig...logのstruct
type LogConfig struct {
	// log level: default: info
	Level string `toml:"level"`

	// default: stdout
	OutputPaths []string `toml:"outputPaths"`

	// default: stderr
	ErrorOutputPaths []string `toml:"errorOutputPaths"`
}

// DataStoreConfig...redis/mysqlなどdatastoreのstruct
type DataStoreConfig struct {
	Url      string `toml:"url"`
	Port     string `toml:"port"`
	DbName   string `toml:"dbName"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}
