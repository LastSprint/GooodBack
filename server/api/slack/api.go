package slack

import (
	"encoding/json"
	"errors"
	slackErrors "github.com/LastSprint/GooodBack/api/slack/errors"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type CommandHandler interface {
	Handle(input url.Values) error
}

type WebhookHandler interface {
	Handle(payload string) error
}

type Api struct {
	CmdHandler         CommandHandler
	InteractionHandler WebhookHandler
}

func (a *Api) Start(r chi.Router) {
	r.Route("/feedback/slack", func(r chi.Router) {
		r.Post("/modal-form", a.handleFeedbackCommand)
		r.Post("/webhook", a.handleInteractivityWebhook)
	})
}

func (a *Api) handleFeedbackCommand(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		data, _ := ioutil.ReadAll(r.Body)
		log.Printf("[ERR] couldn't parse slack request %s with error %s", string(data), err.Error())
		http.Error(w, "can't parse slack request", http.StatusBadRequest)
		return
	}

	if err := a.CmdHandler.Handle(r.Form); err != nil {
		log.Println("[ERR] handleFeedbackCommand ->", err.Error())
		http.Error(w, "cant handle feedback command", http.StatusBadRequest)
	}
}

func (a *Api) handleInteractivityWebhook(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println("[ERR] can't parse webhook form ->", err.Error())
		http.Error(w, "Can't parse payload", http.StatusBadRequest)
		return
	}

	rawPayload := r.Form.Get("payload")

	if len(rawPayload) == 0 {
		log.Println("[ERR] payload is null Request:", r.URL.String())
		http.Error(w, "Payload is empty", http.StatusBadRequest)
		return
	}

	err := a.InteractionHandler.Handle(rawPayload)

	if err == nil {
		return
	}

	log.Printf("[ERR] when handling interaction %s -> %s\n", r.URL.String(), err.Error())

	var slackErr *slackErrors.SlackViewSubmissionError

	if errors.As(err, &slackErr) {
		sendError(w, slackErr.Key, slackErr.Value)
		return
	}

	http.Error(w, "something went wrong", http.StatusInternalServerError)
}

func sendError(w http.ResponseWriter, text, key string) {
	w.Header().Add("Content-Type", "application/json")

	object := struct {
		ResponseAction string            `json:"response_action"`
		Errors         map[string]string `json:"errors"`
	}{
		ResponseAction: "errors",
		Errors: map[string]string{
			key: text,
		},
	}

	if err := json.NewEncoder(w).Encode(object); err != nil {
		log.Println("[ERR] couldn't encode error object ->", err.Error())
	}
}
