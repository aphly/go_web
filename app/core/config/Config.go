package config

type Config struct {
	Cors  *Cors
	Http  *Http
	Db    *map[string]DbGroup
	Log   *Log
	Redis *map[string]RedisGroup
	Oss   *Oss
}

func NewConfig() *Config {
	var config = &Config{}
	config.Cors = CorsConfigLoad()
	config.Http = HttpConfigLoad()
	config.Db = DbConfigLoad()
	config.Log = LogConfigLoad()
	config.Redis = RedisConfigLoad()
	config.Oss = OssConfigLoad()
	return config
}
