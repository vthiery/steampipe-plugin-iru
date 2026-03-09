package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// PrismFileVault represents a FileVault record from the Prism filevault endpoint.
type PrismFileVault struct {
	DeviceID              string `json:"device_id"`
	DeviceName            string `json:"device_name"`
	FileVaultEnabled      bool   `json:"filevault_enabled"`
	RecoveryKeyType       string `json:"recovery_key_type"`
	PRKEscrowed           bool   `json:"prk_escrowed"`
	NextRotation          string `json:"next_rotation"`
	RegenRequired         bool   `json:"regen_required"`
	UpdatedAt             string `json:"updated_at"`
}

// prismFileVaultResponse is the paginated response from GET /api/v1/prism/filevault.
type prismFileVaultResponse struct {
	Next    *string          `json:"next"`
	Results []PrismFileVault `json:"results"`
}

//// TABLE DEFINITION

func tableIruPrismFileVault() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_prism_filevault",
		Description: "Cross-fleet FileVault encryption status from the Prism telemetry endpoint. Requires Prism to be enabled on the API token.",
		List: &plugin.ListConfig{
			Hydrate: listPrismFileVault,
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_STRING, Description: "The device ID."},
			{Name: "device_name", Type: proto.ColumnType_STRING, Description: "The device name."},
			{Name: "filevault_enabled", Type: proto.ColumnType_BOOL, Description: "Whether FileVault is enabled."},
			{Name: "recovery_key_type", Type: proto.ColumnType_STRING, Description: "Type of recovery key (Personal, Institutional, etc.)."},
			{Name: "prk_escrowed", Type: proto.ColumnType_BOOL, Description: "Whether the personal recovery key is escrowed in Iru."},
			{Name: "next_rotation", Type: proto.ColumnType_STRING, Description: "Next scheduled key rotation timestamp."},
			{Name: "regen_required", Type: proto.ColumnType_BOOL, Description: "Whether key regeneration is required."},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Description: "Last updated timestamp."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listPrismFileVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{"limit": "300"}
	for {
		var resp prismFileVaultResponse
		if err := client.get(ctx, "/api/v1/prism/filevault", params, &resp); err != nil {
			return nil, fmt.Errorf("listing prism filevault: %w", err)
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
