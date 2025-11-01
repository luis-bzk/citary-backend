package constants

// DatabaseConfig contains database connection pool configuration
var DatabaseConfig = struct {
	MaxOpenConnections int
	MaxIdleConnections int
}{
	MaxOpenConnections: 25,
	MaxIdleConnections: 5,
}
