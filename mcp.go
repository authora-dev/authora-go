package authora

import (
	"context"
	"fmt"
	"net/http"
)

// McpService handles MCP (Model Context Protocol) related API endpoints.
type McpService struct {
	client *httpClient
}

// RegisterServer registers a new MCP server. POST /mcp/servers
func (s *McpService) RegisterServer(ctx context.Context, input *RegisterMcpServerInput) (*McpServer, error) {
	var resp McpServer
	if err := s.client.request(ctx, http.MethodPost, "/mcp/servers", input, &resp); err != nil {
		return nil, fmt.Errorf("mcp.RegisterServer: %w", err)
	}
	return &resp, nil
}

// ListServers returns registered MCP servers. GET /mcp/servers
func (s *McpService) ListServers(ctx context.Context, input *ListMcpServersInput) (*PaginatedResponse[McpServer], error) {
	q := queryString(map[string]interface{}{
		"workspaceId": input.WorkspaceID,
		"page":        input.Page,
		"limit":       input.Limit,
	})
	// The API may return a paginated response or a raw array.
	// Try paginated first; if that fails, try raw array.
	var resp PaginatedResponse[McpServer]
	if err := s.client.request(ctx, http.MethodGet, "/mcp/servers"+q, nil, &resp); err != nil {
		// Fallback: try as raw array
		var items []McpServer
		if err2 := s.client.request(ctx, http.MethodGet, "/mcp/servers"+q, nil, &items); err2 != nil {
			return nil, fmt.Errorf("mcp.ListServers: %w", err)
		}
		return &PaginatedResponse[McpServer]{Items: items, Total: len(items)}, nil
	}
	return &resp, nil
}

// GetServer retrieves a single MCP server by ID. GET /mcp/servers/:serverId
func (s *McpService) GetServer(ctx context.Context, serverID string) (*McpServer, error) {
	var resp McpServer
	if err := s.client.request(ctx, http.MethodGet, "/mcp/servers/"+serverID, nil, &resp); err != nil {
		return nil, fmt.Errorf("mcp.GetServer: %w", err)
	}
	return &resp, nil
}

// UpdateServer modifies an existing MCP server.
// PATCH /mcp/servers/:serverId
func (s *McpService) UpdateServer(ctx context.Context, serverID string, input *UpdateMcpServerInput) (*McpServer, error) {
	var resp McpServer
	if err := s.client.request(ctx, http.MethodPatch, "/mcp/servers/"+serverID, input, &resp); err != nil {
		return nil, fmt.Errorf("mcp.UpdateServer: %w", err)
	}
	return &resp, nil
}

// ListTools returns tools registered on an MCP server.
// GET /mcp/servers/:serverId/tools
func (s *McpService) ListTools(ctx context.Context, serverID string) ([]McpTool, error) {
	var resp []McpTool
	if err := s.client.request(ctx, http.MethodGet, "/mcp/servers/"+serverID+"/tools", nil, &resp); err != nil {
		return nil, fmt.Errorf("mcp.ListTools: %w", err)
	}
	return resp, nil
}

// RegisterTool registers a tool on an MCP server.
// POST /mcp/servers/:serverId/tools
func (s *McpService) RegisterTool(ctx context.Context, serverID string, input *RegisterMcpToolInput) (*McpTool, error) {
	var resp McpTool
	if err := s.client.request(ctx, http.MethodPost, "/mcp/servers/"+serverID+"/tools", input, &resp); err != nil {
		return nil, fmt.Errorf("mcp.RegisterTool: %w", err)
	}
	return &resp, nil
}

// Proxy forwards a request to an MCP server. POST /mcp/proxy
// Builds a JSON-RPC 2.0 request with _authora metadata.
func (s *McpService) Proxy(ctx context.Context, input *McpProxyInput) (*McpProxyResponse, error) {
	params := make(map[string]interface{})
	for k, v := range input.Params {
		params[k] = v
	}
	// Ensure _authora metadata includes the server ID
	authora, ok := params["_authora"].(map[string]interface{})
	if !ok {
		authora = make(map[string]interface{})
	}
	if _, exists := authora["mcpServerId"]; !exists {
		authora["mcpServerId"] = input.ServerID
	}
	params["_authora"] = authora

	rpcBody := &mcpProxyJsonRpc{
		Jsonrpc: "2.0",
		Method:  input.Method,
		ID:      1,
		Params:  params,
	}

	var resp McpProxyResponse
	if err := s.client.request(ctx, http.MethodPost, "/mcp/proxy", rpcBody, &resp); err != nil {
		return nil, fmt.Errorf("mcp.Proxy: %w", err)
	}
	return &resp, nil
}
