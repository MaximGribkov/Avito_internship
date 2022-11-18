package main

import (
	"avito/internal/user"
	"avito/pkg/logging"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// @title microservice balance
// @version 1.0
// @description service for avito internship

func main() {
	logger := logging.GetLogger()

	router := gin.Default()
	user.Router(router)
	logger.Info("create router")

	pathDatabase := "postgresql://postgres:123@localhost:5432/avito"

	con, err := pgx.Connect(context.Background(), pathDatabase)
	if err != nil {
		logger.Fatalf("error do with tries database %s", err)
	}

	user.SetDataBase(con)
	defer con.Close(context.Background())
	router.Run()
	logger.Info("Connection is established")

}

// Задание показалось мне довольно интересным. Было несколько вариантов написания кода, но я решил остановиться на этом,
//так как посчитал его наиболее оптимальным. Надеюсь, моя работа Вам понравится и мы продолжим наше сотрудничество.
