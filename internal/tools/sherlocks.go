package tools

import (
	"context"
	"fmt"

	"github.com/NoASLR/htb-mcp-server/pkg/htb"
	"github.com/NoASLR/htb-mcp-server/pkg/mcp"
)

// ListSherlocks tool for listing HTB Sherlocks (DFIR challenges)
type ListSherlocks struct {
	client *htb.Client
}

func NewListSherlocks(client *htb.Client) *ListSherlocks {
	return &ListSherlocks{client: client}
}

func (t *ListSherlocks) Name() string {
	return "list_sherlocks"
}

func (t *ListSherlocks) Description() string {
	return "Get a paginated list of HackTheBox Sherlocks (DFIR investigation challenges)"
}

func (t *ListSherlocks) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"per_page": {
				Type:        "integer",
				Description: "Number of Sherlocks per page",
				Default:     15,
			},
		},
	}
}

func (t *ListSherlocks) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	perPage := 15
	if pp, ok := args["per_page"].(float64); ok {
		perPage = int(pp)
	}

	endpoint := fmt.Sprintf("/sherlocks?per_page=%d", perPage)

	data, err := t.client.GetWithParsing(ctx, endpoint, "data")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sherlocks: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// ListSherlockCategories tool for listing Sherlock categories
type ListSherlockCategories struct {
	client *htb.Client
}

func NewListSherlockCategories(client *htb.Client) *ListSherlockCategories {
	return &ListSherlockCategories{client: client}
}

func (t *ListSherlockCategories) Name() string {
	return "list_sherlock_categories"
}

func (t *ListSherlockCategories) Description() string {
	return "Get the list of Sherlock categories (DFIR, SOC, Malware Analysis, etc.)"
}

func (t *ListSherlockCategories) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type:       "object",
		Properties: map[string]mcp.Property{},
	}
}

func (t *ListSherlockCategories) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	data, err := t.client.GetWithParsing(ctx, "/sherlocks/categories/list", "info")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sherlock categories: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// GetSherlockInfo tool for getting Sherlock details by ID
type GetSherlockInfo struct {
	client *htb.Client
}

func NewGetSherlockInfo(client *htb.Client) *GetSherlockInfo {
	return &GetSherlockInfo{client: client}
}

func (t *GetSherlockInfo) Name() string {
	return "get_sherlock_info"
}

func (t *GetSherlockInfo) Description() string {
	return "Get detailed information about a Sherlock including description and linked Academy modules"
}

func (t *GetSherlockInfo) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"sherlock_id": {
				Type:        "integer",
				Description: "The ID of the Sherlock",
			},
		},
		Required: []string{"sherlock_id"},
	}
}

func (t *GetSherlockInfo) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	sherlockID, ok := args["sherlock_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("sherlock_id is required")
	}

	endpoint := fmt.Sprintf("/sherlocks/%d/info", int(sherlockID))

	data, err := t.client.GetWithParsing(ctx, endpoint, "data")
	if err != nil {
		return nil, fmt.Errorf("failed to get sherlock info: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// GetSherlockPlay tool for getting the Sherlock play view (scenario, creators, instance status)
type GetSherlockPlay struct {
	client *htb.Client
}

func NewGetSherlockPlay(client *htb.Client) *GetSherlockPlay {
	return &GetSherlockPlay{client: client}
}

func (t *GetSherlockPlay) Name() string {
	return "get_sherlock_play"
}

func (t *GetSherlockPlay) Description() string {
	return "Get Sherlock play view including scenario text, creators, first blood, file info, and Docker instance status (play_info)"
}

func (t *GetSherlockPlay) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"sherlock_id": {
				Type:        "integer",
				Description: "The ID of the Sherlock",
			},
		},
		Required: []string{"sherlock_id"},
	}
}

func (t *GetSherlockPlay) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	sherlockID, ok := args["sherlock_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("sherlock_id is required")
	}

	endpoint := fmt.Sprintf("/sherlocks/%d/play", int(sherlockID))

	data, err := t.client.GetWithParsing(ctx, endpoint, "data")
	if err != nil {
		return nil, fmt.Errorf("failed to get sherlock play view: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// GetSherlockTasks tool for listing tasks/questions for a Sherlock
type GetSherlockTasks struct {
	client *htb.Client
}

func NewGetSherlockTasks(client *htb.Client) *GetSherlockTasks {
	return &GetSherlockTasks{client: client}
}

func (t *GetSherlockTasks) Name() string {
	return "get_sherlock_tasks"
}

func (t *GetSherlockTasks) Description() string {
	return "Get the list of tasks (questions) for a Sherlock, including completion status and expected flag format"
}

func (t *GetSherlockTasks) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"sherlock_id": {
				Type:        "integer",
				Description: "The ID of the Sherlock",
			},
		},
		Required: []string{"sherlock_id"},
	}
}

func (t *GetSherlockTasks) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	sherlockID, ok := args["sherlock_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("sherlock_id is required")
	}

	endpoint := fmt.Sprintf("/sherlocks/%d/tasks", int(sherlockID))

	data, err := t.client.GetWithParsing(ctx, endpoint, "data")
	if err != nil {
		return nil, fmt.Errorf("failed to get sherlock tasks: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// GetSherlockProgress tool for getting user progress on a Sherlock
type GetSherlockProgress struct {
	client *htb.Client
}

func NewGetSherlockProgress(client *htb.Client) *GetSherlockProgress {
	return &GetSherlockProgress{client: client}
}

func (t *GetSherlockProgress) Name() string {
	return "get_sherlock_progress"
}

func (t *GetSherlockProgress) Description() string {
	return "Get user progress on a Sherlock including tasks answered, total tasks, and ownership status"
}

func (t *GetSherlockProgress) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"sherlock_id": {
				Type:        "integer",
				Description: "The ID of the Sherlock",
			},
		},
		Required: []string{"sherlock_id"},
	}
}

func (t *GetSherlockProgress) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	sherlockID, ok := args["sherlock_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("sherlock_id is required")
	}

	endpoint := fmt.Sprintf("/sherlocks/%d/progress", int(sherlockID))

	data, err := t.client.GetWithParsing(ctx, endpoint, "data")
	if err != nil {
		return nil, fmt.Errorf("failed to get sherlock progress: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// SubmitSherlockFlag tool for submitting a task answer for a Sherlock
type SubmitSherlockFlag struct {
	client *htb.Client
}

func NewSubmitSherlockFlag(client *htb.Client) *SubmitSherlockFlag {
	return &SubmitSherlockFlag{client: client}
}

func (t *SubmitSherlockFlag) Name() string {
	return "submit_sherlock_flag"
}

func (t *SubmitSherlockFlag) Description() string {
	return "Submit an answer (flag) for a specific task in a Sherlock"
}

func (t *SubmitSherlockFlag) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"sherlock_id": {
				Type:        "integer",
				Description: "The ID of the Sherlock",
			},
			"task_id": {
				Type:        "integer",
				Description: "The ID of the task to submit the answer for",
			},
			"flag": {
				Type:        "string",
				Description: "The answer/flag to submit",
			},
		},
		Required: []string{"sherlock_id", "task_id", "flag"},
	}
}

func (t *SubmitSherlockFlag) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	sherlockID, ok := args["sherlock_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("sherlock_id is required")
	}

	taskID, ok := args["task_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("task_id is required")
	}

	flag, ok := args["flag"].(string)
	if !ok {
		return nil, fmt.Errorf("flag is required")
	}

	endpoint := fmt.Sprintf("/sherlocks/%d/tasks/%d/flag", int(sherlockID), int(taskID))

	payload := htb.SherlockFlagRequest{
		Flag: flag,
	}

	data, err := t.client.PostWithParsing(ctx, endpoint, payload, "")
	if err != nil {
		return nil, fmt.Errorf("failed to submit sherlock flag: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// DownloadSherlock tool for getting a signed download URL for Sherlock artifacts
type DownloadSherlock struct {
	client *htb.Client
}

func NewDownloadSherlock(client *htb.Client) *DownloadSherlock {
	return &DownloadSherlock{client: client}
}

func (t *DownloadSherlock) Name() string {
	return "download_sherlock"
}

func (t *DownloadSherlock) Description() string {
	return "Get a signed, time-limited download URL for Sherlock investigation artifacts"
}

func (t *DownloadSherlock) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"sherlock_id": {
				Type:        "integer",
				Description: "The ID of the Sherlock",
			},
		},
		Required: []string{"sherlock_id"},
	}
}

func (t *DownloadSherlock) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	sherlockID, ok := args["sherlock_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("sherlock_id is required")
	}

	endpoint := fmt.Sprintf("/sherlocks/%d/download_link", int(sherlockID))

	// This endpoint returns {url, expires_in} at the top level (not wrapped in "data")
	data, err := t.client.GetWithParsing(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get sherlock download link: %w", err)
	}

	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}
