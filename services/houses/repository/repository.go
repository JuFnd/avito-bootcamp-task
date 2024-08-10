package repository

import (
	"bootcamp-task/pkg/models"
	"bootcamp-task/pkg/variables"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

type Repository struct {
	db *sql.DB
}

func GetHousesRepository(configDatabase *variables.RelationalDataBaseConfig, logger *slog.Logger) (*Repository, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password= %s host=%s port=%d sslmode=%s",
		configDatabase.User, configDatabase.DbName, configDatabase.Password, configDatabase.Host, configDatabase.Port, configDatabase.Sslmode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error(variables.SqlOpenError+"%w", "repo", "err", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error(variables.SqlPingError+"%w", "repo", "err", err)
		return nil, err
	}

	db.SetMaxOpenConns(configDatabase.MaxOpenConns)

	profileDb := Repository{
		db: db,
	}

	errs := make(chan error)
	go func() {
		errs <- profileDb.pingDb(configDatabase.Timer, logger)
	}()

	if err := <-errs; err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return &profileDb, nil
}

func (repository *Repository) pingDb(timer uint32, logger *slog.Logger) error {
	var err error
	var retries int

	for retries < variables.MaxRetries {
		err = repository.db.Ping()
		if err == nil {
			return nil
		}

		retries++
		logger.Error(variables.SqlPingError+"%w", "repo", "err", err)
		time.Sleep(time.Duration(timer) * time.Second)
	}

	logger.Error(variables.SqlMaxPingRetriesError, err)
	return fmt.Errorf(fmt.Sprintf(variables.SqlMaxPingRetriesError+" %v", err))
}

func (repository *Repository) CreateHouse(ctx context.Context, address string, yearBuilt int64, developer string) (models.House, error) {
	query := `
	INSERT INTO House (address, year_built, developer)
	VALUES ($1, $2, $3)
	RETURNING house_id, created_at;
`
	var houseID int
	var createdAt time.Time
	err := repository.db.QueryRow(query, address, yearBuilt, developer).Scan(&houseID, &createdAt)
	if err != nil {
		return models.House{}, err
	}

	return models.House{
		HouseID:            houseID,
		Address:            address,
		YearBuilt:          int(yearBuilt),
		Developer:          developer,
		CreatedAt:          createdAt,
		LastApartmentAdded: createdAt,
	}, nil
}

func (repository *Repository) GetHouseFlats(ctx context.Context, houseId int64, userRole string) ([]models.HouseFlat, error) {
	var query string
	if userRole == "user" {
		query = `
            SELECT A.apartment_id, A.apartment_number, A.price, A.rooms, H.house_id, H.address, AH.status
            FROM Apartment A
            JOIN House H ON A.house_id = H.house_id
            JOIN Apartment_House AH ON A.apartment_id = AH.apartment_id
            WHERE H.house_id = $1 AND AH.status = 'approved'
        `
	} else {
		query = `
            SELECT A.apartment_id, A.apartment_number, A.price, A.rooms, H.house_id, H.address, AH.status
            FROM Apartment A
            JOIN House H ON A.house_id = H.house_id
            JOIN Apartment_House AH ON A.apartment_id = AH.apartment_id
            WHERE H.house_id = $1
        `
	}

	rows, err := repository.db.QueryContext(ctx, query, houseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flats []models.HouseFlat
	for rows.Next() {
		var flat models.HouseFlat
		err := rows.Scan(&flat.ApartmentID, &flat.ApartmentNumber, &flat.Price, &flat.Rooms, &flat.HouseID, &flat.Address, &flat.Status)
		if err != nil {
			return nil, err
		}
		flats = append(flats, flat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return flats, nil
}

func (repository *Repository) CreateFlat(number int64, price int64, rooms int64, houseId int64) (models.HouseFlat, error) {
	query := `
    INSERT INTO Apartment (apartment_number, price, rooms, house_id)
    VALUES ($1, $2, $3, $4)
    RETURNING apartment_id;
`
	var apartmentID int
	err := repository.db.QueryRow(query, number, price, rooms, houseId).Scan(&apartmentID)
	if err != nil {
		return models.HouseFlat{}, err
	}
	return models.HouseFlat{
		ApartmentID:     apartmentID,
		ApartmentNumber: int(number),
		Price:           int(price),
		Rooms:           int(rooms),
		HouseID:         int(houseId),
	}, nil
}

func (repository *Repository) UpdateFlat(number int64, price int64, rooms int64, houseId int64, status string) (models.HouseFlat, error) {
	query := `
    UPDATE Apartment
    SET apartment_number = $1, price = $2, rooms = $3
    WHERE house_id = $4 AND apartment_number = $5
    RETURNING apartment_id;
    `
	var apartmentID int
	err := repository.db.QueryRow(query, number, price, rooms, houseId, number).Scan(&apartmentID)
	if err != nil {
		return models.HouseFlat{}, err
	}

	// Update the status in the separate status table
	statusQuery := `
    UPDATE Apartment_House
    SET status = $1
    WHERE apartment_id = $2;
    `
	_, err = repository.db.Exec(statusQuery, status, apartmentID)
	if err != nil {
		return models.HouseFlat{}, err
	}

	return models.HouseFlat{
		ApartmentID:     apartmentID,
		ApartmentNumber: int(number),
		Price:           int(price),
		Rooms:           int(rooms),
		HouseID:         int(houseId),
	}, nil
}
