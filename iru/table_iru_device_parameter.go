package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// DeviceParameter represents one parameter (MDM setting/compliance check) on a device.
type DeviceParameter struct {
	DeviceID    string `json:"device_id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
}

// deviceStatusResponse is the shape returned by GET /api/v1/devices/{id}/status.
// It contains both library_items and parameters. We only use parameters here;
// iru_device_library_item covers the library_items section via the dedicated endpoint.
type deviceStatusResponse struct {
	Parameters []DeviceParameter `json:"parameters"`
}

//// TABLE DEFINITION

func tableIruDeviceParameter() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_device_parameter",
		Description: "Lists per-device MDM parameter compliance status for all managed devices.",
		List: &plugin.ListConfig{
			Hydrate: listDeviceParameters,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "device_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_STRING, Description: "The device ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Parameter name (e.g. Set Computer Name)."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Compliance status (pass, fail, error, etc.)."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Top-level parameter category."},
			{Name: "subcategory", Type: proto.ColumnType_STRING, Description: "Parameter subcategory."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listDeviceParameters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	if deviceID := d.EqualsQualString("device_id"); deviceID != "" {
		return nil, fetchParametersForDevice(ctx, client, d, deviceID)
	}

	deviceIDs, err := getAllDeviceIDs(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("listing device IDs: %w", err)
	}

	for _, deviceID := range deviceIDs {
		if err := fetchParametersForDevice(ctx, client, d, deviceID); err != nil {
			plugin.Logger(ctx).Warn("iru_device_parameter: failed to fetch", "device_id", deviceID, "error", err)
			continue
		}
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

func fetchParametersForDevice(ctx context.Context, client *Client, d *plugin.QueryData, deviceID string) error {
	var resp deviceStatusResponse
	if err := client.get(ctx, fmt.Sprintf("/api/v1/devices/%s/status", deviceID), nil, &resp); err != nil {
		return err
	}
	for _, param := range resp.Parameters {
		param.DeviceID = deviceID
		d.StreamListItem(ctx, param)
		if d.RowsRemaining(ctx) == 0 {
			return nil
		}
	}
	return nil
}
