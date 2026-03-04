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
	var resp PaginatedResponse[McpServer]
	if err := s.client.request(ctx, http.MethodGet, "/mcp/servers"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("mcp.ListServers: %w", err)
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
func (s *McpService) Proxy(ctx context.Context, input *McpProxyInput) (*McpProxyResponse, error) {
	var resp McpProxyResponse
	if err := s.client.request(ctx, http.MethodPost, "/mcp/proxy", input, &resp); err != nil {
		return nil, fmt.Errorf("mcp.Proxy: %w", err)
	}
	return &resp, nil
}
