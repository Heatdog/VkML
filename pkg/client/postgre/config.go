package postgre

const (
	defaultPort        = 5432
	defaultHost        = "postgre"
	defaultTimePrepare = 5
	defaultTimeWait    = 3
)

type Config struct {
	Host     string `mapstructure:"host"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`

	TimePrepare int `mapstructure:"time_prepare"`
	TimeWait    int `mapstructure:"time_wait"`
}

func (conf *Config) WithDefaults() {
	if conf.Port == 0 {
		conf.Port = defaultPort
	}

	if conf.Host == "" {
		conf.Host = defaultHost
	}

	if conf.TimePrepare == 0 {
		conf.TimePrepare = defaultTimePrepare
	}

	if conf.TimeWait == 0 {
		conf.TimeWait = defaultTimeWait
	}
}
