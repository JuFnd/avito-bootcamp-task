package configs

import (
	"bootcamp-task/pkg/variables"
	"flag"
	"fmt"
	"os"
	"syscall"

	"gopkg.in/yaml.v2"
)

func readYAMLFile[T any](filePath string) (*T, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var config T
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML data: %w", err)
	}

	return &config, nil
}

func ParseFlagsAndReadYAMLFile[T any](fileName string, defaultFilePath string, flags *flag.FlagSet) (*T, error) {
	flag.Parse()
	var path string
	flag.StringVar(&path, fileName, defaultFilePath, "Путь к конфигу"+fileName)

	config, err := readYAMLFile[T](path)
	if err == syscall.ENOENT {
		return nil, fmt.Errorf(fmt.Sprintf("Failed to parse %s from provided path: %v", fileName, err))
	}

	return config, nil
}

func ReadAuthAppConfig() (*variables.AppConfig, error) {
	return ParseFlagsAndReadYAMLFile[variables.AppConfig]("auth_config_path", "configs/AuthorizationAppConfig.yml", flag.CommandLine)
}

func ReadGrpcConfig() (*variables.GrpcConfig, error) {
	return ParseFlagsAndReadYAMLFile[variables.GrpcConfig]("grpc_config_path", "configs/GrpcConfig.yml", flag.CommandLine)
}

func ReadHousesAppConfig() (*variables.AppConfig, error) {
	return ParseFlagsAndReadYAMLFile[variables.AppConfig]("houses_config_path", "configs/HousesAppConfig.yml", flag.CommandLine)
}

func ReadRelationalAuthDataBaseConfig() (*variables.RelationalDataBaseConfig, error) {
	return ParseFlagsAndReadYAMLFile[variables.RelationalDataBaseConfig]("sql_config_auth_path", "configs/AuthorizationSqlDataBaseConfig.yml", flag.CommandLine)
}

func ReadRelationalHousesDataBaseConfig() (*variables.RelationalDataBaseConfig, error) {
	return ParseFlagsAndReadYAMLFile[variables.RelationalDataBaseConfig]("sql_config_houses_path", "configs/HousesSqlDataBaseConfig.yml", flag.CommandLine)
}

func ReadCacheDatabaseConfig() (*variables.CacheDataBaseConfig, error) {
	return ParseFlagsAndReadYAMLFile[variables.CacheDataBaseConfig]("cache_config_path", "configs/AuthorizationCacheDataBaseConfig.yml", flag.CommandLine)
}
