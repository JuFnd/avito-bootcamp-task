package usecase

import (
	"bootcamp-task/pkg/models"
	"bootcamp-task/pkg/variables"
	"bootcamp-task/services/authorization/proto/authorization"
	"bootcamp-task/services/houses/repository"
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IRepository interface {
	CreateHouse(ctx context.Context, address string, yearBuilt int64, developer string) (models.House, error)
	GetHouseFlats(ctx context.Context, houseId int64, userRole string) ([]models.HouseFlat, error)
	CreateFlat(number int64, price int64, rooms int64, houseId int64) (models.HouseFlat, error)
	UpdateFlat(number int64, price int64, rooms int64, houseId int64, status string) (models.HouseFlat, error)
}

type Core struct {
	HousesRepository IRepository
	logger           *slog.Logger
	client           authorization.AuthorizationClient
}

func GetClient(address string) (authorization.AuthorizationClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc connect err: %w", err)
	}
	client := authorization.NewAuthorizationClient(conn)

	return client, nil
}

func GetCore(HousesRelConfig *variables.RelationalDataBaseConfig, grpcCfg *variables.GrpcConfig, logger *slog.Logger) (*Core, error) {
	repository, err := repository.GetHousesRepository(HousesRelConfig, logger)

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Repository can't create %v", err))
	}

	HousesGrpcClient, err := GetClient(grpcCfg.Address + ":" + grpcCfg.Port)
	if err != nil {
		return nil, fmt.Errorf("grpc connect err: %w", err)
	}

	return &Core{
		HousesRepository: repository,
		logger:           logger,
		client:           HousesGrpcClient,
	}, nil
}

func (core *Core) CreateHouse(ctx context.Context, address string, yearBuilt int64, developer string) (models.House, error) {
	house, err := core.HousesRepository.CreateHouse(ctx, address, yearBuilt, developer)
	if err != nil {
		core.logger.Error(variables.HouseCreationError, err)
		return models.House{}, fmt.Errorf(variables.HouseCreationError, err)
	}
	return house, nil
}

func (core *Core) GetHouseFlats(ctx context.Context, houseId int64, userId int64) ([]models.HouseFlat, error) {
	var err error
	var userRole string

	userRole, err = core.GetUserRole(ctx, userId)
	if err != nil {
		userRole = "user"
	}
	flats, err := core.HousesRepository.GetHouseFlats(ctx, houseId, userRole)
	if err != nil {
		core.logger.Error(variables.HouseFlatsError, err)
		return nil, fmt.Errorf(variables.HouseFlatsError, err)
	}
	return flats, nil
}

func (core *Core) CreateFlat(number int64, price int64, rooms int64, houseId int64) (models.HouseFlat, error) {
	flat, err := core.HousesRepository.CreateFlat(number, price, rooms, houseId)
	if err != nil {
		core.logger.Error(variables.FlatCreationError, err)
		return models.HouseFlat{}, fmt.Errorf(variables.FlatCreationError, err)
	}
	return flat, nil
}

func (core *Core) UpdateFlat(number int64, price int64, rooms int64, houseId int64, status string) (models.HouseFlat, error) {
	flat, err := core.HousesRepository.UpdateFlat(number, price, rooms, houseId, status)
	if err != nil {
		core.logger.Error(variables.FlatUpdateError, err)
		return models.HouseFlat{}, fmt.Errorf(variables.FlatUpdateError, err)
	}
	return flat, nil
}

func (core *Core) GetUserRole(ctx context.Context, id int64) (string, error) {
	grpcRequest := authorization.RoleRequest{Id: id}

	grpcResponse, err := core.client.GetRole(ctx, &grpcRequest)
	if err != nil {
		core.logger.Error(variables.GrpcRecievError, err)
		return "", fmt.Errorf(variables.GrpcRecievError, err)
	}
	return grpcResponse.GetRole(), nil
}

func (core *Core) GetUserId(ctx context.Context, sid string) (int64, error) {
	grpcRequest := authorization.FindIdRequest{Sid: sid}

	grpcResponse, err := core.client.GetId(ctx, &grpcRequest)
	if err != nil {
		core.logger.Error(variables.GrpcRecievError, err)
		return 0, fmt.Errorf(variables.GrpcRecievError, err)
	}
	return grpcResponse.Value, nil
}
