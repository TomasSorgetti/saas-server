package handlers

import (
	"net/http"

	"luthierSaas/internal/application/usecases/auth"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/http/middlewares"

	customErr "luthierSaas/internal/interfaces/http/errors"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
    loginUC   *auth.LoginUseCase
    registerUC   *auth.RegisterUserUseCase
    checkEmailUC *auth.CheckEmailUseCase
	verifyEmailUC *auth.VerifyEmailUseCase
	resendVerificationCodeUC *auth.ResendVerificationCodeUseCase
	refreshTokenUC *auth.RefreshTokenUseCase
	googleLoginUC *auth.LoginGoogleUseCase
	googleCallbackUC *auth.GoogleCallbackUseCase
}


func NewAuthHandler(login *auth.LoginUseCase, register *auth.RegisterUserUseCase, checkEmail *auth.CheckEmailUseCase, verifyEmail *auth.VerifyEmailUseCase, resendVerificationCode *auth.ResendVerificationCodeUseCase, refreshToken *auth.RefreshTokenUseCase, googleLoginUC *auth.LoginGoogleUseCase, googleCallbackUC *auth.GoogleCallbackUseCase) *AuthHandler {
    return &AuthHandler{
		loginUC:          login,
        registerUC:          register,
        checkEmailUC:       checkEmail,
		verifyEmailUC:      verifyEmail,
		resendVerificationCodeUC: resendVerificationCode,
		refreshTokenUC: refreshToken,
		googleLoginUC: googleLoginUC,
		googleCallbackUC: googleCallbackUC,
    }
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dtos.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	
	result, err := h.loginUC.Execute(input)
	if err != nil {
		c.Error(customErr.New(http.StatusBadRequest, "Error to login", err.Error()))
		return
	}
	if result.VerificationRequired {
        c.JSON(http.StatusOK, gin.H{
            "verificationRequired": true,
            "verificationToken":   result.VerificationToken,
			"verificationCodeExpiresAt": result.VerificationExpiresAt,
            "redirect":            result.Redirect,
        })
        return
    }

	// set domain to cookie with secure and httpOnly flags
	// c.SetCookie("access_token", result.AccessToken, 3600, "/", "", true, true) 
	c.SetCookie("access_token", result.AccessToken, 3600, "/", "", false, true) 
	c.SetCookie("refresh_token", result.RefreshToken, 604800, "/", "", false, true) 

	c.JSON(http.StatusOK, result.Profile)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dtos.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	result, err := h.registerUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
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

    exists, err := h.checkEmailUC.Execute(c.Request.Context(), input.Email)
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
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
        return
    }

    _, err := h.verifyEmailUC.Execute(input.VerificationToken, input.VerificationCode)
    if err != nil {
		status := http.StatusBadRequest
        c.Error(customErr.New(status, "error to verify email", err.Error()))
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *AuthHandler) ResendVerificationCode(c *gin.Context) {
	var input dtos.VerifyEmailResendInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(customErr.New(http.StatusBadRequest, "Invalid input data", err.Error()))
		return
	}

	result, err := h.resendVerificationCodeUC.Execute(input.VerificationToken)
	if err != nil {
		c.Error(customErr.New(http.StatusInternalServerError, "Failed to resend verification code", err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userIDVal, exists := c.Get(middlewares.RefreshUserIDKey)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    userID, ok := userIDVal.(int)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
        return
    }

	result, err := h.refreshTokenUC.Execute(userID)

	if err != nil {
		c.Error(customErr.New(http.StatusInternalServerError, "Failed to resend verification code", err.Error()))
		return
	}

	c.SetCookie("access_token", result.AccessToken, 3600, "/", "", false, true) 
	c.SetCookie("refresh_token", result.RefreshToken, 604800, "/", "", false, true) 
	
	c.JSON(http.StatusOK, result.Profile)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Session logout success"})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
    err := h.googleLoginUC.Execute()
	
	c.Error(customErr.New(http.StatusInternalServerError, "not implemented", err.Error()))

}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
    err := h.googleCallbackUC.Execute()
	c.Error(customErr.New(http.StatusInternalServerError, "not implemented", err.Error()))
}