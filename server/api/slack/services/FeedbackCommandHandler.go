package services

import (
	"fmt"
	"net/url"
)

var dialogString = `{
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

type SlackFeedbackFormSender interface {
	SendView(view string, triggerID string) error
}

type FeedbackCommandHandler struct {
	Repo SlackFeedbackFormSender
}

func (h *FeedbackCommandHandler) Handle(input url.Values) error {
	triggerID := input.Get("trigger_id")

	if len(triggerID) == 0 {
		return fmt.Errorf("FeedbackCommandHandler couldn't get trigger_id from input %v", input)
	}

	if err := h.Repo.SendView(dialogString, triggerID); err != nil {
		return fmt.Errorf("couldn't send view with error -> %w", err)
	}

	return nil
}
