package delivery

import (
	"bootcamp-task/pkg/middleware"
	"bootcamp-task/pkg/models"
	"bootcamp-task/pkg/util"
	"bootcamp-task/pkg/variables"
	"context"
	"log/slog"
	"net/http"
)

type ICore interface {
	CreateHouse(ctx context.Context, address string, yearBuilt int64, developer string) (models.House, error)
	GetHouseFlats(ctx context.Context, houseId int64) ([]models.HouseFlat, error)
	CreateFlat(number int64, price int64, rooms int64, houseId int64) (models.HouseFlat, error)
	UpdateFlat(number int64, price int64, rooms int64, houseId int64) (models.HouseFlat, error)
	GetUserRole(ctx context.Context, id int64) (string, error)
	GetUserId(ctx context.Context, sid string) (int64, error)
}

type API struct {
	logger *slog.Logger
	mux    *http.ServeMux
	core   ICore
}

func GetHousesApi(core ICore, logger *slog.Logger) *API {
	api := &API{
		logger: logger,
		mux:    http.NewServeMux(),
		core:   core,
	}

	api.mux.HandleFunc("/api/v1/house/create", api.createHouse)
	api.mux.HandleFunc("/api/v1/flat/update", api.updateFlat)
	createHouseHandler := middleware.PermissionsMiddleware(api.mux, api.core, variables.ModeratorRole, api.logger)
	createHouseHandler = middleware.AuthorizationMiddleware(createHouseHandler, api.core, api.logger)
	createHouseHandler = middleware.MethodMiddleware(createHouseHandler, variables.MethodPost, api.logger)

	houseMux := http.NewServeMux()
	houseMux.Handle("/", createHouseHandler)
	houseMux.HandleFunc("/api/v1/house/", api.getHouse)

	flatMux := http.NewServeMux()
	flatMux.Handle("/", houseMux)
	flatMux.HandleFunc("/api/v1/flat/create", api.createFlat)

	flatHandler := middleware.AuthorizationMiddleware(flatMux, api.core, api.logger)
	flatHandler = middleware.MethodMiddleware(flatHandler, variables.MethodPost, api.logger)

	panicMux := http.NewServeMux()
	panicMux.Handle("/", flatHandler)
	panicHandler := middleware.PanicMiddleware(panicMux, api.logger)

	siteMux := http.NewServeMux()
	siteMux.Handle("/", panicHandler)
	api.mux = siteMux

	return api
}

func (api *API) ListenAndServe(appConfig *variables.AppConfig) error {
	err := http.ListenAndServe(appConfig.Address, api.mux)
	if err != nil {
		api.logger.Error(variables.ListenAndServeError, "%w", err.Error())
		return err
	}
	return nil
}

func (api *API) createHouse(w http.ResponseWriter, r *http.Request) {
	var houseRequest models.House

	err := util.GetRequestBody(w, r, &houseRequest, api.logger)
	if err != nil {
		return
	}

	year64 := int64(houseRequest.YearBuilt)
	house, err := api.core.CreateHouse(r.Context(), houseRequest.Address, year64, houseRequest.Developer)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.StatusInternalServerError, err, api.logger)
		return
	}

	util.SendResponse(w, r, http.StatusOK, house, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) getHouse(w http.ResponseWriter, r *http.Request) {
	houseId, err := util.GetHouseIdFromRequest(r.URL.Path)
	if err != nil {
		util.SendResponse(w, r, http.StatusBadRequest, nil, variables.StatusBadRequestError, err, api.logger)
		return
	}

	flats, err := api.core.GetHouseFlats(r.Context(), houseId)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.StatusInternalServerError, err, api.logger)
		return
	}

	util.SendResponse(w, r, http.StatusOK, flats, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) createFlat(w http.ResponseWriter, r *http.Request) {
	var flatRequest models.HouseFlat

	err := util.GetRequestBody(w, r, &flatRequest, api.logger)
	if err != nil {
		return
	}

	houseId, err := util.GetHouseIdFromRequest(r.URL.Path)
	if err != nil {
		util.SendResponse(w, r, http.StatusBadRequest, nil, variables.StatusBadRequestError, err, api.logger)
		return
	}

	flat, err := api.core.CreateFlat(int64(flatRequest.ApartmentNumber), int64(flatRequest.Price), int64(flatRequest.Rooms), houseId)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.StatusInternalServerError, err, api.logger)
		return
	}

	util.SendResponse(w, r, http.StatusOK, flat, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) updateFlat(w http.ResponseWriter, r *http.Request) {
	var flatRequest models.HouseFlat

	err := util.GetRequestBody(w, r, &flatRequest, api.logger)
	if err != nil {
		return
	}

	houseId, err := util.GetHouseIdFromRequest(r.URL.Path)
	if err != nil {
		util.SendResponse(w, r, http.StatusBadRequest, nil, variables.StatusBadRequestError, err, api.logger)
		return
	}

	flat, err := api.core.UpdateFlat(int64(flatRequest.ApartmentNumber), int64(flatRequest.Price), int64(flatRequest.Rooms), houseId)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.StatusInternalServerError, err, api.logger)
		return
	}

	util.SendResponse(w, r, http.StatusOK, flat, variables.StatusOkMessage, nil, api.logger)
}
