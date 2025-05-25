package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qslfrs/Go-API/controllers"
)

func UserRoutes(r *gin.Engine) {
	grp := r.Group("/users")
	{
		grp.GET("", controllers.GetUsers)
		grp.GET("/:id", controllers.GetUserByID)
		grp.POST("", controllers.CreateUser)
		grp.PUT("/:id", controllers.UpdateUser)
		grp.DELETE("/:id", controllers.DeleteUser)
	}
}
