package middleware

import (
	"bootcamp-task/pkg/util"
	"bootcamp-task/pkg/variables"
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type ICore interface {
	GetUserId(ctx context.Context, sid string) (int64, error)
	GetUserRole(ctx context.Context, id int64) (string, error)
}

func PanicMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.StatusInternalServerError, fmt.Errorf(fmt.Sprintf("%v", err)), logger)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func MethodMiddleware(next http.Handler, methods []string, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isMethod := false
		for _, val := range methods {
			if r.Method == val {
				isMethod = true
				break
			}
		}

		if !isMethod {
			util.SendResponse(w, r, http.StatusMethodNotAllowed, nil, variables.StatusMethodNotAllowedError, nil, logger)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthorizationMiddleware(next http.Handler, core ICore, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie(variables.SessionCookieName)
		if err != nil {
			util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.StatusUnauthorizedError, nil, logger)
			return
		}

		userId, err := core.GetUserId(r.Context(), session.Value)
		if err != nil || userId == 0 {
			util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.StatusUnauthorizedError, nil, logger)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), variables.UserIDKey, userId))
		next.ServeHTTP(w, r)
	})
}

func PermissionsMiddleware(next http.Handler, core ICore, roles []string, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, isAuth := r.Context().Value(variables.UserIDKey).(int64)
		if !isAuth {
			util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.StatusUnauthorizedError, nil, logger)
			return
		}

		userRole, err := core.GetUserRole(r.Context(), userId)
		if err != nil {
			util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.StatusInternalServerError, err, logger)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), variables.RoleKey, userRole))

		isPermitted := false
		for _, val := range roles {
			if userRole == val {
				isPermitted = true
				break
			}
		}

		if !isPermitted {
			util.SendResponse(w, r, http.StatusForbidden, nil, variables.StatusForbiddenError, nil, logger)
			return
		}

		next.ServeHTTP(w, r)
	})
}
