package npp

// TransportError is an error that is returned when the underlying gRPC
// transport is broken.
type TransportError struct {
	error
}

func (m TransportError) Error() string {
	return m.error.Error()
}
