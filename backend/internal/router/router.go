package router

import (
	"token-hub/internal/handler"
	"token-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())

	userHandler := handler.NewUserHandler()
	providerHandler := handler.NewProviderHandler()
	paymentHandler := handler.NewPaymentHandler()
	pointsHandler := handler.NewPointsHandler()
	apiHandler := handler.NewAPIHandler()
	callLogHandler := handler.NewCallLogHandler()

	api := r.Group("/api")
	{
		public := api.Group("/v1")
		{
			public.POST("/register", userHandler.Register)
			public.POST("/login", userHandler.Login)

			public.GET("/providers", providerHandler.GetAllProviders)
			public.GET("/payment-methods", paymentHandler.GetPaymentMethods)
			public.GET("/points/rate", pointsHandler.GetPointsRate)
		}

		auth := api.Group("/v1")
		auth.Use(middleware.Auth())
		{
			auth.GET("/user/profile", userHandler.GetProfile)
			auth.PUT("/user/profile", userHandler.UpdateProfile)
			auth.GET("/user/points", userHandler.GetPoints)
			auth.GET("/user/statistics", userHandler.GetUserStatistics)

			auth.GET("/points/records", pointsHandler.GetPointsRecords)

			auth.GET("/recharge/records", paymentHandler.GetUserRechargeRecords)
			auth.POST("/recharge/create", paymentHandler.CreateRechargeOrder)
			auth.POST("/recharge/simulate/:order_no", paymentHandler.SimulatePayment)

			auth.GET("/api-keys", apiHandler.GetAPIKeys)
			auth.POST("/api-keys", apiHandler.CreateAPIKey)
			auth.DELETE("/api-keys/:id", apiHandler.DeleteAPIKey)
			auth.PUT("/api-keys/:id/status", apiHandler.UpdateAPIKeyStatus)

			auth.GET("/call-logs", callLogHandler.GetUserCallLogs)
			auth.GET("/call-logs/statistics", callLogHandler.GetUserStatistics)
		}

		apiKeyAuth := api.Group("/v1")
		apiKeyAuth.Use(middleware.APIKeyAuth())
		{
			apiKeyAuth.GET("/models", apiHandler.ListModels)
			apiKeyAuth.POST("/chat/completions", apiHandler.ChatCompletion)
			apiKeyAuth.POST("/completions", apiHandler.Completion)
		}

		admin := api.Group("/v1/admin")
		admin.Use(middleware.Auth(), middleware.Admin())
		{
			admin.GET("/users", userHandler.GetUserList)
			admin.PUT("/users/:id/status", userHandler.UpdateUserStatus)

			admin.GET("/providers", providerHandler.GetProviderList)
			admin.GET("/providers/:id", providerHandler.GetProviderByID)
			admin.POST("/providers", providerHandler.CreateProvider)
			admin.PUT("/providers/:id", providerHandler.UpdateProvider)
			admin.DELETE("/providers/:id", providerHandler.DeleteProvider)

			admin.GET("/models", providerHandler.GetModelList)
			admin.GET("/models/:id", providerHandler.GetModelByID)
			admin.POST("/models", providerHandler.CreateModel)
			admin.PUT("/models/:id", providerHandler.UpdateModel)
			admin.DELETE("/models/:id", providerHandler.DeleteModel)

			admin.GET("/payment-methods", paymentHandler.GetPaymentMethods)
			admin.POST("/payment-methods", paymentHandler.CreatePaymentMethod)
			admin.PUT("/payment-methods/:id", paymentHandler.UpdatePaymentMethod)
			admin.DELETE("/payment-methods/:id", paymentHandler.DeletePaymentMethod)

			admin.GET("/recharge/records", paymentHandler.GetAllRechargeRecords)

			admin.GET("/points/rate", pointsHandler.GetPointsRate)
			admin.PUT("/points/rate", pointsHandler.SetPointsRate)
			admin.GET("/points/statistics", pointsHandler.GetAllStatistics)

			admin.GET("/call-logs", callLogHandler.GetAllCallLogs)
			admin.GET("/call-logs/daily", callLogHandler.GetDailyStatistics)
			admin.GET("/call-logs/overview", callLogHandler.GetOverviewStatistics)
		}
	}

	return r
}
