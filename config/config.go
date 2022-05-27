package config

import "github.com/spf13/viper"

type Config struct {
	PORT         string `mapstructure:"PORT"`
	CQL_KEYSPACE string `mapstructure:"CQL_KEYSPACE"`
	CQL_HOSTS    string `mapstructure:"CQL_HOSTS"`
	NATS_CLUSTER string `mapstructure:"NATS_CLUSTER"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
