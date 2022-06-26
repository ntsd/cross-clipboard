package config

// Config config struct for cross clipbaord
type Config struct {
	GroupName  string
	ProtocolID string
	ListenHost string
	ListenPort int

	MaxHistory      int  // max number of history clipboard
	HiddenClipboard bool // hidden clipboard text
}
