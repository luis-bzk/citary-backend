package constants

// ErrorMessages contiene mensajes de error estandarizados
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
	NotFound:          "No se encontró el registro solicitado",
	BadRequest:        "Solicitud inválida",
	InternalError:     "Error interno del servidor",
	Unauthorized:      "No autorizado",
	AlreadyExists:     "El recurso ya existe",
	InvalidEmail:      "El email proporcionado no es válido",
	InvalidPassword:   "La contraseña no cumple con los requisitos mínimos",
	UserAlreadyExists: "Ya existe un usuario con ese email",
}

// SuccessMessages contiene mensajes de éxito estandarizados
var SuccessMessages = struct {
	UserCreated string
	UserUpdated string
	UserDeleted string
}{
	UserCreated: "Registro creado exitosamente",
	UserUpdated: "Registro actualizado exitosamente",
	UserDeleted: "Registro eliminado exitosamente",
}
