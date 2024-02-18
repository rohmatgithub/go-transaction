package config

type DevelopmentConfig struct {
	Configuration
	Server                Server          `json:"server"`
	Postgresql            Postgresql      `json:"postgresql"`
	PostgresqlView        PostgresqlView  `json:"postgresql_view"`
	Redis                 Redis           `json:"redis"`
	LanguageDirectoryPath string          `json:"language_directory_path"`
	SqlMigrateDirPath     string          `json:"sql_migrate_dir_path"`
	DirFileResource       DirFileResource `json:"dir_file_resource"`
	Discord               Discord         `json:"discord"`
	Jwt                   Jwt             `json:"jwt"`
	UriResource           UriResource     `json:"uri_resource"`
}

func (input DevelopmentConfig) GetServerConfig() Server {
	return input.Server
}

func (input DevelopmentConfig) GetPostgresqlConfig() Postgresql {
	return input.Postgresql
}

func (input DevelopmentConfig) GetPostgresqlViewConfig() PostgresqlView {
	return input.PostgresqlView
}

func (input DevelopmentConfig) GetRedisConfig() Redis {
	return input.Redis
}

func (input DevelopmentConfig) GetLanguageDirectoryPath() string {
	return input.LanguageDirectoryPath
}

func (input DevelopmentConfig) GetSqlMigrateDirPath() string {
	return input.SqlMigrateDirPath
}

func (input DevelopmentConfig) GetFileResourceDir() DirFileResource {
	return input.DirFileResource
}

func (input DevelopmentConfig) GetDiscordConfig() Discord {
	return input.Discord
}

func (input DevelopmentConfig) GetJwtConfig() Jwt {
	//return Jwt{
	//	CodeKey:          "",
	//	TokenKey:         "",
	//	TokenInternalKey: "",
	//}
	return input.Jwt
}
func (input DevelopmentConfig) GetUriResouce() UriResource {
	return input.UriResource
}
