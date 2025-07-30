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

func WebSearch(query string) (*ToolResult, error) {
	reqData := map[string]string{"query": query}

	jsonReq, err := json.Marshal(reqData)

	if err != nil {
		return nil, err
	}

	// TODO: rename this route 
	url := os.Getenv("CODE_URL")

	if url == "" {
		return nil, errors.New("No Code URL in ENV")
	}

	bodyReader := bytes.NewBuffer(jsonReq)

	req, err := http.NewRequest(http.MethodPost, url + "/websearch", bodyReader)

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

	var result ToolResult
	err = json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println("Error when marshaling")
		fmt.Println(err.Error())
		fmt.Println(string(body))
		return nil, err
	}

	return &result, nil
}
