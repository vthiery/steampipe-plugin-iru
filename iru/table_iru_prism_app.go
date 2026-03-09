package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// PrismApp represents a single application record from the Prism app inventory endpoint.
// Field names are best-effort; exact shapes are confirmed once Prism is enabled on the token.
type PrismApp struct {
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	AppName    string `json:"app_name"`
	BundleID   string `json:"bundle_id"`
	Version    string `json:"version"`
	Path       string `json:"path"`
	Source     string `json:"source"`
	Process    string `json:"process"`
	UpdatedAt  string `json:"updated_at"`
}

// prismAppsResponse is the paginated response from GET /api/v1/prism/apps.
type prismAppsResponse struct {
	Next    *string    `json:"next"`
	Results []PrismApp `json:"results"`
}

//// TABLE DEFINITION

func tableIruPrismApp() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_prism_app",
		Description: "Cross-fleet application inventory from the Prism telemetry endpoint. Requires Prism to be enabled on the API token.",
		List: &plugin.ListConfig{
			Hydrate: listPrismApps,
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_STRING, Description: "The device ID."},
			{Name: "device_name", Type: proto.ColumnType_STRING, Description: "The device name."},
			{Name: "app_name", Type: proto.ColumnType_STRING, Description: "Application name."},
			{Name: "bundle_id", Type: proto.ColumnType_STRING, Description: "Bundle identifier."},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "Application version."},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Installation path."},
			{Name: "source", Type: proto.ColumnType_STRING, Description: "Application source (e.g. Apple, Identified Developer)."},
			{Name: "process", Type: proto.ColumnType_STRING, Description: "Process name."},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Description: "Last updated timestamp."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listPrismApps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{"limit": "300"}
	for {
		var resp prismAppsResponse
		if err := client.get(ctx, "/api/v1/prism/apps", params, &resp); err != nil {
			return nil, fmt.Errorf("listing prism apps: %w", err)
		}
		for _, item := range resp.Results {
			d.StreamListItem(ctx, item)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if resp.Next == nil || *resp.Next == "" {
			break
		}
		cursor, err := parseNextCursor(*resp.Next)
		if err != nil || cursor == "" {
			break
		}
		params["cursor"] = cursor
	}
	return nil, nil
}
