package main

import (
	"fmt"
	"os"

	"github.com/authora-dev/authora-go"
)

func main() {
	client := authora.NewClient(os.Getenv("AUTHORA_API_KEY"))

	agent, err := client.Agents.Create(authora.CreateAgentInput{
		Name:        "my-agent",
		WorkspaceId: "ws_...",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created agent:", agent.Id)
}
