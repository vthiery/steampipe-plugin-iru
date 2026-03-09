package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// LibraryItem represents a library item in an Iru tenant.
type LibraryItem struct {
	Active      bool     `json:"active"`
	BlueprintIDs []string `json:"blueprint_ids"`
	CreatedAt   string   `json:"created_at"`
	ID          string   `json:"id"`
	Info        string   `json:"info"`
	IsManaged   bool     `json:"is_managed"`
	ItemType    string   `json:"item_type"`
	Name        string   `json:"name"`
	RunsOn      []string `json:"runs_on"`
	UpdatedAt   string   `json:"updated_at"`
}

// libraryItemsResponse is the cursor-paginated envelope returned by GET /api/v1/library/custom-apps.
type libraryItemsResponse struct {
	Next     *string       `json:"next"`
	Previous *string       `json:"previous"`
	Results  []LibraryItem `json:"results"`
}

//// TABLE DEFINITION

func tableIruLibraryItem() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_library_item",
		Description: "List library items configured in your Iru tenant.",
		List: &plugin.ListConfig{
			Hydrate: listLibraryItems,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "item_type", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getLibraryItem,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the library item."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the library item."},
			{Name: "item_type", Type: proto.ColumnType_STRING, Description: "Type of library item (profile, auto-app, custom-app, etc.)."},
			{Name: "info", Type: proto.ColumnType_STRING, Description: "Additional information or description for the library item."},
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Whether the library item is active."},
			{Name: "is_managed", Type: proto.ColumnType_BOOL, Description: "Whether the library item is managed (read-only) by Iru."},
			{Name: "runs_on", Type: proto.ColumnType_JSON, Description: "Device platforms this library item applies to."},
			{Name: "blueprint_ids", Type: proto.ColumnType_JSON, Description: "IDs of blueprints that include this library item."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the library item was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the library item was last updated."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listLibraryItems(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"limit": "300",
	}
	if v := d.EqualsQualString("item_type"); v != "" {
		params["item_type"] = v
	}

	for {
		var result libraryItemsResponse
		if err := client.get(ctx, "/api/v1/library/custom-apps", params, &result); err != nil {
			return nil, fmt.Errorf("listing library items: %w", err)
		}

		for _, item := range result.Results {
			d.StreamListItem(ctx, item)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.Next == nil || len(result.Results) == 0 {
			break
		}
		cursor, err := parseNextCursor(*result.Next)
		if err != nil || cursor == "" {
			break
		}
		params["cursor"] = cursor
	}

	return nil, nil
}

func getLibraryItem(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result LibraryItem
		if err := client.get(ctx, fmt.Sprintf("/api/v1/library/custom-apps/%s", id), nil, &result); err != nil {
		return nil, fmt.Errorf("getting library item %s: %w", id, err)
	}

	return result, nil
}
