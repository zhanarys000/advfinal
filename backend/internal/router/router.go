package router

import (
	"adv/internal/handlers"
	"adv/pkg/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Rourer struct {
	user handlers.UserHandlers
	book handlers.BookHandlers
}

func NewRourer(user handlers.UserHandlers, book handlers.BookHandlers) *Rourer {
	return &Rourer{user: user, book: book}
}

func (r *Rourer) Setup(app *gin.Engine) {
	api := app.Group("/api")
	{
		limiter := rate.NewLimiter(1, 3)
		api.Use(middleware.RateLimiterMiddleware(limiter))
		api.POST("/register", r.user.Register)
		api.POST("/login", r.user.Login)
		api.POST("/getverificatoincode", r.user.GetVerificationCode)
		api.POST("checkverificatiocode", r.user.CheckVerificationCode)

		api.GET("/profile", middleware.RequireAuthMiddleware, r.user.Profile)
		api.PUT("/update", middleware.RequireAuthMiddleware, r.user.UpdateProfile)
		api.POST("/subscribe", middleware.RequireAuthMiddleware, r.user.Subscribe)
		api.GET("/books", r.book.GetAllBooks)
		api.POST("/send", middleware.RequireAuthMiddleware, middleware.AdminMiddleware, r.user.SendEmail)
		book := api.Group("/books", middleware.RequireAuthMiddleware)
		{
			book.GET("/all", r.book.GetByParams)
			book.GET(`/:id`, r.book.GetBook)
			book.POST("/buy/:id", r.book.BuyBook)
			book.GET("/my", r.book.GetUserBooks)
			book.POST("/", middleware.AdminMiddleware, r.book.CreateBook)
			book.PUT("/:id", middleware.AdminMiddleware, r.book.UpdateBook)
			book.DELETE("/:id", middleware.AdminMiddleware, r.book.DeleteBook)

		}
	}
}
