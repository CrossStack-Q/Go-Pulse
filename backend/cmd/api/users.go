package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/CrossStack-Q/Go-Pulse/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {

	followerUser := getUserFromContext(r)

	var payload FollowUser

	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	// Temp checking User exists
	_, err := app.store.Users.GetUserByID(ctx, payload.UserID)
	if err != nil {
		log.Println("user_id wala user is not present")
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	err = app.store.Followers.Follow(ctx, followerUser.ID, payload.UserID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	followerUser := getUserFromContext(r)

	var payload FollowUser

	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, followerUser.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetUserByID(ctx, userID)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
