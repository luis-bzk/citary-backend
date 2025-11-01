package constants

// StatusCode contains common HTTP status codes
var StatusCode = struct {
	Ok                  int
	Created             int
	BadRequest          int
	Unauthorized        int
	Forbidden           int
	NotFound            int
	Conflict            int
	InternalServerError int
}{
	Ok:                  200,
	Created:             201,
	BadRequest:          400,
	Unauthorized:        401,
	Forbidden:           403,
	NotFound:            404,
	Conflict:            409,
	InternalServerError: 500,
}
