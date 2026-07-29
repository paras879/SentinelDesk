package windowsservice

import (
	"fmt"
	"log"


	"sentineldesk/agent/internal/api"
	"sentineldesk/agent/internal/deviceid"
)

type PendingCommand struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	Action      string `json:"action"`
	Status      string `json:"status"`
}

type PendingCommandsResponse struct {
	Count    int              `json:"count"`
	Commands []PendingCommand `json:"commands"`
}

type CommandResultRequest struct {
	CommandID    string `json:"command_id"`
	Result       string `json:"result"`
	ErrorMessage string `json:"error_message"`
}

func PollAndExecuteCommands() {

	client := api.NewClient()

	var respData PendingCommandsResponse

	resp, err := client.R().
		SetQueryParam("device_id", deviceid.Get()).
		SetResult(&respData).
		Get("/api/v1/agent/commands")

	if err != nil {
		log.Println("CommandPoll: HTTP error:", err)
		return
	}

	if resp.StatusCode() != 200 {
		return
	}

	for _, cmd := range respData.Commands {
		executeAndReport(cmd)
	}
}

func executeAndReport(cmd PendingCommand) {

	result := "success"
	errMsg := ""

	switch cmd.Action {
	case "start":
		errMsg = startService(cmd.ServiceName)
	case "stop":
		errMsg = stopService(cmd.ServiceName)
	case "restart":
		errMsg = restartService(cmd.ServiceName)
	default:
		result = "failure"
		errMsg = "unknown action: " + cmd.Action
	}

	if errMsg != "" {
		result = "failure"
	}

	client := api.NewClient()

	req := CommandResultRequest{
		CommandID:    cmd.ID,
		Result:       result,
		ErrorMessage: errMsg,
	}

	resp, err := client.R().
		SetBody(req).
		Post("/api/v1/agent/commands/result")

	if err != nil {
		log.Printf("CommandPoll: Failed to report result for %s: %v", cmd.ID, err)
		return
	}

	fmt.Printf("CommandPoll: %s %s -> %s (%s)\n",
		cmd.Action, cmd.ServiceName, result, resp.Status())
}


