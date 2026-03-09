package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// DeviceApp represents a single application installed on a managed device.
type DeviceApp struct {
	DeviceID         string `json:"device_id"`
	AppID            string `json:"app_id"`
	AppName          string `json:"app_name"`
	BundleID         string `json:"bundle_id"`
	BundleSize       string `json:"bundle_size"`
	AppStoreVendable string `json:"app_store_vendable"`
	DeviceBasedVPP   string `json:"device_based_vpp"`
	Source           string `json:"source"`
	Process          string `json:"process"`
	Signature        string `json:"signature"`
	CreationDate     string `json:"creation_date"`
	ModificationDate string `json:"modification_date"`
	Path             string `json:"path"`
	Version          string `json:"version"`
}

// deviceAppsResponse is the shape returned by GET /api/v1/devices/{id}/apps.
type deviceAppsResponse struct {
	DeviceID string      `json:"device_id"`
	Apps     []DeviceApp `json:"apps"`
}

//// TABLE DEFINITION

func tableIruDeviceApp() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_device_app",
		Description: "Lists all applications installed on each managed device.",
		List: &plugin.ListConfig{
			Hydrate: listDeviceApps,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "device_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_STRING, Description: "The device ID."},
			{Name: "app_id", Type: proto.ColumnType_STRING, Description: "The app record ID."},
			{Name: "app_name", Type: proto.ColumnType_STRING, Description: "Application name."},
			{Name: "bundle_id", Type: proto.ColumnType_STRING, Description: "Bundle identifier."},
			{Name: "bundle_size", Type: proto.ColumnType_STRING, Description: "Bundle size in bytes."},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "Application version."},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Filesystem installation path."},
			{Name: "source", Type: proto.ColumnType_STRING, Description: "Application source (e.g. Apple, Identified Developer)."},
			{Name: "process", Type: proto.ColumnType_STRING, Description: "Process name."},
			{Name: "creation_date", Type: proto.ColumnType_STRING, Description: "Creation date of the app record."},
			{Name: "modification_date", Type: proto.ColumnType_STRING, Description: "Last modification date."},
			{Name: "app_store_vendable", Type: proto.ColumnType_STRING, Description: "Whether the app is App Store vendable."},
			{Name: "device_based_vpp", Type: proto.ColumnType_STRING, Description: "Device-based VPP status."},
			{Name: "signature", Type: proto.ColumnType_STRING, Description: "Code signature information."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listDeviceApps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// If a specific device_id is requested, only fetch apps for that device.
	if deviceID := d.EqualsQualString("device_id"); deviceID != "" {
		if err := fetchAppsForDevice(ctx, client, d, deviceID); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// Otherwise, fetch all device IDs first, then apps for each.
	deviceIDs, err := getAllDeviceIDs(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("listing device IDs: %w", err)
	}

	for _, deviceID := range deviceIDs {
		if err := fetchAppsForDevice(ctx, client, d, deviceID); err != nil {
			// Log and skip rather than aborting the whole list.
			plugin.Logger(ctx).Warn("iru_device_app: failed to fetch apps", "device_id", deviceID, "error", err)
			continue
		}
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

// fetchAppsForDevice calls /api/v1/devices/{id}/apps and streams each app as a row.
func fetchAppsForDevice(ctx context.Context, client *Client, d *plugin.QueryData, deviceID string) error {
	var resp deviceAppsResponse
	if err := client.get(ctx, fmt.Sprintf("/api/v1/devices/%s/apps", deviceID), nil, &resp); err != nil {
		return err
	}
	for _, app := range resp.Apps {
		app.DeviceID = deviceID
		d.StreamListItem(ctx, app)
		if d.RowsRemaining(ctx) == 0 {
			return nil
		}
	}
	return nil
}

// getAllDeviceIDs fetches device IDs from the device list using offset pagination.
func getAllDeviceIDs(ctx context.Context, client *Client) ([]string, error) {
	limitInt := 300
	offset := 0
	var ids []string
	for {
		params := map[string]string{
			"limit":  fmt.Sprintf("%d", limitInt),
			"offset": fmt.Sprintf("%d", offset),
		}
		var batch []Device
		if err := client.get(ctx, "/api/v1/devices", params, &batch); err != nil {
			return nil, err
		}
		for _, dev := range batch {
			ids = append(ids, dev.DeviceID)
		}
		if len(batch) < limitInt {
			break
		}
		offset += len(batch)
	}
	return ids, nil
}
