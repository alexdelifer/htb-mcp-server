package tools

import (
	"context"
	"fmt"

	"github.com/NoASLR/htb-mcp-server/pkg/htb"
	"github.com/NoASLR/htb-mcp-server/pkg/mcp"
)

// GetVPNStatus tool for checking HTB VPN connection status
type GetVPNStatus struct {
	client *htb.Client
}

func NewGetVPNStatus(client *htb.Client) *GetVPNStatus {
	return &GetVPNStatus{client: client}
}

func (t *GetVPNStatus) Name() string {
	return "get_vpn_status"
}

func (t *GetVPNStatus) Description() string {
	return "Get HTB VPN connection status including assigned server, IP addresses, and bandwidth. Use this to verify VPN connectivity before starting machines."
}

func (t *GetVPNStatus) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type:       "object",
		Properties: map[string]mcp.Property{},
	}
}

func (t *GetVPNStatus) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	// /connection/status returns an array at top level, not wrapped in a key
	data, err := t.client.GetWithParsing(ctx, "/connection/status", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get VPN status: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// GetActiveResources tool for checking what machines/containers are currently running
type GetActiveResources struct {
	client *htb.Client
}

func NewGetActiveResources(client *htb.Client) *GetActiveResources {
	return &GetActiveResources{client: client}
}

func (t *GetActiveResources) Name() string {
	return "get_active_resources"
}

func (t *GetActiveResources) Description() string {
	return "Check what machines and challenge containers are currently active/running. Useful before spawning to avoid conflicts."
}

func (t *GetActiveResources) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type:       "object",
		Properties: map[string]mcp.Property{},
	}
}

func (t *GetActiveResources) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	result := make(map[string]interface{})

	// Check active machine
	machineData, err := t.client.GetWithParsing(ctx, "/machine/active", "info")
	if err != nil {
		result["machine_error"] = err.Error()
	} else if machineData == nil {
		result["active_machine"] = nil
	} else {
		result["active_machine"] = machineData
	}

	// Check VPN connections (shows what types of connections are active)
	connData, err := t.client.GetWithParsing(ctx, "/connection/status", "")
	if err != nil {
		result["connections_error"] = err.Error()
	} else {
		result["vpn_connections"] = connData
	}

	content, err := mcp.CreateJSONContent(result)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// ListChallengeCategories tool for listing challenge categories
type ListChallengeCategories struct {
	client *htb.Client
}

func NewListChallengeCategories(client *htb.Client) *ListChallengeCategories {
	return &ListChallengeCategories{client: client}
}

func (t *ListChallengeCategories) Name() string {
	return "list_challenge_categories"
}

func (t *ListChallengeCategories) Description() string {
	return "Get the list of challenge categories (Web, Pwn, Crypto, Forensics, Reversing, etc.) with their IDs"
}

func (t *ListChallengeCategories) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type:       "object",
		Properties: map[string]mcp.Property{},
	}
}

func (t *ListChallengeCategories) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	data, err := t.client.GetWithParsing(ctx, "/challenge/categories/list", "info")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch challenge categories: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// GetRecommended tool for getting HTB-recommended machine and challenge
type GetRecommended struct {
	client *htb.Client
}

func NewGetRecommended(client *htb.Client) *GetRecommended {
	return &GetRecommended{client: client}
}

func (t *GetRecommended) Name() string {
	return "get_recommended"
}

func (t *GetRecommended) Description() string {
	return "Get HTB's recommended machine and challenge picks (staff picks, seasonal, etc.)"
}

func (t *GetRecommended) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"type": {
				Type:        "string",
				Description: "Get recommendations for machines or challenges",
				Enum:        []string{"machines", "challenges"},
				Default:     "machines",
			},
		},
	}
}

func (t *GetRecommended) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	recType := "machines"
	if rt, ok := args["type"].(string); ok {
		recType = rt
	}

	var endpoint string
	switch recType {
	case "challenges":
		endpoint = "/challenge/recommended"
	default:
		endpoint = "/machine/recommended"
	}

	data, err := t.client.GetWithParsing(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}
