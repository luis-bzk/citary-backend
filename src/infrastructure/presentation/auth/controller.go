package authController

import (
	dtosAuth "citary-backend/src/domain/dtos/auth"
	useCaseAuth "citary-backend/src/domain/use_cases/auth"
	"citary-backend/src/infrastructure/infraHttp"
	"citary-backend/src/shared/constants"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	signupUserUseCase *useCaseAuth.SignupUserUseCase
}

func NewAuthController(signupUserUseCase *useCaseAuth.SignupUserUseCase) *AuthController {
	return &AuthController{signupUserUseCase: signupUserUseCase}
}

func (a *AuthController) SignupUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		infraHttp.SendError(w, constants.StatusCode.BadRequest, "Método no permitido")
		return
	}

	var dto dtosAuth.SignupUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		infraHttp.SendError(w, constants.StatusCode.BadRequest, "JSON inválido")
		return
	}
	if err := dto.Validate(); err != nil {
		infraHttp.SendError(w, constants.StatusCode.BadRequest, err.Error())
		return
	}

	user, err := a.signupUserUseCase.Execute(dto)
	if err != nil {
		infraHttp.HandleDomainError(w, err)
		return
	}

	authResponse := map[string]interface{}{
		"id":            user.ID,
		"email":         user.Email,
		"emailVerified": user.EmailVerified,
		"createdDate":   user.CreatedDate,
	}

	infraHttp.SendSuccess(w, constants.StatusCode.Created, constants.SuccessMessages.UserCreated, authResponse)
}
