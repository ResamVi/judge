package main

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"
	"time"
)

//go:embed submission
var sourceCode string

func main() {
	token := createSubmission()

	for {
		status, stdout := getSubmission(token)
		fmt.Println(status, stdout)

		if status == "ACCEPTED" {
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func getSubmission(token string) (string, string) {
	str := fmt.Sprintf("http://localhost:2358/submissions/%s?base64_encoded=true", token)

	req, err := http.NewRequest("GET", str, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Auth-Token", "mypassword")
	//req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}

	type response struct {
		Stdout  string `json:"stdout"`
		Time    string `json:"time"`
		Memory  int    `json:"memory"`
		Stderr  string `json:"stderr"`
		Token   string `json:"token"`
		Message string `json:"message"`
		Status  struct {
			Id          int    `json:"id"`
			Description string `json:"description"`
		} `json:"status"`
	}

	var body response
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		panic(err)
	}

	fmt.Println(body)

	return body.Status.Description, body.Stdout
}

func createSubmission() (token string) {
	// expected_output
	str := fmt.Sprintf(`{"source_code": "%s", "language_id": 60}`, base64.StdEncoding.EncodeToString([]byte(sourceCode)))

	req, err := http.NewRequest("POST", "http://localhost:2358/submissions?base64_encoded=true", bytes.NewBuffer([]byte(str)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Auth-Token", "mypassword")
	req.Header.Set("Content-Type", "application/json")

	// TODO: better client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		panic(resp.Status)
	}

	type response struct {
		Token string `json:"token"`
	}

	var body response
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		panic(err)
	}

	return body.Token
}
