package handlers

import (
	"luthierSaas/internal/application/usecases/user"
	"luthierSaas/internal/interfaces/http/dtos"
	customErr "luthierSaas/internal/interfaces/http/errors"
	"luthierSaas/internal/interfaces/http/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	profileUC   *user.ProfileUseCase
    changePasswordUC *user.ChangePasswordUseCase
}

func NewUserHandler(profile *user.ProfileUseCase, changePassword *user.ChangePasswordUseCase) *UserHandler {
    return &UserHandler{
		profileUC:          profile,
        changePasswordUC:    changePassword,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
    userIDVal, exists := c.Get(middlewares.UserIDKey)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    userID, ok := userIDVal.(int)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
        return
    }
    
    result, err := h.profileUC.Execute(userID)
    if err != nil {
        c.Error(customErr.New(http.StatusBadRequest, "Error to get profile", err.Error()))
        return
    }

    c.JSON(http.StatusOK, result)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
    userIDVal, exists := c.Get(middlewares.UserIDKey)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    userID, ok := userIDVal.(int)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
        return
    }

    var input dtos.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(customErr.New(http.StatusBadRequest, "Invalid input data", err.Error()))
		return
	}

    err := h.changePasswordUC.Execute(userID, input)
    if err != nil {
        c.Error(customErr.New(http.StatusBadRequest, "Error to get profile", err.Error()))
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}