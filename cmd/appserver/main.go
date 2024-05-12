package main

// Точка входа в приложение, инициализация сервера

import (
	"database/sql"
	"log"
	"myapp/internal/handlers"
	"myapp/internal/repositories"
	"myapp/internal/services"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Инициализация подключения к БД и других ресурсов...

	db, err := sql.Open("postgres", "host=db user=postgres password=15328 dbname=planner sslmode=disable port=5432")
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	folderRepo := repositories.NewFolderRepositories(db)
	folderService := services.NewFolderService(folderRepo)
	folderHandler := handlers.NewFolderHandler(folderService)

	noteRepo := repositories.NewNoteRepositories(db)
	noteService := services.NewNoteService(noteRepo)
	noteHandler := handlers.NewNoteHandlers(noteService)

	financeRepo := repositories.NewFinanceRepositories(db)
	financeService := services.NewFinanceService(financeRepo)
	financeHandler := handlers.NewFinanceHandler(financeService)

	taskRepo := repositories.NewTaskRepositories(db)
	taskServices := services.NewTaskServices(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskServices)

	router := gin.Default()
	router.GET("/", handlers.HomePage)

	// Регистрация маршрута с использованием Gin
	router.POST("/users", userHandler.CreateUser)

	router.GET("/auth", userHandler.AuthenticateUser)

	//---------------------------------------------------------
	router.POST("/folder/create", folderHandler.CreateFolder)

	router.DELETE("/folder/delete", folderHandler.DeleteFolder)

	router.GET("/folder/fetch", folderHandler.FetchFolder)

	router.PUT("/folder/update", folderHandler.UpdateFolder)

	//---------------------------------------------------------
	router.POST("/note/create", noteHandler.CreateNote)

	router.DELETE("/note/delete", noteHandler.DeleteNote)

	router.GET("/note/fetch", noteHandler.FetchNote)

	//---------------------------------------------------------
	router.POST("/finance/create", financeHandler.CreateFinance)

	router.DELETE("/finance/delete", financeHandler.DeleteFinance)

	router.GET("/finance/fetch", financeHandler.FetchFinance)

	router.GET("/finance/fetch-income", financeHandler.FetchFinanceIncome)

	router.GET("/finance/fetch-expense", financeHandler.FetchFinanceExpense)

	//---------------------------------------------------------
	router.POST("/task/create", taskHandler.CreateTask)

	router.DELETE("/task/delete", taskHandler.DeleteTask)

	router.GET("/task/fetch", taskHandler.FetchTask)

	router.PUT("/task/update", taskHandler.UpdateTask)

	router.GET("/task/count", taskHandler.CountTask)

	router.GET("/task/count/favourites", taskHandler.CountTaskFavourites)

	router.GET("/task/fetch/favourites", taskHandler.FetchTaskFavourites)

	// Запуск сервера
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
