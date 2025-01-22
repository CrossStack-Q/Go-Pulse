package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/CrossStack-Q/Go-Pulse/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	userId := 1

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// Change after auth
		UserID: int64(userId),
	}

	ctx := r.Context()

	err := app.store.Posts.Create(ctx, post)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = writeJSON(w, http.StatusCreated, post)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, id)

	log.Println(post)

	if err == nil {
		writeJSON(w, http.StatusOK, post)
		return
	}

	comments, err := app.store.Comments.GetByPostID(ctx, id)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}

	post.Comments = comments

	writeJSONError(w, http.StatusInternalServerError, err.Error())

	return

}
