package handlers

// обработчики HTTP-запросов
// такие как GET, POST, PUT, DELETE

import (
	"myapp/internal/models"
	"myapp/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// структура для инъекции зависимостей сервисов.
type UserHandler struct {
	UserService *services.UserService
}

// создает экземпляр UserHandler.
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызов метода CreateUser из UserService для создания пользователя
	if err := h.UserService.CreateUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *UserHandler) AuthenticateUser(c *gin.Context) {
	var loginDetails struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, token, ok, err := h.UserService.Authenticate(loginDetails.Login, loginDetails.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid login details"})
		return
	}

	// Отправка ответа клиенту, например, токена для доступа
	c.JSON(http.StatusOK, gin.H{"userID": userID, "token": token, "message": "Authentication successful"})
}
