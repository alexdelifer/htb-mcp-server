package tools

import (
	"context"
	"fmt"

	"github.com/NoASLR/htb-mcp-server/pkg/htb"
	"github.com/NoASLR/htb-mcp-server/pkg/mcp"
)

// GetChallengeInfo tool for getting detailed challenge information
type GetChallengeInfo struct {
	client *htb.Client
}

func NewGetChallengeInfo(client *htb.Client) *GetChallengeInfo {
	return &GetChallengeInfo{client: client}
}

func (t *GetChallengeInfo) Name() string {
	return "get_challenge_info"
}

func (t *GetChallengeInfo) Description() string {
	return "Get detailed challenge info including description, creator, docker status, download availability, first blood, and difficulty breakdown"
}

func (t *GetChallengeInfo) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"challenge_id": {
				Type:        "integer",
				Description: "The ID of the challenge",
			},
		},
		Required: []string{"challenge_id"},
	}
}

func (t *GetChallengeInfo) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	challengeID, ok := args["challenge_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("challenge_id is required")
	}

	endpoint := fmt.Sprintf("/challenge/info/%d", int(challengeID))

	data, err := t.client.GetWithParsing(ctx, endpoint, "challenge")
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge info: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// SpawnChallengeContainer tool for starting a Docker challenge instance
type SpawnChallengeContainer struct {
	client *htb.Client
}

func NewSpawnChallengeContainer(client *htb.Client) *SpawnChallengeContainer {
	return &SpawnChallengeContainer{client: client}
}

func (t *SpawnChallengeContainer) Name() string {
	return "spawn_challenge_container"
}

func (t *SpawnChallengeContainer) Description() string {
	return "Spawn a Docker container instance for a challenge. Returns instance ID and connection details. Use get_challenge_info to check docker_ip/docker_ports after spawning."
}

func (t *SpawnChallengeContainer) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"challenge_id": {
				Type:        "integer",
				Description: "The ID of the challenge to spawn a container for",
			},
		},
		Required: []string{"challenge_id"},
	}
}

func (t *SpawnChallengeContainer) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	challengeID, ok := args["challenge_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("challenge_id is required")
	}

	payload := htb.ContainerActionRequest{
		ContainerableID: int(challengeID),
	}

	data, err := t.client.PostWithParsing(ctx, "/container/start", payload, "")
	if err != nil {
		return nil, fmt.Errorf("failed to spawn challenge container: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// StopChallengeContainer tool for stopping a Docker challenge instance
type StopChallengeContainer struct {
	client *htb.Client
}

func NewStopChallengeContainer(client *htb.Client) *StopChallengeContainer {
	return &StopChallengeContainer{client: client}
}

func (t *StopChallengeContainer) Name() string {
	return "stop_challenge_container"
}

func (t *StopChallengeContainer) Description() string {
	return "Stop a running Docker container instance for a challenge"
}

func (t *StopChallengeContainer) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"challenge_id": {
				Type:        "integer",
				Description: "The ID of the challenge whose container to stop",
			},
		},
		Required: []string{"challenge_id"},
	}
}

func (t *StopChallengeContainer) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	challengeID, ok := args["challenge_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("challenge_id is required")
	}

	payload := htb.ContainerActionRequest{
		ContainerableID: int(challengeID),
	}

	data, err := t.client.PostWithParsing(ctx, "/container/stop", payload, "")
	if err != nil {
		return nil, fmt.Errorf("failed to stop challenge container: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// DownloadChallenge tool for getting a signed download URL for challenge files
type DownloadChallenge struct {
	client *htb.Client
}

func NewDownloadChallenge(client *htb.Client) *DownloadChallenge {
	return &DownloadChallenge{client: client}
}

func (t *DownloadChallenge) Name() string {
	return "download_challenge"
}

func (t *DownloadChallenge) Description() string {
	return "Get a signed, time-limited download URL for challenge files (zip archive)"
}

func (t *DownloadChallenge) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"challenge_id": {
				Type:        "integer",
				Description: "The ID of the challenge",
			},
		},
		Required: []string{"challenge_id"},
	}
}

func (t *DownloadChallenge) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	challengeID, ok := args["challenge_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("challenge_id is required")
	}

	endpoint := fmt.Sprintf("/challenges/%d/download_link", int(challengeID))

	// This endpoint returns {url, expires_in} at the top level
	data, err := t.client.GetWithParsing(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge download link: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}
