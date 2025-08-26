package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qslfrs/Go-API/controllers"
)

func Routes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("", controllers.GetUsers)
		users.GET("/:id", controllers.GetUserByID)
		users.POST("", controllers.CreateUser)
		users.PUT("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
	}

	transactions := r.Group("/transaction")
	{
		transactions.GET("", controllers.GetTransaction)
	}
}
