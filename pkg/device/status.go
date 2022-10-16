package device

// DeviceStatus device status
type DeviceStatus string

const (
	// StatusPending the device waiting to handshake and trust the device, they send you the public key but you didn't
	StatusPending DeviceStatus = "pending"
	// StatusConnected the device is trusted and connected
	StatusConnected DeviceStatus = "connected"
	// StatusDisconnected the device is trusted but disconnected or offline
	StatusDisconnected DeviceStatus = "disconnected"
	// StatusError found a error in the device should disconnect and reconnect
	StatusError DeviceStatus = "error"
	// StatusBlocked the device is blocked by the user
	StatusBlocked DeviceStatus = "blocked"
)
