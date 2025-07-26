package handlers

import (
	"net/http"

	"luthierSaas/internal/application/usecases/auth"
	"luthierSaas/internal/interfaces/http/dtos"

	customErr "luthierSaas/internal/interfaces/http/errors"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
    registerUC   *auth.RegisterUserUseCase
    checkEmailUC *auth.CheckEmailUseCase
	verifyEmailUC *auth.VerifyEmailUseCase
}


func NewAuthHandler(register *auth.RegisterUserUseCase, checkEmail *auth.CheckEmailUseCase, verifyEmail *auth.VerifyEmailUseCase) *AuthHandler {
    return &AuthHandler{
        registerUC:   register,
        checkEmailUC: checkEmail,
		verifyEmailUC: verifyEmail,
    }
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.Error(customErr.New(http.StatusBadRequest, "Invalid input data", "Expected JSON body with email and password"))
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dtos.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	result, err := h.registerUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el usuario"})
		return
	}

	c.JSON(http.StatusCreated, result)
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

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var input dtos.VerifyEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(customErr.New(http.StatusBadRequest, "Invalid input data", err.Error()))
		return
	}

	success, err := h.verifyEmailUC.Execute(input.VerificationToken, input.VerificationCode)
	if err != nil {
		c.Error(customErr.New(http.StatusInternalServerError, "Failed to verify email", err.Error()))
		return
	}

	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Verification failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}