package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Requests struct {
	authToken  string
	authPrefix string
}

func NewRequests(authToken string) *Requests {
	return &Requests{
		authToken:  authToken,
		authPrefix: "Bearer",
	}
}

func (r *Requests) post(url string, params any, response any) error {

	newData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(string(newData)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if r.authToken != "" {
		req.Header.Set("Authorization", r.authPrefix+" "+r.authToken)
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		return err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	return nil
}
