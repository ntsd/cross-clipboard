package config

import (
	"log"
	"os/user"

	gopenpgp "github.com/ProtonMail/gopenpgp/v2/crypto"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/ntsd/cross-clipboard/pkg/crypto"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"github.com/spf13/viper"
)

// Config config struct for cross clipbaord
type Config struct {
	// Network Config
	GroupName  string `mapstructure:"group_name"`
	ProtocolID string `mapstructure:"protocal_id"`
	ListenHost string `mapstructure:"listen_host"`
	ListenPort int    `mapstructure:"listen_port"`

	// Clipbaord Config
	EncryptEnabled bool `mapstructure:"encrypt_enabled"` // encryption clipbaord enabled
	MaxSize        int  `mapstructure:"max_size"`        // max size to send clipboard
	MaxHistory     int  `mapstructure:"max_history"`     // max number of history clipboard

	// UI Config
	TerminalMode bool `mapstructure:"terminal_mode"` // is terminal mode or ui mode
	HiddenText   bool `mapstructure:"hidden_text"`   // hidden clipboard text in UI

	// Device Config
	Username             string            `mapstructure:"-"`           // username of the device
	ID                   p2pcrypto.PrivKey `mapstructure:"-"`           // id private key of this device
	IDPem                string            `mapstructure:"id"`          // id private key pem
	PGPPrivateKey        *gopenpgp.Key     `mapstructure:"-"`           // pgp private key for e2e encryption
	PGPPrivateKeyArmored string            `mapstructure:"private_key"` // armor pgp private key
	AutoTrust            bool              `mapstructure:"auto_trust"`  // auto trust device
}

// Save save config to file
func (c Config) Save() error {
	err := viper.WriteConfig()
	if err != nil {
		return xerror.NewRuntimeError("failed to viper.WriteConfig").Wrap(err)
	}
	return nil
}

func LoadConfig() (Config, error) {
	user, err := user.Current()
	if err != nil {
		return Config{}, xerror.NewFatalError("error to get user").Wrap(err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("group_name", "default")
	viper.SetDefault("protocal_id", "/cross-clipboard/0.0.1")
	viper.SetDefault("listen_host", "0.0.0.0")
	viper.SetDefault("listen_port", 4001)

	viper.SetDefault("encrypt_enabled", true)
	viper.SetDefault("max_size", 1<<24) // 16MB
	viper.SetDefault("max_history", 10)

	viper.SetDefault("terminal_mode", false)
	viper.SetDefault("hidden_text", false)

	idPem, err := crypto.GenerateIDPem()
	if err != nil {
		return Config{}, xerror.NewFatalError("failed to generate default id pem").Wrap(err)
	}
	viper.SetDefault("id", idPem)
	armoredPrivkey, err := crypto.GeneratePGPKey(user.Username)
	if err != nil {
		return Config{}, xerror.NewFatalError("failed to generate default pgp key").Wrap(err)
	}
	viper.SetDefault("private_key", armoredPrivkey)
	viper.SetDefault("auto_trust", true)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SafeWriteConfig()
		} else {
			return Config{}, xerror.NewFatalError("failed to viper.ReadInConfig").Wrap(err)
		}
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return Config{}, xerror.NewFatalError("failed to viper.Unmarshal").Wrap(err)
	}
	log.Println("loaded config:", cfg)

	// save config after load default
	err = viper.WriteConfig()
	if err != nil {
		return Config{}, xerror.NewFatalError("failed to viper.WriteConfig").Wrap(err)
	}

	// set home username
	cfg.Username = user.Username

	// unmarshal id
	idPK, err := crypto.UnmarshalIDPrivateKey(cfg.IDPem)
	if err != nil {
		return cfg, xerror.NewFatalError("failed to unmarshal id private key").Wrap(err)
	}
	cfg.ID = idPK

	// unmarshal pgp private key
	pgpPrivateKey, err := crypto.UnmarshalPGPKey(cfg.PGPPrivateKeyArmored, nil)
	if err != nil {
		return cfg, xerror.NewFatalError("failed to unmarshal gpg private key").Wrap(err)
	}
	cfg.PGPPrivateKey = pgpPrivateKey

	return cfg, nil
}
