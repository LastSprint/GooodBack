package services

import (
	"encoding/json"
	"fmt"
	"github.com/LastSprint/GooodBack/api/feedback/entries"
	"github.com/LastSprint/GooodBack/api/slack/errors"
	"strconv"
)

type FeedbackService interface {
	Write(feedback entries.NewFeedback) error
}

type SlackInteractivityWebHandler struct {
	FeedbackSrv FeedbackService
}

func (s *SlackInteractivityWebHandler) Handle(rawPayload string) error {

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
		return fmt.Errorf("couldn't parse payload as json -> %w", err.Error())
	}

	switch payload.Type {
	case "view_submission":
		if payload.View.CallbackId == "feedback-form" {
			feedback, err := s.handleViewSubmission(payload.View.State.Values)
			if err != nil {
				return fmt.Errorf("can't handle view submission -> %w", err)
			}

			if err = s.FeedbackSrv.Write(*feedback); err != nil {
				return fmt.Errorf("couldn't write feedback %v -> %w", feedback, err)
			}
		}
	default:
		return fmt.Errorf("payload with type %s is not supported", payload.Type)
	}

	return nil
}

func (s *SlackInteractivityWebHandler) handleViewSubmission(form map[string]interface{}) (*entries.NewFeedback, error) {

	// im sorry for that
	// feedback_form_type -> type -> selected_option -> value

	typeObj, ok := form["feedback_form_type"].(map[string]interface{})

	if !ok {
		return nil, &errors.SlackViewSubmissionError{
			Key:   "feedback_form_type",
			Value: "Can't read field",
		}
	}

	typeObj, ok = typeObj["type"].(map[string]interface{})

	if !ok {
		return nil, &errors.SlackViewSubmissionError{
			Key:   "feedback_form_type",
			Value: "Can't read reaction type",
		}
	}

	typeObj, ok = typeObj["selected_option"].(map[string]interface{})

	if !ok {
		return nil, &errors.SlackViewSubmissionError{
			Key:   "feedback_form_type",
			Value: "Can't read reaction type",
		}
	}

	typeRawVal, ok := typeObj["value"].(string)

	if !ok {
		return nil, &errors.SlackViewSubmissionError{
			Key:   "feedback_form_type",
			Value: "Can't read reaction type",
		}
	}

	typeVal, err := strconv.Atoi(typeRawVal)

	if err != nil {
		return nil, &errors.SlackViewSubmissionError{
			Key:   "feedback_form_type",
			Value: "Can't parse reaction type",
		}
	}

	// feedback_form_target -> target -> value

	target, ok := form["feedback_form_target"].(map[string]interface{})

	if !ok {
		return nil, &errors.SlackViewSubmissionError{
			Key:   "feedback_form_target",
			Value: "Can't read field",
		}
	}

	target, ok = target["target"].(map[string]interface{})

	if !ok {
		return nil, &errors.SlackViewSubmissionError{
			Key:   "feedback_form_target",
			Value: "Can't read target",
		}
	}

	targetVal, ok := target["value"].(string)

	if !ok {
		return nil, &errors.SlackViewSubmissionError{
			Key:   "feedback_form_target",
			Value: "Can't read target",
		}
	}

	// feedback_form_message -> message -> value

	message, ok := form["feedback_form_message"].(map[string]interface{})

	if !ok {
		return nil, &errors.SlackViewSubmissionError{
			Key:   "feedback_form_message",
			Value: "Can't read field",
		}
	}

	message, ok = message["message"].(map[string]interface{})

	if !ok {
		return nil, &errors.SlackViewSubmissionError{
			Key:   "feedback_form_message",
			Value: "Can't read feedback_form_message.message",
		}
	}

	messageVal, ok := message["value"].(string)

	if !ok {
		return nil, &errors.SlackViewSubmissionError{
			Key:   "feedback_form_message",
			Value: "Can't read read feedback_form_message.message.value as string",
		}
	}

	return &entries.NewFeedback{
		Message: messageVal,
		Target:  targetVal,
		Type:    typeVal,
	}, nil

}
