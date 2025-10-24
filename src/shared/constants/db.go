package constants

var DBConstants = struct {
	MaxOpenConnections int
	MaxIdleConnections int
}{
	MaxOpenConnections: 25,
	MaxIdleConnections: 5,
}
