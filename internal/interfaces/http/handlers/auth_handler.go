package handlers

import (
	"errors"
	"net/http"

	"luthierSaas/internal/application/usecases/auth"
	"luthierSaas/internal/infrastructure/persistance/repositories"
	"luthierSaas/internal/interfaces/http/dtos"

	customErr "luthierSaas/internal/interfaces/http/errors"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
    registerUC   *auth.RegisterUserUseCase
    checkEmailUC *auth.CheckEmailUseCase
}


func NewAuthHandler(register *auth.RegisterUserUseCase, checkEmail *auth.CheckEmailUseCase) *AuthHandler {
    return &AuthHandler{
        registerUC:   register,
        checkEmailUC: checkEmail,
    }
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.Error(customErr.New(http.StatusBadRequest, "Invalid input data", "Expected JSON body with email and password"))
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dtos.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(customErr.New(http.StatusBadRequest, "Invalid input data", err.Error()))
		return
	}

	err := h.registerUC.Execute(input)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrEmailAlreadyExists):
			c.Error(customErr.New(http.StatusConflict, "Email already registered"))
		default:
			c.Error(customErr.New(http.StatusInternalServerError, "Failed to register user", err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) CheckEmail(c *gin.Context) {
	var input dtos.CheckEmailInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(customErr.New(http.StatusBadRequest, "Invalid input data", err.Error()))
		return
	}

	exists, err := h.checkEmailUC.Execute(input.Email)
	if err != nil {
		c.Error(customErr.New(http.StatusInternalServerError, "Failed to check email", err.Error()))
		return
	}

	if exists {
		c.JSON(http.StatusOK, gin.H{"exists": true, "message": "Email already registered"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"exists": false, "message": "Email available"})
}