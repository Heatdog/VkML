package redis

const (
	defaultPort     = 6379
	defaultHost     = "redis"
	defaultDataBase = 0
	defaultPassword = ""
	defaultTTL      = 5
)

type Config struct {
	Host     string
	Password string
	DataBase int
	Port     int
	TTL      int
}

func (conf *Config) WithDefaults() {
	if conf.Port == 0 {
		conf.Port = defaultPort
	}

	if conf.Host == "" {
		conf.Host = defaultHost
	}

	if conf.Password == "" {
		conf.Password = defaultPassword
	}

	if conf.DataBase == 0 {
		conf.DataBase = defaultDataBase
	}

	if conf.TTL == 0 {
		conf.TTL = defaultTTL
	}
}
