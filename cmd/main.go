package main

import (
	"fmt"
	api "velocityApi/cmd/api/v1"
	"velocityApi/config"
	"velocityApi/connection"
	"velocityApi/logs"
)

func main() {
	enviroment := "development"
	err := config.ReadConfig(enviroment) // Finalizing the enviroment to run
	if err != nil {
		fmt.Println("An issue occurred, while reading the Enviroment from config file @velocity.go", err.Error())
		return
	}
	logger := logs.New()
	defer logger.Close(0)

	clickhouseConn, clickError := connection.ConnectClickhouse() // Connect to clickhouse
	if clickError != nil {
		logger.Error("An issue occurred, while initiating the connection  to clickhouse db @velocity.go", clickError.Error())
		return
	}
	defer connection.DisconnectClickhouse(clickhouseConn)
	server := api.NewAPIServer(":8080", clickhouseConn, logger)
	if err := server.Run(); err != nil {
		logger.Write("Error while running server. Error: " + err.Error())
	}
}
