package constants

// ErrorMessages contains standardized error messages
var ErrorMessages = struct {
	NotFound          string
	BadRequest        string
	InternalError     string
	Unauthorized      string
	AlreadyExists     string
	InvalidEmail      string
	InvalidPassword   string
	UserAlreadyExists string
}{
	NotFound:          "The requested record was not found",
	BadRequest:        "Invalid request",
	InternalError:     "Internal server error",
	Unauthorized:      "Unauthorized",
	AlreadyExists:     "The resource already exists",
	InvalidEmail:      "The provided email is not valid",
	InvalidPassword:   "The password does not meet minimum requirements",
	UserAlreadyExists: "A user with that email already exists",
}

// SuccessMessages contains standardized success messages
var SuccessMessages = struct {
	UserCreated string
	UserUpdated string
	UserDeleted string
}{
	UserCreated: "User created successfully",
	UserUpdated: "User updated successfully",
	UserDeleted: "User deleted successfully",
}
