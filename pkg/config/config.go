package config

// Config config struct for cross clipbaord
type Config struct {
	GroupName  string
	ProtocolID string
	ListenHost string
	ListenPort int

	MaxSize         uint32 // max size to send clipboard
	MaxHistory      int    // max number of history clipboard
	HiddenClipboard bool   // hidden clipboard text
}
