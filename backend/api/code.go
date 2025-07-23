package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type CodeResult struct {
	Result string `json:"result"`
	Error string `json:"error"`
}

func executePython(code string) (*CodeResult, error) {
	reqData := map[string]string{"code": code}

	jsonReq, err := json.Marshal(reqData)

	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewBuffer(jsonReq)

	url := os.Getenv("CODE_URL")

	if url == "" {
		return nil, errors.New("No Code URL in ENV")
	}

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var result *CodeResult
	err = json.Unmarshal(body, result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
