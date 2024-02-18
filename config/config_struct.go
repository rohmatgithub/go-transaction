package config

import (
	"fmt"
	"os"

	"github.com/tkanos/gonfig"
)

var ApplicationConfiguration Configuration

type Configuration interface {
	GetServerConfig() Server
	GetPostgresqlConfig() Postgresql
	GetPostgresqlViewConfig() PostgresqlView
	GetRedisConfig() Redis
	GetLanguageDirectoryPath() string
	GetSqlMigrateDirPath() string
	GetFileResourceDir() DirFileResource
	GetDiscordConfig() Discord
	GetJwtConfig() Jwt
	GetUriResouce() UriResource
}

type Server struct {
	Protocol      string `json:"protocol"`
	Ethernet      string `json:"ethernet"`
	AutoAddHost   bool   `json:"auto_add_host"`
	AutoAddClient bool   `json:"auto_add_client"`
	Host          string `json:"host" envconfig:"$(AUTH_HOST)"`
	Port          int    `json:"port" envconfig:"$(AUTH_PORT)"`
	Version       string `json:"version" envconfig:"$(AUTH_RESOURCE_ID)"`
	ResourceID    string `json:"resource_id"`
	PrefixPath    string `json:"prefix_path"`
	LogLevel      int    `json:"log_level"`
	MaxProcessNCO int    `json:"max_process_nco"`
	IsDevelopment bool   `json:"is_development"`
	IsPraRelease  bool   `json:"is_pra_release"`
}
type Postgresql struct {
	Address           string `json:"address" envconfig:"$(AUTH_DB_CONNECTION)"`
	DefaultSchema     string `json:"default_schema" envconfig:"$(AUTH_DB_PARAM)"`
	MaxOpenConnection int    `json:"max_open_connection"`
	MaxIdleConnection int    `json:"max_idle_connection"`
}
type PostgresqlView struct {
	Address           string `json:"address" envconfig:"$(AUTH_DB_CONNECTION)"`
	DefaultSchema     string `json:"default_schema" envconfig:"$(AUTH_DB_PARAM)"`
	MaxOpenConnection int    `json:"max_open_connection"`
	MaxIdleConnection int    `json:"max_idle_connection"`
}
type Redis struct {
	Host                   string `json:"host" envconfig:"$(AUTH_REDIS_HOST)"`
	Port                   int    `json:"port" envconfig:"$(AUTH_REDIS_PORT)"`
	Db                     int    `json:"db" envconfig:"$(AUTH_REDIS_DB)"`
	Password               string `json:"password" envconfig:"$(AUTH_REDIS_PASSWORD)"`
	Timeout                int    `json:"timeout"`
	RequestVolumeThreshold int    `json:"request_volume_threshold"`
	SleepWindow            int    `json:"sleep_window"`
	ErrorPercentThreshold  int    `json:"error_percent_threshold"`
	MaxConcurrentRequests  int    `json:"max_concurrent_requests"`
}

type Email struct {
	Address  string `json:"address" envconfig:"$(AUTH_SMTP_EMAIL)"`
	Password string `json:"password" envconfig:"$(AUTH_SMTP_PASSWORD)"`
	Port     string `json:"port" envconfig:"$(AUTH_SMTP_PORT)"`
	Host     string `json:"host" envconfig:"$(AUTH_SMTP_HOST)"`
}

type DirFileResource struct {
	Path string `json:"path"`
}

type Discord struct {
	Token     string `json:"token"`
	ChannelID string `json:"channel_id"`
}

type Jwt struct {
	TokenKey string `json:"token_key" envconfig:"$(AUTH_JWT_TOKEN_KEY)"`
}

type UriResource struct {
	MasterData string `json:"master_data"`
}

func GenerateConfiguration(arguments string) {
	var err error
	// enviName := os.Getenv("master-config")
	if arguments == "development" {
		temp := DevelopmentConfig{}
		err = gonfig.GetConf("./config/config_development.json", &temp)
		ApplicationConfiguration = &temp
		//var filename = "config_production.json"
		//err = envconfig.Process(enviName+"/"+filename, &temp)
	} else {
		//temp := ProductionConfig{}
		//var filename = "config_production.json"
		//switch arguments {
		//case "sandbox":
		//	filename = "config_sandbox.json"
		//case "staging":
		//	filename = "config_staging.json"
		//case "dev":
		//	filename = "config_dev.json"
		//}
		//err = gonfig.GetConf(enviName+"/"+filename, &temp)
		//if err != nil {
		//	fmt.Print(err)
		//	os.Exit(2)
		//}
		//err = envconfig.Process(enviName+"/"+filename, &temp)
		//ApplicationConfiguration = &temp
	}
	//ApplicationConfiguration.setFileResourceDir()
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
}
