package config

import (
	"fmt"
	"os"
	"os/user"

	gopenpgp "github.com/ProtonMail/gopenpgp/v2/crypto"
	p2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/ntsd/cross-clipboard/pkg/crypto"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"github.com/ntsd/go-utils/pkg/maputil"
	"github.com/ntsd/go-utils/pkg/stringutil"
	"github.com/spf13/viper"
)

const configDirName = ".cross-clipboard"

// Config config struct for cross clipbaord
type Config struct {
	// Network Config
	GroupName  string `mapstructure:"group_name"`
	ListenHost string `mapstructure:"listen_host"`
	ListenPort int    `mapstructure:"listen_port"`

	// Clipbaord Config
	MaxSize    int `mapstructure:"max_size"`    // limit clipboard size (bytes) to send
	MaxHistory int `mapstructure:"max_history"` // limit number of clipboard history

	// UI Config
	HiddenText bool `mapstructure:"hidden_text"` // hidden clipboard text in UI

	// Device Config
	Username             string            `mapstructure:"-"`           // username of the device
	ID                   p2pcrypto.PrivKey `mapstructure:"-"`           // id private key of this device
	IDPem                string            `mapstructure:"id"`          // id private key pem
	PGPPrivateKey        *gopenpgp.Key     `mapstructure:"-"`           // pgp private key for e2e encryption
	PGPPrivateKeyArmored string            `mapstructure:"private_key"` // armor pgp private key
	AutoTrust            bool              `mapstructure:"auto_trust"`  // auto trust device

	ConfigDirPath string `mapstructure:"config_dir_path"` // config directory path
}

// Save save config to file
func (c *Config) Save() error {
	// set viper value from struct
	m, err := maputil.ToMapString(c, "mapstructure")
	if err != nil {
		return xerror.NewRuntimeError("can not convert config to map").Wrap(err)
	}
	for k, v := range m {
		if k == "-" {
			continue
		}
		viper.Set(k, v)
	}

	err = viper.WriteConfig()
	if err != nil {
		return xerror.NewRuntimeError(fmt.Sprintf(
			"failed to write config at path %s",
			viper.ConfigFileUsed(),
		)).Wrap(err)
	}
	return nil
}

// Save save config to file
func (c *Config) ResetToDefault() error {
	err := os.RemoveAll(c.ConfigDirPath)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig() (*Config, error) {
	thisUser, err := user.Current()
	if err != nil {
		return nil, xerror.NewFatalError("error to get user").Wrap(err)
	}

	dotPath := stringutil.JoinURL(thisUser.HomeDir, configDirName)
	// make directory if not exists
	os.MkdirAll(dotPath, 0777)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dotPath)

	viper.SetDefault("group_name", "default")
	viper.SetDefault("listen_host", "0.0.0.0")
	viper.SetDefault("listen_port", 4001)

	viper.SetDefault("max_size", 5<<20) // 5MB
	viper.SetDefault("max_history", 10)

	viper.SetDefault("hidden_text", true)

	viper.SetDefault("config_dir_path", dotPath)

	idPem, err := crypto.GenerateIDPem()
	if err != nil {
		return nil, xerror.NewFatalError("failed to generate default id pem").Wrap(err)
	}
	viper.SetDefault("id", idPem)
	armoredPrivkey, err := crypto.GeneratePGPKey(thisUser.Username)
	if err != nil {
		return nil, xerror.NewFatalError("failed to generate default pgp key").Wrap(err)
	}
	viper.SetDefault("private_key", armoredPrivkey)
	viper.SetDefault("auto_trust", true)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SafeWriteConfig()
		} else {
			return nil, xerror.NewFatalError("failed to viper.ReadInConfig").Wrap(err)
		}
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, xerror.NewFatalError("failed to viper.Unmarshal").Wrap(err)
	}

	// save config after load default
	err = viper.WriteConfig()
	if err != nil {
		return nil, xerror.NewFatalError("failed to viper.WriteConfig").Wrap(err)
	}

	// set home username
	cfg.Username = thisUser.Username

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
