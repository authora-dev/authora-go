package main

import (
	"context"
	"fmt"
	"os"

	"github.com/authora-dev/authora-go"
)

func main() {
	client := authora.NewClient(os.Getenv("AUTHORA_API_KEY"))
	ctx := context.Background()

	agent, err := client.Agents.Create(ctx, &authora.CreateAgentInput{
		WorkspaceID: "ws_...",
		Name:        "my-agent",
		CreatedBy:   "quickstart",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created agent:", agent.ID)
}
