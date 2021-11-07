package repos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SlackRepo struct {
	SlackToken string
}

func (s *SlackRepo) SendView(view string, triggerID string) error {
	requestData := map[string]interface{}{
		"trigger_id": triggerID,
		"view":       view,
	}

	requestPayload, err := json.Marshal(requestData)

	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, "https://slack.com/api/views.open", bytes.NewBuffer(requestPayload))

	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+s.SlackToken)
	request.Header.Set("Content-type", "application/json")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("POST %s returned status code %v with body %s", request.URL.String(), response.StatusCode, string(body))
	}

	return nil
}
