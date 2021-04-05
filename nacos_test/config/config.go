package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string    `mapstructure:"name" json:"name"`
}

type ServerConfig struct {
	Port        int           `mapstructure:"port" json:"port"`
	Name        string        `mapstructure:"name" json:"name"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ConsulConfig struct {
	Host        string           `mapstructure:"host" json:"host"`
	Port        int           `mapstructure:"port" json:"port"`
}
