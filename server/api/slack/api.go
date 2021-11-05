package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LastSprint/GooodBack/api/feedback/entries"
	"github.com/LastSprint/GooodBack/common"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type FeedbackService interface {
	Write(feedback entries.NewFeedback) error
}

type Api struct {
	Srv        FeedbackService
	SlackToken string
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

	fmt.Println("handleFeedbackCommand")

	fmt.Println(r.FormValue("payload"))

	if err := a.sendFeedbackForm(r.FormValue("trigger_id")); err != nil {
		log.Println("[ERR] Couldn't send form to slack ->", err.Error())
	}
}

func (a *Api) handleInteractivityWebhook(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println("[ERR] can't parse webhook form ->", err.Error())
		http.Error(w, "Can't parse payload", http.StatusBadRequest)
		return
	}

	rawPayload := r.Form.Get("payload")

	fmt.Println(r.Form)

	payload := struct {
		Type       string `json:"type"`
		CallbackId string `json:"callback_id"`
		TriggerId  string `json:"trigger_id"`
		View       struct {
			CallbackId string `json:"callback_id"`
			State      struct {
				Values map[string]interface{} `json:"values"`
			} `json:"state"`
		} `json:"view"`
	}{}

	if err := json.Unmarshal([]byte(rawPayload), &payload); err != nil {
		log.Println("[ERR] Can't parse payload as json ->", err.Error())
		http.Error(w, "Can't parse payload", http.StatusBadRequest)
		return
	}

	switch payload.Type {
	case "view_submission":
		if payload.View.CallbackId == "feedback-form" {
			if err := a.handleViewSubmission(payload.View.State.Values, w); err != nil {
				log.Println("[ERR] can't handle view submission ->", err.Error())
			}
		}
	}
}

func (a *Api) handleViewSubmission(form map[string]interface{}, w http.ResponseWriter) error {

	// im sorry for that
	// feedback_form_type -> type -> selected_option -> value

	typeObj, ok := form["feedback_form_type"].(map[string]interface{})

	if !ok {
		sendError(w, "Can't read field", "feedback_form_type")
		return fmt.Errorf("can't read reaction type")
	}

	typeObj, ok = typeObj["type"].(map[string]interface{})

	if !ok {
		sendError(w, "Can't read field", "feedback_form_type")
		return fmt.Errorf("can't read reaction type")
	}

	typeObj, ok = typeObj["selected_option"].(map[string]interface{})

	if !ok {
		sendError(w, "Can't read field", "feedback_form_type")
		return fmt.Errorf("can't read reaction type")
	}

	typeRawVal, ok := typeObj["value"].(string)

	if !ok {
		sendError(w, "Can't read field", "feedback_form_type")
		return fmt.Errorf("can't read reaction type")
	}

	typeVal, err := strconv.Atoi(typeRawVal)

	if err != nil {
		sendError(w, "Can't read field", "feedback_form_type")
		return fmt.Errorf("can't read reaction type -> %w", err)
	}

	// feedback_form_target -> target -> value

	target, ok := form["feedback_form_target"].(map[string]interface{})

	if !ok {
		sendError(w, "Can't read field", "feedback_form_target")
		return fmt.Errorf("can't read target")
	}

	target, ok = target["target"].(map[string]interface{})

	if !ok {
		sendError(w, "Can't read field", "feedback_form_target")
		return fmt.Errorf("can't read target")
	}

	targetVal, ok := target["value"].(string)

	if !ok {
		sendError(w, "Can't read field", "feedback_form_target")
		return fmt.Errorf("can't read target")
	}

	// feedback_form_message -> message -> value

	message, ok := form["feedback_form_message"].(map[string]interface{})

	if !ok {
		sendError(w, "Can't read field", "feedback_form_message")
		return fmt.Errorf("can't read message")
	}

	message, ok = message["message"].(map[string]interface{})

	if !ok {
		sendError(w, "Can't read field", "feedback_form_message")
		return fmt.Errorf("can't read message")
	}

	messageVal, ok := message["value"].(string)

	if !ok {
		sendError(w, "Can't read field", "feedback_form_message")
		return fmt.Errorf("can't read message")
	}

	fmt.Println(typeVal, targetVal, messageVal)

	err = a.Srv.Write(entries.NewFeedback{
		Message: messageVal,
		Target:  targetVal,
		Type:    typeVal,
	})

	if errors.Is(err, common.NotFound) {
		sendError(w, "User with this email not found", "feedback_form_target")
		return err
	}

	if err != nil {
		return fmt.Errorf("can't write feedback %w", err)
	}

	return nil
}

func (a *Api) sendFeedbackForm(triggerId string) error {

	fmt.Println(triggerId)

	dialogString := `{
	"callback_id": "feedback-form",
	"type": "modal",
	"submit": {
		"type": "plain_text",
		"text": "Send",
		"emoji": true
	},
	"close": {
		"type": "plain_text",
		"text": "Cancel",
		"emoji": true
	},
	"title": {
		"type": "plain_text",
		"text": "GoodBack Form",
		"emoji": true
	},
	"blocks": [
		{
			"block_id": "feedback_form_type",
			"type": "input",
			"label": {
				"type": "plain_text",
				"text": "Select the emotion",
				"emoji": true
			},
			"element": {
				"type": "static_select",
				"action_id": "type",
				"options": [
					{
						"text": {
							"type": "plain_text",
							"text": "🔥 cool",
							"emoji": true
						},
						"value": "2"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "👍 good",
							"emoji": true
						},
						"value": "0"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "👎 not so good",
							"emoji": true
						},
						"value": "1"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "🤬 awfull",
							"emoji": true
						},
						"value": "3"
					}
				]
			}
		},
		{
			"block_id": "feedback_form_target",
			"type": "input",
			"label": {
				"type": "plain_text",
				"text": "Target",
				"emoji": true
			},
			"element": {
				"action_id": "target",
				"type": "plain_text_input",
				"multiline": false,
				"placeholder": {
					"type": "plain_text",
					"text": "email@surfstudio.ru"
				}
			}
		},
		{
			"block_id": "feedback_form_message",
			"type": "input",
			"label": {
				"type": "plain_text",
				"text": "Message",
				"emoji": true
			},
			"element": {
				"action_id": "message",
				"type": "plain_text_input",
				"multiline": true
			}
		}
	]
}`

	requestData := map[string]interface{}{
		"trigger_id": triggerId,
		"view":       dialogString,
	}

	requestPayload, err := json.Marshal(requestData)

	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, "https://slack.com/api/views.open", bytes.NewBuffer(requestPayload))

	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+a.SlackToken)
	request.Header.Set("Content-type", "application/json")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	return nil
}

func sendError(w http.ResponseWriter, text, key string) {
	w.Header().Add("Content-Type", "application/json")
	//w.WriteHeader(http.StatusBadRequest)
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
