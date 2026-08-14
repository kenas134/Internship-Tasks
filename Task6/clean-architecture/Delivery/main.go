package delivery

import (
	infrastructure "clean-architecture/Infrastructure"
	repositories "clean-architecture/Repositories"
	usecases "clean-architecture/Usecases"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(
			"No .env file found. Using environment variables.",
		)
	}

	database := os.Getenv("DATABASE_NAME")


	client := infrastructure.ConnectMongoDB()
	db := client.Database(database)


	taskService := repositories.NewTaskRepository(db)
	userService := repositories.NewUserRepository(db)


	//pasword and jwt service

	passwordService := infrastructure.NewPasswordService()
	jwtService :=  infrastructure.NewJWTService()


	taskUsecase := usecases.NewTaskUsecase(taskService)
	userUsecase := usecases.NewUserUsecase(userService,passwordService,jwtService)
	

	



}
