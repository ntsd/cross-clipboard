package device

// DeviceStatus device status
type DeviceStatus int

const (
	// StatusUnknown unknown device status
	StatusUnknown DeviceStatus = iota
	// StatusPending the device waiting to handshake and trust the device, they send you the public key but you didn't
	StatusPending
	// StatusConnected the device is trusted and connected
	StatusConnected
	// StatusDisconnected the device is trusted but disconnected or offline
	StatusDisconnected
	// StatusError found a error in the device should disconnect and reconnect
	StatusError
	// StatusBlocked the device is blocked by the user
	StatusBlocked
)

func (ds DeviceStatus) ToString() string {
	switch ds {
	case StatusUnknown:
		return "unknown"
	case StatusPending:
		return "pending"
	case StatusConnected:
		return "connected"
	case StatusDisconnected:
		return "disconnected"
	case StatusError:
		return "error"
	case StatusBlocked:
		return "blocked"
	default:
		return "unknown"
	}
}
