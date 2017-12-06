package status

type HealthStatus int

const (
	// indicating that the component or subsystem is in an unknown state.
	UNKNOWN HealthStatus = iota
	//indicating that the component or subsystem is functioning as expected.
	UP
	// indicating that the component or subsystem has suffered an unexpected failure.
	DOWN
	// indicating that the component or subsystem has been taken out of service and should not be used.
	OUT_OF_SERVICE
)