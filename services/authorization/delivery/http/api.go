package delivery

import (
	"bootcamp-task/pkg/middleware"
	"bootcamp-task/pkg/models"
	communication "bootcamp-task/pkg/requests"
	"bootcamp-task/pkg/util"
	"bootcamp-task/pkg/variables"
	"bootcamp-task/services/authorization/usecase"
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Core interface
type ICore interface {
	KillSession(ctx context.Context, sid string) error
	FindActiveSession(ctx context.Context, sid string) (bool, error)
	CreateSession(ctx context.Context, login string) (models.Session, error)
	CreateUserAccount(login string, password string) error
	FindUserByLogin(login string) (bool, error)
	FindUserAccount(login string, password string) (*models.UserItem, bool, error)
	GetUserId(ctx context.Context, sid string) (int64, error)
	GetUserRole(ctx context.Context, id int64) (string, error)
}

type API struct {
	core   ICore
	logger *slog.Logger
	mux    *http.ServeMux
}

func (api *API) ListenAndServe(appConfig *variables.AppConfig) error {
	err := http.ListenAndServe(appConfig.Address, api.mux)
	if err != nil {
		api.logger.Error(variables.ListenAndServeError, "%w", err.Error())
		return err
	}
	return nil
}

func GetAuthorizationApi(authCore *usecase.Core, authLogger *slog.Logger) *API {
	api := &API{
		core:   authCore,
		logger: authLogger,
		mux:    http.NewServeMux(),
	}

	api.mux.HandleFunc("/login", api.Signin)
	api.mux.HandleFunc("/register", api.Signup)

	signHandler := middleware.MethodMiddleware(api.mux, variables.MethodPost, api.logger)

	logoutMux := http.NewServeMux()
	logoutMux.Handle("/", signHandler)
	logoutMux.HandleFunc("/logout", api.LogoutSession)

	//authHandler := middleware.AuthorizationMiddleware(logoutMux, api.core, api.logger)
	authHandler := middleware.MethodMiddleware(logoutMux, variables.MethodPost, api.logger)
	authHandler = middleware.PanicMiddleware(authHandler, api.logger)

	siteMux := http.NewServeMux()
	siteMux.Handle("/", authHandler)
	api.mux = siteMux

	return api
}

// @Summary SignIn
// @Tags authentication
// @Description Authenticate user by providing login and password credentials
// @ID authenticate-user
// @Accept json
// @Produce json
// @Param input body communication.SigninRequest true "login and password"
// @Success 200 {string} string "Authentication token"
// @Failure 401 {string} string variables.StatusUnauthorizedError
// @Failure 500 {string} string variables.StatusInternalServerError
// @Router /login [post]
func (api *API) Signin(w http.ResponseWriter, r *http.Request) {
	var signinRequest communication.SigninRequest

	err := util.GetRequestBody(w, r, &signinRequest, api.logger)
	if err != nil {
		return
	}

	user, found, err := api.core.FindUserAccount(signinRequest.Login, signinRequest.Password)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.StatusInternalServerError, err, api.logger)
		return
	}

	if !found {
		util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.StatusUnauthorizedError, nil, api.logger)
		return
	}

	session, err := api.core.CreateSession(r.Context(), user.Login)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.SessionCreateError, err, api.logger)
		return
	}

	authorizationCookie := util.GetCookie(variables.SessionCookieName, session.SID, "/", variables.HttpOnly, session.ExpiresAt)
	http.SetCookie(w, authorizationCookie)
	util.SendResponse(w, r, http.StatusOK, nil, variables.StatusOkMessage, nil, api.logger)
}

// @Summary SignUp
// @Tags registration
// @Desription Create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body communication.SignupRequest true "account information"
// @Success 200 {integer} object communication.SignupResponse
// @Failure 400 {string} string variables.InvalidLoginOrPasswordError
// @Failure 401 {string} string variables.InvalidLoginOrPasswordError
// @Failure 500 {string} string variables.StatusInternalServerError
// @Router /register [post]
func (api *API) Signup(w http.ResponseWriter, r *http.Request) {
	var signupRequest communication.SignupRequest

	err := util.GetRequestBody(w, r, &signupRequest, api.logger)
	if err != nil {
		return
	}

	found, err := api.core.FindUserByLogin(signupRequest.Login)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.StatusInternalServerError, err, api.logger)
		return
	}

	if found {
		util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.UserAlreadyExistsError, nil, api.logger)
		return
	}

	err = api.core.CreateUserAccount(signupRequest.Login, signupRequest.Password)
	if err != nil && err.Error() == variables.InvalidLoginOrPasswordError {
		util.SendResponse(w, r, http.StatusBadRequest, variables.InvalidLoginOrPasswordError, variables.InvalidLoginOrPasswordError, err, api.logger)
		return
	}
	if err != nil {
		util.SendResponse(w, r, http.StatusUnauthorized, variables.InvalidLoginOrPasswordError, variables.UserAlreadyExistsError, err, api.logger)
		return
	}

	response := communication.SignupResponse{Login: signupRequest.Login}
	util.SendResponse(w, r, http.StatusOK, response, variables.StatusOkMessage, nil, api.logger)
}

// @Summary Logout
// @Tags authentication
// @Description End current user's active session
// @ID end-current-session
// @Accept json
// @Produce json
// @Header 200 {integer} 1
// @Success 200 {string} string "Session ended successfully."
// @Failure 401 {string} string variables.StatusUnauthorizedError
// @Failure 500 {string} string variables.StatusInternalServerError
// @Router /logout [post]
func (api *API) LogoutSession(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err != nil {
		util.SendResponse(w, r, http.StatusUnauthorized, variables.StatusUnauthorizedError, variables.SessionNotFoundError, nil, api.logger)
		return
	}

	found, err := api.core.FindActiveSession(r.Context(), session.Value)
	if err != nil {
		util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.StatusUnauthorizedError, err, api.logger)
		return
	}

	if !found {
		util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.SessionNotFoundError, nil, api.logger)
		return
	}

	err = api.core.KillSession(r.Context(), session.Value)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.SessionKilledError, err, api.logger)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	util.SendResponse(w, r, http.StatusOK, nil, variables.StatusOkMessage, nil, api.logger)
}
