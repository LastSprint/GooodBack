package feedback

import (
	"encoding/json"
	"errors"
	"github.com/LastSprint/GooodBack/api/feedback/entries"
	"github.com/LastSprint/GooodBack/common/middlewares"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Service interface {
	Write(feedback entries.NewFeedback) error
	Read(userId string) ([]entries.Feedback, error)
}

type Api struct {
	Srv Service
}

func (a *Api) Start(r chi.Router) {
	r.Route("/feedback", func(r chi.Router) {
		r.Post("/", a.handleCreateFeedback)
		r.With(middlewares.RejectInvalidTokens).Get("/", a.readFeedback)
	})
}

func (a *Api) handleCreateFeedback(w http.ResponseWriter, r *http.Request) {
	entry := entries.NewFeedback{}

	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "couldn't parse request", http.StatusBadRequest)
		return
	}

	err := a.Srv.Write(entry)

	if errors.Is(err, FeedbackTargetNotFound) {
		http.Error(w, "target not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}

func (a *Api) readFeedback(w http.ResponseWriter, r *http.Request) {

	result, err := a.Srv.Read(r.Context().Value(middlewares.ContextKeyUserId).(string))

	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("[ERR] can't encode feedback arr to json with error %s", err.Error())
		http.Error(w, "cant encode result", http.StatusInternalServerError)
	}
}
