package routes

import (
	"net/http"

	"dealer_golang_api/internal/middleware"
	"dealer_golang_api/internal/service/auth"
	"dealer_golang_api/internal/service/brand"
	"dealer_golang_api/internal/service/favorite"
	"dealer_golang_api/internal/service/payment"
	"dealer_golang_api/internal/service/report"
	"dealer_golang_api/internal/service/transaction"
	vtype "dealer_golang_api/internal/service/type"
	"dealer_golang_api/internal/service/user"
	"dealer_golang_api/internal/service/vehicle"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {

	// ========== USER & AUTH ==========
	userRepo := user.NewUserRepository(db)
	userService := user.NewService(userRepo)
	userController := user.NewController(userService)

	authService := auth.NewService(userRepo)
	authController := auth.NewController(authService)

	// health cek
	e.GET("/dealer/health", func(c echo.Context) error {
		status := "healthy"
		dbStatus := "connected"

		if err := db.Ping(c.Request().Context()); err != nil {
			status = "unhealthy"
			dbStatus = "disconnected"
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status":   status,
			"database": dbStatus,
		})
	})

	e.POST("/dealer/register", authController.Register)
	e.POST("/dealer/login", authController.Login)

	// ========== BRAND ==========
	brandRepo := brand.NewRepository(db)
	brandService := brand.NewService(brandRepo)
	brandController := brand.NewController(brandService)

	// ========== TYPE ==========
	typeRepo := vtype.NewRepository(db)
	typeService := vtype.NewService(typeRepo)
	typeController := vtype.NewController(typeService)

	// ========== VEHICLE ==========
	vehicleRepo := vehicle.NewVehicleRepository(db)
	vehicleService := vehicle.NewService(vehicleRepo, brandRepo, typeRepo)
	vehicleController := vehicle.NewController(vehicleService)

	// ========== FAVORITE ==========
	favoriteRepo := favorite.NewRepository(db)
	favoriteService := favorite.NewService(favoriteRepo)
	favoriteController := favorite.NewController(favoriteService)

	// ========== REPORT ==========
	reportService := report.NewService(vehicleService, favoriteService)
	reportController := report.NewController(reportService)

	// ========== TRANSACTION ==========
	transactionRepo := transaction.NewRepository(db)
	paymentService := payment.NewMidtransService()
	transactionService := transaction.NewService(db, transactionRepo, vehicleRepo, paymentService)
	transactionController := transaction.NewController(transactionService, db)

	// ========== PAYMENT CALLBACK ==========
	paymentRepo := payment.NewRepository()
	callbackController := payment.NewCallbackController(db, paymentRepo, paymentService)

	// ========== PROTECTED API ==========
	api := e.Group("/dealer")
	api.Use(middleware.JWTMiddleware())

	// CUSTOMER ROUTES
	customer := api.Group("/customer")
	customer.Use(middleware.RoleMiddleware("customer"))

	customer.GET("/vehicles", vehicleController.GetAll)
	customer.GET("/vehicles/:id", vehicleController.GetByID)
	customer.POST("/favorites", favoriteController.AddFavorite)
	customer.GET("/favorites", favoriteController.GetFavorites)

	// TRANSACTION (CUSTOMER)
	customer.POST("/transactions", transactionController.CreateTransaction)
	customer.GET("/transactions/:order_id", transactionController.GetTransaction)

	// ADMIN ROUTES
	admin := api.Group("/admin")
	admin.Use(middleware.RoleMiddleware("admin"))

	// VEHICLE
	admin.POST("/vehicles", vehicleController.Create)
	admin.PATCH("/vehicles/:id", vehicleController.Update)
	admin.POST("/vehicles/import", vehicleController.Import)
	admin.GET("/vehicles/low-stock", vehicleController.LowStock)

	// BRAND
	admin.POST("/brands", brandController.Create)
	admin.GET("/brands", brandController.GetAll)

	// TYPE
	admin.POST("/types", typeController.Create)
	admin.GET("/types", typeController.GetAll)

	// REPORTS
	admin.GET("/reports/low-stock", reportController.LowStockJSON)
	admin.GET("/reports/low-stock/csv", reportController.LowStockCSV)
	admin.GET("/reports/favorites", reportController.FavoriteJSON)
	admin.GET("/reports/favorites/csv", reportController.FavoriteCSV)

	admin.GET("/users", userController.GetAll)
	admin.GET("/users/:id", userController.GetByID)

	// MIDTRANS CALLBACK (PUBLIC, NO JWT)
	e.POST("/dealer/payments/callback", callbackController.HandleCallback)

}
