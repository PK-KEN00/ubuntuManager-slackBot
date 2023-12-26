package main

import (
	"../Documents/ubuntuManage-slackBot/systemctl"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	slackWebhookURL = "YOUR_SLACK_WEBHOOK_URL"
)

type SlackMessage struct {
	Text string `json:"text"`
}

func handleSlackRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var slackReq struct {
		Text string `json:"text"`
	}

	err = json.Unmarshal(body, &slackReq)
	if err != nil {
		log.Printf("Failed to parse Slack request JSON: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := processUserMessage(slackReq.Text)

	sendSlackMessage(response)

	w.WriteHeader(http.StatusOK)
}

func processUserMessage(message string) string {
	parts := strings.SplitN(message, " ", 3)
	if len(parts) != 3 {
		return "Invalid command. Please use the following format:\n" +
			"`systemctl <status|restart> servicename`"
	}

	command := parts[0]
	service := parts[2]

	switch command {
	case "systemctl":
		switch parts[1] {
		case "status":
			return systemctl.RunServiceStatus(service)
		case "restart":
			return systemctl.RunServiceRestart(service)
		default:
			return "Invalid command. Supported commands are `status` and `restart`."
		}
	default:
		return "Invalid command. Please use the `systemctl` command."
	}
}

func sendSlackMessage(message string) {

	payload, err := json.Marshal(SlackMessage{Text: message})
	if err != nil {
		log.Printf("Failed to marshal Slack message JSON: %v\n", err)
		return
	}

	resp, err := http.Post(slackWebhookURL, "application/json", strings.NewReader(string(payload)))
	if err != nil {
		log.Printf("Failed to send Slack message: %v\n", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error in closing response boby")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send Slack message. Status code: %d\n", resp.StatusCode)
	}
}
