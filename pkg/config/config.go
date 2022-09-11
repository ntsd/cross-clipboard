package config

import (
	"encoding/base64"
	"log"

	"github.com/ntsd/cross-clipboard/pkg/crypto"
	"github.com/spf13/viper"
)

const defaultPassphrase = "PleaseChangeThis!@#$"

// Config config struct for cross clipbaord
type Config struct {
	// Network Config
	GroupName  string `mapstructure:"group_name"`
	ProtocolID string `mapstructure:"protocal_id"`
	ListenHost string `mapstructure:"listen_host"`
	ListenPort int    `mapstructure:"listen_port"`

	// Clipbaord Config
	MaxSize    uint32 `mapstructure:"max_size"`    // max size to send clipboard
	MaxHistory int    `mapstructure:"max_history"` // max number of history clipboard

	// UI Config
	TerminalMode bool `mapstructure:"terminal_mode"` // is terminal mode or ui mode
	HiddenText   bool `mapstructure:"hidden_text"`   // hidden clipboard text in UI

	// Device Config
	ID            string `mapstructure:"id"`          // id of this client
	GPGPrivateKey string `mapstructure:"private_key"` // private key for libp2p and p2p encryption
	AutoTrust     bool   `mapstructure:"auto_trust"`  // auto trust device
	Passphrase    string `mapstructure:"passphrase"`  // passphrase in base64 encoded, will use to encrypt public key
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

	viper.SetDefault("id", getDefaultID())
	viper.SetDefault("private_key", "")
	viper.SetDefault("passphrase", base64.StdEncoding.EncodeToString([]byte(defaultPassphrase)))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SafeWriteConfig()
		} else {
			log.Fatal(err)
		}
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("loaded config:", cfg)

	// save config after load default
	viper.WriteConfig()

	return cfg
}

func getDefaultID() string {
	prvKey, _, err := crypto.NewKeyPair()
	if err != nil {
		log.Fatal(err)
	}
	prvKeyBytes, err := crypto.MarshalPrivateKey(prvKey)
	if err != nil {
		log.Fatal(err)
	}
	return string(prvKeyBytes)
}

func getDefaultPrivateKey() string {
	return ""
}
