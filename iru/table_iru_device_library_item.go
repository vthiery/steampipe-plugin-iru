package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// DeviceLibraryItem is one library item assignment record for a device.
type DeviceLibraryItem struct {
	DeviceID          string `json:"device_id"`
	Name              string `json:"name"`
	Status            string `json:"status"`
	Type              string `json:"type"`
	ReportedAt        string `json:"reported_at"`
	LastAuditRun      string `json:"last_audit_run"`
	LastAuditLog      string `json:"last_audit_log"`
	ControlReportedAt string `json:"control_reported_at"`
	ControlLog        string `json:"control_log"`
	Log               string `json:"log"`
}

// deviceLibraryItemsResponse is the shape returned by GET /api/v1/devices/{id}/library-items.
type deviceLibraryItemsResponse struct {
	LibraryItems []DeviceLibraryItem `json:"library_items"`
}

//// TABLE DEFINITION

func tableIruDeviceLibraryItem() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_device_library_item",
		Description: "Lists the per-device library item assignment status for all managed devices.",
		List: &plugin.ListConfig{
			Hydrate: listDeviceLibraryItems,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "device_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_STRING, Description: "The device ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Library item name."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Compliance/installation status."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Library item type (e.g. profile, script, app)."},
			{Name: "reported_at", Type: proto.ColumnType_STRING, Description: "Timestamp of the last status report."},
			{Name: "last_audit_run", Type: proto.ColumnType_STRING, Description: "Timestamp of the last audit run."},
			{Name: "last_audit_log", Type: proto.ColumnType_STRING, Description: "Log output from the last audit run."},
			{Name: "control_reported_at", Type: proto.ColumnType_STRING, Description: "Timestamp of the last control report."},
			{Name: "control_log", Type: proto.ColumnType_STRING, Description: "Control enforcement log output."},
			{Name: "log", Type: proto.ColumnType_STRING, Description: "Full log output."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listDeviceLibraryItems(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	if deviceID := d.EqualsQualString("device_id"); deviceID != "" {
		return nil, fetchLibraryItemsForDevice(ctx, client, d, deviceID)
	}

	deviceIDs, err := getAllDeviceIDs(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("listing device IDs: %w", err)
	}

	for _, deviceID := range deviceIDs {
		if err := fetchLibraryItemsForDevice(ctx, client, d, deviceID); err != nil {
			plugin.Logger(ctx).Warn("iru_device_library_item: failed to fetch", "device_id", deviceID, "error", err)
			continue
		}
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

func fetchLibraryItemsForDevice(ctx context.Context, client *Client, d *plugin.QueryData, deviceID string) error {
	var resp deviceLibraryItemsResponse
	if err := client.get(ctx, fmt.Sprintf("/api/v1/devices/%s/library-items", deviceID), nil, &resp); err != nil {
		return err
	}
	for _, item := range resp.LibraryItems {
		item.DeviceID = deviceID
		d.StreamListItem(ctx, item)
		if d.RowsRemaining(ctx) == 0 {
			return nil
		}
	}
	return nil
}
