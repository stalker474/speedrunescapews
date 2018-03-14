package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func Test_CreateUserFail(t *testing.T) {
	hc := http.DefaultClient
	var command JSONUser
	data, err := json.Marshal(&command)
	str := string(data)

	req, _ := http.NewRequest("POST", "http://localhost:8080/createuser", strings.NewReader(str))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	resp, err := hc.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}

	dec := json.NewDecoder(resp.Body)
	var res JSONResult
	dec.Decode(&res)

	if resp.StatusCode != 200 {
		t.Fatalf("Unexpected ws answer : %d, %s", resp.StatusCode, res.Result)
	}

	if res.Result != "Created" {
		t.Fatal("Unexpected result : " + res.Result)
	}
}
