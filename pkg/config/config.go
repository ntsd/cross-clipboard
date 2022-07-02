package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config config struct for cross clipbaord
type Config struct {
	GroupName  string `mapstructure:"group_name"`
	ProtocolID string `mapstructure:"protocal_id"`
	ListenHost string `mapstructure:"listen_host"`
	ListenPort int    `mapstructure:"listen_port"`

	MaxSize    uint32 `mapstructure:"max_size"`    // max size to send clipboard
	MaxHistory int    `mapstructure:"max_history"` // max number of history clipboard

	TerminalMode bool `mapstructure:"terminal_mode"` // is terminal mode or ui mode
	HiddenText   bool `mapstructure:"hidden_text"`   // hidden clipboard text in UI

	PrivateKey string `mapstructure:"private_key"` // private key for libp2p and p2p encryption
}

func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("group_name", "default")
	viper.SetDefault("protocal_id", "/cross-clipboard/0.0.1")
	viper.SetDefault("listen_host", "0.0.0.0")
	viper.SetDefault("listen_port", 4001)

	viper.SetDefault("max_size", 1<<24) // 16MB
	viper.SetDefault("max_history", 10)

	viper.SetDefault("terminal_mode", false)
	viper.SetDefault("hidden_text", false)

	viper.SetDefault("private_key", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SafeWriteConfig()
		} else {
			panic(err)
		}
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	log.Println("loaded config:", cfg)

	return cfg
}
