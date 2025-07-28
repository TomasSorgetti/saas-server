package handlers

import (
	"luthierSaas/internal/application/usecases/user"
	customErr "luthierSaas/internal/interfaces/http/errors"
	"luthierSaas/internal/interfaces/http/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	profileUC   *user.ProfileUseCase
}

func NewUserHandler(profile *user.ProfileUseCase) *UserHandler {
    return &UserHandler{
		profileUC:          profile,
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