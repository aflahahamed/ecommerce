package routes

import (
	"github.com/aflahahamed/ecommerce/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine, app *controllers.Application) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	incomingRoutes.GET("/users/productview", controllers.SearchProduct())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
	incomingRoutes.GET("/addtocart", app.AddToCart())
	incomingRoutes.GET("/removeitem", app.RemoveItem())
	incomingRoutes.POST("/addaddress", controllers.AddAddress())
	incomingRoutes.PUT("/edithomeaddress", controllers.EditHomeAddress())
	incomingRoutes.PUT("/editworkaddress", controllers.EditWorkAddress())
	incomingRoutes.GET("/deleteaddresses", controllers.DeleteAddress())
	incomingRoutes.GET("/cartcheckout", app.BuyFromCart())
	incomingRoutes.GET("/instantbuy", app.InstantBuy())
}
