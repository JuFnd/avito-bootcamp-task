package main

import (
	"bootcamp-task/configs"
	"bootcamp-task/pkg/variables"
	delivery "bootcamp-task/services/houses/delivery/http"
	"bootcamp-task/services/houses/usecase"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	logFile, err := os.Create("houses.log")
	if err != nil {
		fmt.Println("Error creating log file")
		return
	}

	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	housesAppConfig, err := configs.ReadHousesAppConfig()
	if err != nil {
		logger.Error(variables.ReadAuthConfigError, err.Error())
		return
	}

	relationalDataBaseConfig, err := configs.ReadRelationalHousesDataBaseConfig()
	if err != nil {
		logger.Error(variables.ReadAuthSqlConfigError, err.Error())
		return
	}

	grpcCfg, err := configs.ReadGrpcConfig()
	if err != nil {
		logger.Error("failed to parse grpc configs file: %s", err.Error())
		return
	}

	core, err := usecase.GetCore(relationalDataBaseConfig, grpcCfg, logger)
	if err != nil {
		logger.Error(variables.CoreInitializeError, err)
		return
	}

	api := delivery.GetHousesApi(core, logger)

	errApi := api.ListenAndServe(housesAppConfig)
	if errApi != nil {
		logger.Error(variables.ListenAndServeError, "%w", errApi.Error())
		return
	}
}
