package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/CrossStack-Q/Go-Pulse/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	userId := 3

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

	post := getPostFromCtx(r)

	// Here is Comment Logic
	// Too Complex
	// Worst Part

	if post == nil {
		writeJSON(w, http.StatusAccepted, post)
		return
	}

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}

	if err == nil {
		writeJSON(w, http.StatusOK, post)
		return
	}

	post.Comments = comments

	writeJSONError(w, http.StatusInternalServerError, err.Error())

	return

}

func (app *application) deletePost(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "postID")

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	if err := app.store.Posts.DeletePostByID(ctx, id); err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, nil)
	return

}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePost(w http.ResponseWriter, r *http.Request) {

	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	ctx := r.Context()

	if err := app.store.Posts.UpdatePostByID(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusAccepted, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetByID(ctx, id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {

	post, ok := r.Context().Value(postCtx).(store.Post)
	if !ok {
		fmt.Println("Error: No post found in context")
		fmt.Println("Context Value: ", r.Context().Value(postCtx))
		fmt.Println("Type of Context Value: ", reflect.TypeOf(r.Context().Value(postCtx)))
		return nil
	}
	return &post
}
