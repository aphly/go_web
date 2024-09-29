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
	config.Cors = CorsConfigLoad()
	config.Http = HttpConfigLoad()
	config.Db = DbConfigLoad()
	config.Log = LogConfigLoad()
	config.Redis = RedisConfigLoad()
	return config
}
