//go:build integration

package tools_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/NoASLR/htb-mcp-server/internal/tools"
	"github.com/NoASLR/htb-mcp-server/pkg/config"
	"github.com/NoASLR/htb-mcp-server/pkg/htb"
)

func setupClient(t *testing.T) *htb.Client {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Skipf("Skipping integration test (no HTB_TOKEN): %v", err)
	}
	return htb.NewClient(cfg)
}

// parseJSONContent extracts and parses the JSON from a tool response's first content block.
func parseJSONContent(t *testing.T, content string) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// Try as array — some responses are arrays
		t.Logf("Response is not a JSON object (may be array or text): %s", content[:min(200, len(content))])
		return nil
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ==================== Tier 1: Pure reads ====================

func TestGetServerStatus(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewGetServerStatus(client)

	resp, err := tool.Execute(context.Background(), nil)
	if err != nil {
		t.Fatalf("get_server_status failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	data := parseJSONContent(t, resp.Content[0].Text)
	if data == nil {
		t.Fatal("could not parse response as JSON object")
	}

	if status, ok := data["status"].(string); !ok || status != "running" {
		t.Errorf("expected status=running, got %v", data["status"])
	}
	if htbStatus, ok := data["htb_api_status"].(string); !ok || htbStatus != "healthy" {
		t.Errorf("expected htb_api_status=healthy, got %v", data["htb_api_status"])
	}
	t.Logf("Server status: %s, HTB API: %s, uptime: %v", data["status"], data["htb_api_status"], data["uptime"])
}

func TestGetUserProfile(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewGetUserProfile(client)

	resp, err := tool.Execute(context.Background(), nil)
	if err != nil {
		t.Fatalf("get_user_profile failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	data := parseJSONContent(t, resp.Content[0].Text)
	if data == nil {
		t.Fatal("could not parse response as JSON object")
	}

	// Verify essential fields
	if _, ok := data["id"]; !ok {
		t.Error("missing 'id' field in user profile")
	}
	if name, ok := data["name"].(string); !ok || name == "" {
		t.Error("missing or empty 'name' field in user profile")
	}

	t.Logf("User: %v (id=%v)", data["name"], data["id"])
}

func TestGetUserProgress(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewGetUserProgress(client)

	resp, err := tool.Execute(context.Background(), map[string]interface{}{
		"type": "overview",
	})
	if err != nil {
		t.Fatalf("get_user_progress failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	// Just verify it returns valid JSON
	data := parseJSONContent(t, resp.Content[0].Text)
	if data == nil {
		t.Fatal("could not parse response as JSON object")
	}
	t.Logf("Progress response keys: %v", mapKeys(data))
}

func TestSearchContent(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewSearchContent(client)

	resp, err := tool.Execute(context.Background(), map[string]interface{}{
		"query": "Lame",
		"type":  "all",
	})
	if err != nil {
		t.Fatalf("search_content failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	data := parseJSONContent(t, resp.Content[0].Text)
	if data == nil {
		t.Fatal("could not parse response as JSON object")
	}

	// Should have a machines key with results
	if machines, ok := data["machines"]; ok {
		t.Logf("Search returned machines: %v", machines)
	} else {
		t.Logf("Search response keys: %v", mapKeys(data))
	}
}

func TestListMachinesActive(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewListMachines(client)

	resp, err := tool.Execute(context.Background(), map[string]interface{}{
		"status":   "active",
		"per_page": float64(5),
	})
	if err != nil {
		t.Fatalf("list_machines (active) failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	// Try parsing as array
	var machines []interface{}
	if err := json.Unmarshal([]byte(resp.Content[0].Text), &machines); err != nil {
		// Might be wrapped in an object
		t.Logf("Response (first 300 chars): %s", resp.Content[0].Text[:min(300, len(resp.Content[0].Text))])
		t.Fatalf("could not parse machines response: %v", err)
	}

	if len(machines) == 0 {
		t.Error("no active machines returned")
	}
	if len(machines) > 5 {
		t.Logf("WARNING: asked for per_page=5 but got %d machines (server-side pagination may not enforce limit)", len(machines))
	}

	// Check first machine has expected fields
	if len(machines) > 0 {
		if m, ok := machines[0].(map[string]interface{}); ok {
			t.Logf("First machine: name=%v, os=%v, difficulty=%v", m["name"], m["os"], m["difficultyText"])
		}
	}
}

func TestListMachinesRetired(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewListMachines(client)

	resp, err := tool.Execute(context.Background(), map[string]interface{}{
		"status":   "retired",
		"per_page": float64(5),
	})
	if err != nil {
		t.Fatalf("list_machines (retired) failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}
	t.Logf("Retired machines response length: %d bytes", len(resp.Content[0].Text))
}

func TestListChallenges(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewListChallenges(client)

	resp, err := tool.Execute(context.Background(), map[string]interface{}{
		"status": "active",
	})
	if err != nil {
		t.Fatalf("list_challenges failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	// Try parsing as array
	var challenges []interface{}
	if err := json.Unmarshal([]byte(resp.Content[0].Text), &challenges); err != nil {
		t.Logf("Response (first 300 chars): %s", resp.Content[0].Text[:min(300, len(resp.Content[0].Text))])
		t.Fatalf("could not parse challenges response: %v", err)
	}

	if len(challenges) == 0 {
		t.Error("no active challenges returned")
	}

	if len(challenges) > 0 {
		if c, ok := challenges[0].(map[string]interface{}); ok {
			t.Logf("First challenge: name=%v, category=%v, difficulty=%v", c["name"], c["challenge_category_id"], c["difficulty"])
		}
	}
}

func TestGetMachineInfo(t *testing.T) {
	client := setupClient(t)

	// First, search for a known machine to get its ID
	searchTool := tools.NewSearchContent(client)
	searchResp, err := searchTool.Execute(context.Background(), map[string]interface{}{
		"query": "Lame",
		"type":  "machines",
	})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	searchData := parseJSONContent(t, searchResp.Content[0].Text)
	if searchData == nil {
		t.Skip("could not parse search response")
	}

	// Extract machine ID from search results
	var machineID float64
	if machines, ok := searchData["machines"].([]interface{}); ok && len(machines) > 0 {
		if m, ok := machines[0].(map[string]interface{}); ok {
			if id, ok := m["id"].(float64); ok {
				machineID = id
			}
		}
	}
	if machineID == 0 {
		// Fallback: Lame is machine ID 1
		machineID = 1
		t.Logf("Could not extract machine ID from search, using fallback ID=1")
	}

	infoTool := tools.NewGetMachineInfo(client)
	resp, err := infoTool.Execute(context.Background(), map[string]interface{}{
		"machine_id": machineID,
	})
	if err != nil {
		t.Fatalf("get_machine_info failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	data := parseJSONContent(t, resp.Content[0].Text)
	if data == nil {
		t.Fatal("could not parse machine info as JSON object")
	}

	if name, ok := data["name"].(string); ok {
		t.Logf("Machine info: name=%s, os=%v, difficulty=%v", name, data["os"], data["difficultyText"])
	} else {
		t.Logf("Machine info keys: %v", mapKeys(data))
	}
}

// ==================== Tier 2: Reversible state changes ====================

func TestMachineLifecycle(t *testing.T) {
	if os.Getenv("RUN_LIFECYCLE_TESTS") != "1" {
		t.Skip("Skipping machine lifecycle test (set RUN_LIFECYCLE_TESTS=1 to enable)")
	}

	client := setupClient(t)

	// List active machines to find one we can start
	listTool := tools.NewListMachines(client)
	listResp, err := listTool.Execute(context.Background(), map[string]interface{}{
		"status":   "active",
		"per_page": float64(1),
	})
	if err != nil {
		t.Fatalf("list_machines failed: %v", err)
	}

	var machines []map[string]interface{}
	if err := json.Unmarshal([]byte(listResp.Content[0].Text), &machines); err != nil {
		t.Fatalf("could not parse machines: %v", err)
	}
	if len(machines) == 0 {
		t.Skip("no active machines available")
	}

	machineID := machines[0]["id"].(float64)
	machineName := machines[0]["name"]
	t.Logf("Testing lifecycle with machine: %v (id=%v)", machineName, machineID)

	// Start the machine
	startTool := tools.NewStartMachine(client)
	startResp, err := startTool.Execute(context.Background(), map[string]interface{}{
		"machine_id": machineID,
	})
	if err != nil {
		t.Fatalf("start_machine failed: %v", err)
	}
	t.Logf("Start response: %s", startResp.Content[0].Text[:min(200, len(startResp.Content[0].Text))])

	// Wait for machine to spin up
	t.Log("Waiting 10s for machine to initialize...")
	time.Sleep(10 * time.Second)

	// Get the machine IP
	ipTool := tools.NewGetMachineIP(client)
	ipResp, err := ipTool.Execute(context.Background(), nil)
	if err != nil {
		t.Errorf("get_machine_ip failed: %v", err)
	} else {
		t.Logf("Machine IP response: %s", ipResp.Content[0].Text[:min(200, len(ipResp.Content[0].Text))])
	}

	// Stop the machine
	stopTool := tools.NewStopMachine(client)
	stopResp, err := stopTool.Execute(context.Background(), map[string]interface{}{
		"machine_id": machineID,
	})
	if err != nil {
		t.Fatalf("stop_machine failed: %v", err)
	}
	t.Logf("Stop response: %s", stopResp.Content[0].Text[:min(200, len(stopResp.Content[0].Text))])
}

// ==================== Tier 3: Flag tests (gated) ====================

func TestSubmitDummyFlag(t *testing.T) {
	if os.Getenv("RUN_FLAG_TESTS") != "1" {
		t.Skip("Skipping flag submission test (set RUN_FLAG_TESTS=1 to enable)")
	}

	client := setupClient(t)
	tool := tools.NewSubmitUserFlag(client)

	// Submit a known-bad flag to machine ID 1 (Lame)
	resp, err := tool.Execute(context.Background(), map[string]interface{}{
		"machine_id": float64(1),
		"flag":       "HTB{this_is_an_invalid_test_flag}",
	})
	if err != nil {
		// An error response from the API is actually expected
		t.Logf("Flag submission returned error (expected): %v", err)
		return
	}

	t.Logf("Flag submission response: %s", resp.Content[0].Text)
}

// ==================== Sherlock read tests ====================

func TestListSherlocks(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewListSherlocks(client)

	resp, err := tool.Execute(context.Background(), map[string]interface{}{
		"per_page": float64(5),
	})
	if err != nil {
		t.Fatalf("list_sherlocks failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	var sherlocks []interface{}
	if err := json.Unmarshal([]byte(resp.Content[0].Text), &sherlocks); err != nil {
		// Might be wrapped
		data := parseJSONContent(t, resp.Content[0].Text)
		t.Logf("Sherlocks response keys: %v", mapKeys(data))
	} else {
		t.Logf("Got %d sherlocks", len(sherlocks))
		if len(sherlocks) > 0 {
			if s, ok := sherlocks[0].(map[string]interface{}); ok {
				t.Logf("First sherlock: name=%v, difficulty=%v, category=%v", s["name"], s["difficulty"], s["category_name"])
			}
		}
	}
}

func TestListSherlockCategories(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewListSherlockCategories(client)

	resp, err := tool.Execute(context.Background(), nil)
	if err != nil {
		t.Fatalf("list_sherlock_categories failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	t.Logf("Categories response: %s", resp.Content[0].Text[:min(500, len(resp.Content[0].Text))])
}

func TestGetSherlockTasks(t *testing.T) {
	client := setupClient(t)

	// Use Brutus (a well-known Sherlock, ID may vary — we'll list first)
	listTool := tools.NewListSherlocks(client)
	listResp, err := listTool.Execute(context.Background(), map[string]interface{}{
		"per_page": float64(1),
	})
	if err != nil {
		t.Skipf("could not list sherlocks: %v", err)
	}

	var sherlocks []map[string]interface{}
	if err := json.Unmarshal([]byte(listResp.Content[0].Text), &sherlocks); err != nil {
		t.Skipf("could not parse sherlocks list: %v", err)
	}
	if len(sherlocks) == 0 {
		t.Skip("no sherlocks available")
	}

	sherlockID := sherlocks[0]["id"].(float64)
	t.Logf("Testing with sherlock: %v (id=%v)", sherlocks[0]["name"], sherlockID)

	tasksTool := tools.NewGetSherlockTasks(client)
	tasksResp, err := tasksTool.Execute(context.Background(), map[string]interface{}{
		"sherlock_id": sherlockID,
	})
	if err != nil {
		t.Fatalf("get_sherlock_tasks failed: %v", err)
	}

	t.Logf("Tasks response: %s", tasksResp.Content[0].Text[:min(500, len(tasksResp.Content[0].Text))])
}

// ==================== Challenge instance tests ====================

func TestGetChallengeInfo(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewGetChallengeInfo(client)

	// Challenge 1042 = Magical Palindrome (Web, Very Easy, Docker+Download)
	resp, err := tool.Execute(context.Background(), map[string]interface{}{
		"challenge_id": float64(1042),
	})
	if err != nil {
		t.Fatalf("get_challenge_info failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	data := parseJSONContent(t, resp.Content[0].Text)
	if data == nil {
		t.Fatal("could not parse response")
	}

	if name, ok := data["name"].(string); ok {
		t.Logf("Challenge: %s, category=%v, difficulty=%v, docker=%v, download=%v",
			name, data["category_name"], data["difficulty"], data["docker"], data["download"])
	}

	// Verify docker field exists
	if _, ok := data["docker"]; !ok {
		t.Error("missing 'docker' field in challenge info")
	}
}

func TestDownloadChallenge(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewDownloadChallenge(client)

	// Challenge 1042 = Magical Palindrome (has downloadable files)
	resp, err := tool.Execute(context.Background(), map[string]interface{}{
		"challenge_id": float64(1042),
	})
	if err != nil {
		t.Fatalf("download_challenge failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	data := parseJSONContent(t, resp.Content[0].Text)
	if data == nil {
		t.Fatal("could not parse response")
	}

	if url, ok := data["url"].(string); ok {
		t.Logf("Download URL (first 100 chars): %s", url[:min(100, len(url))])
	} else {
		t.Error("missing 'url' field in download response")
	}

	if exp, ok := data["expires_in"].(float64); ok {
		t.Logf("Expires in: %v seconds", exp)
	}
}

func TestChallengeContainerLifecycle(t *testing.T) {
	if os.Getenv("RUN_LIFECYCLE_TESTS") != "1" {
		t.Skip("Skipping container lifecycle test (set RUN_LIFECYCLE_TESTS=1 to enable)")
	}

	client := setupClient(t)
	challengeID := float64(1042) // Magical Palindrome

	// Spawn
	spawnTool := tools.NewSpawnChallengeContainer(client)
	spawnResp, err := spawnTool.Execute(context.Background(), map[string]interface{}{
		"challenge_id": challengeID,
	})
	if err != nil {
		t.Fatalf("spawn_challenge_container failed: %v", err)
	}
	t.Logf("Spawn response: %s", spawnResp.Content[0].Text)

	// Wait a bit for container to come up
	t.Log("Waiting 5s for container...")
	time.Sleep(5 * time.Second)

	// Check challenge info to see docker_ip/docker_ports
	infoTool := tools.NewGetChallengeInfo(client)
	infoResp, err := infoTool.Execute(context.Background(), map[string]interface{}{
		"challenge_id": challengeID,
	})
	if err != nil {
		t.Errorf("get_challenge_info after spawn failed: %v", err)
	} else {
		data := parseJSONContent(t, infoResp.Content[0].Text)
		t.Logf("After spawn: docker_ip=%v, docker_ports=%v, docker_status=%v",
			data["docker_ip"], data["docker_ports"], data["docker_status"])
	}

	// Stop
	stopTool := tools.NewStopChallengeContainer(client)
	stopResp, err := stopTool.Execute(context.Background(), map[string]interface{}{
		"challenge_id": challengeID,
	})
	if err != nil {
		t.Fatalf("stop_challenge_container failed: %v", err)
	}
	t.Logf("Stop response: %s", stopResp.Content[0].Text)
}

// ==================== Platform tests ====================

func TestGetVPNStatus(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewGetVPNStatus(client)

	resp, err := tool.Execute(context.Background(), nil)
	if err != nil {
		t.Fatalf("get_vpn_status failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	t.Logf("VPN status: %s", resp.Content[0].Text[:min(500, len(resp.Content[0].Text))])
}

func TestGetActiveResources(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewGetActiveResources(client)

	resp, err := tool.Execute(context.Background(), nil)
	if err != nil {
		t.Fatalf("get_active_resources failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	data := parseJSONContent(t, resp.Content[0].Text)
	if data == nil {
		t.Fatal("could not parse response")
	}
	t.Logf("Active machine: %v", data["active_machine"])
	if connErr, ok := data["connections_error"]; ok {
		t.Logf("Connections endpoint error (v5 may need different base URL): %v", connErr)
	} else {
		t.Logf("Connections: %v", data["connections"])
	}
}

func TestListChallengeCategories(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewListChallengeCategories(client)

	resp, err := tool.Execute(context.Background(), nil)
	if err != nil {
		t.Fatalf("list_challenge_categories failed: %v", err)
	}
	if len(resp.Content) == 0 {
		t.Fatal("empty response content")
	}

	t.Logf("Categories: %s", resp.Content[0].Text[:min(800, len(resp.Content[0].Text))])
}

func TestGetRecommended(t *testing.T) {
	client := setupClient(t)
	tool := tools.NewGetRecommended(client)

	// Test machine recommendations
	resp, err := tool.Execute(context.Background(), map[string]interface{}{
		"type": "machines",
	})
	if err != nil {
		t.Fatalf("get_recommended (machines) failed: %v", err)
	}
	t.Logf("Recommended machines: %s", resp.Content[0].Text[:min(500, len(resp.Content[0].Text))])

	// Test challenge recommendations
	resp2, err := tool.Execute(context.Background(), map[string]interface{}{
		"type": "challenges",
	})
	if err != nil {
		t.Fatalf("get_recommended (challenges) failed: %v", err)
	}
	t.Logf("Recommended challenges: %s", resp2.Content[0].Text[:min(500, len(resp2.Content[0].Text))])
}

// ==================== Helpers ====================

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
