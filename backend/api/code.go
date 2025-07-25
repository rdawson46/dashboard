package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CodeResult struct {
	Result string `json:"result,omitempty"`
	Error string `json:"error,omitempty"`
}

func ExecutePython(code string) (*CodeResult, error) {
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

	req, err := http.NewRequest(http.MethodPost, url + "/run", bodyReader)

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

	var result CodeResult
	err = json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println("Error when marshaling")
		fmt.Println(err.Error())
		fmt.Println(string(body))
		return nil, err
	}

	return &result, nil
}
