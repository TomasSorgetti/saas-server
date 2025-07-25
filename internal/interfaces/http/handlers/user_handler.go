package handlers

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
    return &UserHandler{}
}

func (h *UserHandler) Create(c *gin.Context) {
	c.JSON(http.StatusOK, "User created")
}

func (h *UserHandler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, "All users")
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid id")
		return
	}
	c.JSON(http.StatusOK, "User by id" + strconv.Itoa(id))
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid id")
		return
	}
	c.JSON(http.StatusOK, "User updated" + strconv.Itoa(id))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid id")
		return
	}
	c.JSON(http.StatusOK, "User deleted" + strconv.Itoa(id))
}