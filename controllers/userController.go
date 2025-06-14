package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	database "github.com/qslfrs/Go-API/db"
	"github.com/qslfrs/Go-API/models"
)

type response struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type UserResponse struct {
	ID    int    `json:"id`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GET /users
func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)

	var result []UserResponse
	for _, user := range users {
		result = append(result, UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	c.JSON(http.StatusOK, response{200, result, "success"})
}

// GET /users/:id
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response{404, nil, "user not found"})
		return
	}

	// userResponse := UserResponse{
	// 	ID:    user.ID,
	// 	Name:  user.Name,
	// 	Email: user.Email,
	// }

	data := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}

	c.JSON(http.StatusOK, response{200, data, "success"})
}

// POST /users
func CreateUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response{400, nil, err.Error()})
		return
	}
	user := models.User{Name: input.Name, Email: input.Email}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response{500, nil, err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{200, user, "user created"})
}

// PUT /users/:id
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response{404, nil, "user not found"})
		return
	}
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response{400, nil, err.Error()})
		return
	}
	database.DB.Model(&user).Updates(models.User{Name: input.Name, Email: input.Email})
	c.JSON(http.StatusOK, response{200, user, "user updated"})
}

// DELETE /users/:id
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response{500, nil, err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{200, nil, "user deleted"})
}
