package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (a api) GetRawPayload(requestURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("response status %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("readAll error: %s", err)
	}

	return body, nil
}

func (a api) GetUnmarshalJson(requestURL string, v any) error {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err
	}

	response, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("readAll error: %s", err)
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("invalid output (HTTP Code %d): %s", response.StatusCode, string(body))
	}

	if !json.Valid(body) {
		return fmt.Errorf("response is not JSON")
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("data unmarshal error: %s", err)
	}

	return nil
}
