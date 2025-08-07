package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"luthierSaas/internal/application/usecases/auth"
	"luthierSaas/internal/interfaces/http/dtos"

	customErr "luthierSaas/internal/interfaces/http/errors"

	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
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
	logoutUC *auth.LogoutUseCase
}


func NewAuthHandler(
	login *auth.LoginUseCase, 
	register *auth.RegisterUserUseCase, 
	checkEmail *auth.CheckEmailUseCase, 
	verifyEmail *auth.VerifyEmailUseCase, 
	resendVerificationCode *auth.ResendVerificationCodeUseCase, 
	refreshToken *auth.RefreshTokenUseCase, 
	googleLoginUC *auth.LoginGoogleUseCase, 
	googleCallbackUC *auth.GoogleCallbackUseCase, 
	logoutUC *auth.LogoutUseCase,
	) *AuthHandler {
    
	return &AuthHandler{
		loginUC:          login,
        registerUC:          register,
        checkEmailUC:       checkEmail,
		verifyEmailUC:      verifyEmail,
		resendVerificationCodeUC: resendVerificationCode,
		refreshTokenUC: refreshToken,
		googleLoginUC: googleLoginUC,
		googleCallbackUC: googleCallbackUC,
		logoutUC: logoutUC,
    }
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dtos.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	ua := useragent.New(c.GetHeader("User-Agent"))
    browser, version := ua.Browser()
    deviceType := "Desktop"
    if ua.Mobile() {
        deviceType = "Mobile"
    }
    deviceInfo := fmt.Sprintf("%s %s, %s, %s", browser, version, ua.OS(), deviceType)
	
	result, err := h.loginUC.Execute(input, deviceInfo)
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
	refreshToken, err := c.Cookie("refresh_token")
    if err != nil || refreshToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no refresh token provided"})
        return
    }

    ua := useragent.New(c.GetHeader("User-Agent"))
    browser, version := ua.Browser()
    deviceType := "Desktop"
    if ua.Mobile() {
        deviceType = "Mobile"
    }
    deviceInfo := fmt.Sprintf("%s %s, %s, %s", browser, version, ua.OS(), deviceType)

    ctx := c.Request.Context()
    result, err := h.refreshTokenUC.Execute(ctx, refreshToken, deviceInfo)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

	c.SetCookie("access_token", result.AccessToken, 3600, "/", "", false, true) 
	c.SetCookie("refresh_token", result.RefreshToken, 604800, "/", "", false, true) 
	
	c.JSON(http.StatusOK, result.Profile)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
    if err != nil || accessToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no access token provided"})
        return
    }

    ctx := c.Request.Context()
    err = h.logoutUC.Execute(ctx, accessToken)
    if err != nil {
        if errors.Is(err, errors.New("session not found")) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "session not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Session logout success"})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	googleUrl, err := h.googleLoginUC.Execute(c.Request.Context())
	if err != nil {
        redirectURL, _ := url.Parse("http://localhost:5173/auth/login")
        query := redirectURL.Query()
        query.Set("error", "failed_to_generate_google_url")
        query.Set("message", "No se pudo generar la URL de inicio de sesión con Google")
        redirectURL.RawQuery = query.Encode()
        c.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
        return
    }

	c.Redirect(http.StatusTemporaryRedirect, googleUrl)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	redirectErrorURL, err := url.Parse("http://localhost:5173/auth/login")
	if err != nil {
        log.Printf("Failed to parse error redirect URL: %v", err)
        c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173/auth/login?error=internal_error&message=Error%20interno%20del%20servidor")
        return
    }

	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		c.Error(customErr.New(http.StatusBadRequest, "Missing code or state", ""))
		return
	}

	ua := useragent.New(c.GetHeader("User-Agent"))
	browser, version := ua.Browser()
	deviceType := "Desktop"
	if ua.Mobile() {
		deviceType = "Mobile"
	}
	deviceInfo := fmt.Sprintf("%s %s, %s, %s", browser, version, ua.OS(), deviceType)

	result, err := h.googleCallbackUC.Execute(c.Request.Context(), code, state, deviceInfo)
	if err != nil {
        query := redirectErrorURL.Query()
        query.Set("error", "google_callback_failed")
        query.Set("message", "Error al procesar la autenticación con Google")
        redirectErrorURL.RawQuery = query.Encode()
        c.Redirect(http.StatusTemporaryRedirect, redirectErrorURL.String())
        return
    }

	redirectURL, err := url.Parse("http://localhost:5173/auth/google/callback")
	if err != nil {
        query := redirectErrorURL.Query()
        query.Set("error", "invalid_redirect_url")
        query.Set("message", "Error al procesar la URL de redirección")
        redirectErrorURL.RawQuery = query.Encode()
        c.Redirect(http.StatusTemporaryRedirect, redirectErrorURL.String())
        return
    }

	query := redirectURL.Query()
	if result.VerificationRequired {
		query.Set("verificationRequired", "true")
		query.Set("verificationToken", result.VerificationToken)
		query.Set("verificationCodeExpiresAt", result.VerificationExpiresAt.Format("2006-01-02T15:04:05Z07:00"))
	} else {
		c.SetCookie("access_token", result.AccessToken, 3600, "/", "", false, true)
		c.SetCookie("refresh_token", result.RefreshToken, 604800, "/", "", false, true)
		profileJSON, err := json.Marshal(result.Profile)
		if err != nil {
            query := redirectErrorURL.Query()
            query.Set("error", "profile_marshal_failed")
            query.Set("message", "Error al procesar el perfil del usuario")
            redirectErrorURL.RawQuery = query.Encode()
            c.Redirect(http.StatusTemporaryRedirect, redirectErrorURL.String())
            return
        }
		query.Set("profile", string(profileJSON))
	}
	redirectURL.RawQuery = query.Encode()

	c.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
}