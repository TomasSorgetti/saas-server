package handlers

import (
	"luthierSaas/internal/application/usecases/user"
	customErr "luthierSaas/internal/interfaces/http/errors"
	"net/http"
	"strconv"

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
	userIDStr := c.Param("user_id")

	userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

	result, err := h.profileUC.Execute(userID)

	if err != nil {
		c.Error(customErr.New(http.StatusBadRequest, "Invalid input data", err.Error()))
		return
	}

	c.JSON(http.StatusFound, result)
}