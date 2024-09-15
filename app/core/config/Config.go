package config

type Config struct {
	Cors  *Cors
	Http  *Http
	Db    *map[string]DbGroup
	Log   *Log
	Redis *map[string]RedisGroup
}

func NewConfig() *Config {
	var config = &Config{}
	config.Cors = CorsConfig()
	config.Http = HttpConfig()
	config.Db = DbConfig()
	config.Log = LogConfigs()
	config.Redis = RedisConfig()
	return config
}
