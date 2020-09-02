package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getAllocatedBrowserURL(dpUrl string) string {
	resp, err := http.Get(fmt.Sprintf("%s/json/version", "http://0.0.0.0:9222/json/version"))
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	return result["webSocketDebuggerUrl"].(string)
}
