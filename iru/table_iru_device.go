package iru

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// Device represents an endpoint managed by Iru (Kandji).
type Device struct {
	AgentInstalled           bool          `json:"agent_installed"`
	AssetTag                 string        `json:"asset_tag"`
	AssignedUser             *AssignedUser `json:"assigned_user"`
	BlueprintID              string        `json:"blueprint_id"`
	BlueprintName            string        `json:"blueprint_name"`
	DeviceID                 string        `json:"device_id"`
	DeviceName               string        `json:"device_name"`
	FirstEnrollment          string        `json:"first_enrollment"`
	IsActivationLocked       bool          `json:"is_activation_locked"`
	IsMissing                bool          `json:"is_missing"`
	IsRemoved                bool          `json:"is_removed"`
	LastCheckIn              string        `json:"last_check_in"`
	LastEnrollment           string        `json:"last_enrollment"`
	MDMEnabled               bool          `json:"mdm_enabled"`
	Model                    string        `json:"model"`
	OSVersion                string        `json:"os_version"`
	Platform                 string        `json:"platform"`
	SerialNumber             string        `json:"serial_number"`
	SupplementalBuildVersion string        `json:"supplemental_build_version"`
	SystemVersion            string        `json:"system_version"`
	Tags                     []string      `json:"tags"`
	User                     *AssignedUser `json:"-"`
}

// deviceJSON is the raw JSON shape of a device; "user" can be an object or "".
type deviceJSON struct {
	AgentInstalled           bool            `json:"agent_installed"`
	AssetTag                 string          `json:"asset_tag"`
	AssignedUser             *AssignedUser   `json:"assigned_user"`
	BlueprintID              string          `json:"blueprint_id"`
	BlueprintName            string          `json:"blueprint_name"`
	DeviceID                 string          `json:"device_id"`
	DeviceName               string          `json:"device_name"`
	FirstEnrollment          string          `json:"first_enrollment"`
	IsActivationLocked       bool            `json:"is_activation_locked"`
	IsMissing                bool            `json:"is_missing"`
	IsRemoved                bool            `json:"is_removed"`
	LastCheckIn              string          `json:"last_check_in"`
	LastEnrollment           string          `json:"last_enrollment"`
	MDMEnabled               bool            `json:"mdm_enabled"`
	Model                    string          `json:"model"`
	OSVersion                string          `json:"os_version"`
	Platform                 string          `json:"platform"`
	SerialNumber             string          `json:"serial_number"`
	SupplementalBuildVersion string          `json:"supplemental_build_version"`
	SystemVersion            string          `json:"system_version"`
	Tags                     []string        `json:"tags"`
	UserRaw                  json.RawMessage `json:"user"`
}

// UnmarshalJSON handles the polymorphic "user" field (object or empty string "").
func (d *Device) UnmarshalJSON(b []byte) error {
	var raw deviceJSON
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	d.AgentInstalled = raw.AgentInstalled
	d.AssetTag = raw.AssetTag
	d.AssignedUser = raw.AssignedUser
	d.BlueprintID = raw.BlueprintID
	d.BlueprintName = raw.BlueprintName
	d.DeviceID = raw.DeviceID
	d.DeviceName = raw.DeviceName
	d.FirstEnrollment = raw.FirstEnrollment
	d.IsActivationLocked = raw.IsActivationLocked
	d.IsMissing = raw.IsMissing
	d.IsRemoved = raw.IsRemoved
	d.LastCheckIn = raw.LastCheckIn
	d.LastEnrollment = raw.LastEnrollment
	d.MDMEnabled = raw.MDMEnabled
	d.Model = raw.Model
	d.OSVersion = raw.OSVersion
	d.Platform = raw.Platform
	d.SerialNumber = raw.SerialNumber
	d.SupplementalBuildVersion = raw.SupplementalBuildVersion
	d.SystemVersion = raw.SystemVersion
	d.Tags = raw.Tags
	// Only unmarshal User if it's a JSON object (not "" or null)
	if len(raw.UserRaw) > 0 && raw.UserRaw[0] == '{' {
		var u AssignedUser
		if err := json.Unmarshal(raw.UserRaw, &u); err == nil {
			d.User = &u
		}
	}
	return nil
}

// devicesResponse is a bare JSON array returned by GET /api/v1/devices.
type devicesResponse = []Device

//// TABLE DEFINITION

func tableIruDevice() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_device",
		Description: "List and inspect devices enrolled in your Iru tenant.",
		List: &plugin.ListConfig{
			Hydrate: listDevices,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "blueprint_id", Require: plugin.Optional},
				{Name: "platform", Require: plugin.Optional},
				{Name: "is_missing", Require: plugin.Optional},
				{Name: "is_removed", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("device_id"),
			Hydrate:    getDevice,
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("DeviceID"), Description: "Unique identifier for the device."},
			{Name: "device_name", Type: proto.ColumnType_STRING, Description: "Name of the device."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "Serial number of the device."},
			{Name: "model", Type: proto.ColumnType_STRING, Description: "Hardware model of the device (e.g. MacBook Pro)."},
			{Name: "platform", Type: proto.ColumnType_STRING, Description: "Device platform (Mac, iPhone, iPad, AppleTV)."},
			{Name: "os_version", Type: proto.ColumnType_STRING, Transform: transform.FromField("OSVersion"), Description: "Operating system version installed on the device."},
			{Name: "system_version", Type: proto.ColumnType_STRING, Description: "System version string reported by the device."},
			{Name: "supplemental_build_version", Type: proto.ColumnType_STRING, Description: "Supplemental OS build version."},
			{Name: "asset_tag", Type: proto.ColumnType_STRING, Description: "Asset tag assigned to the device."},
			{Name: "blueprint_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("BlueprintID"), Description: "ID of the blueprint the device is assigned to."},
			{Name: "blueprint_name", Type: proto.ColumnType_STRING, Description: "Name of the blueprint the device is assigned to."},
			{Name: "mdm_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("MDMEnabled"), Description: "Whether MDM is enabled on the device."},
			{Name: "agent_installed", Type: proto.ColumnType_BOOL, Description: "Whether the Iru agent is installed on the device."},
			{Name: "is_missing", Type: proto.ColumnType_BOOL, Description: "Whether the device is marked as missing."},
			{Name: "is_removed", Type: proto.ColumnType_BOOL, Description: "Whether the device has been removed from the tenant."},
			{Name: "is_activation_locked", Type: proto.ColumnType_BOOL, Description: "Whether Activation Lock is enabled on the device."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "Tags assigned to the device."},
			{Name: "first_enrollment", Type: proto.ColumnType_STRING, Description: "Time the device first enrolled in Iru."},
			{Name: "last_enrollment", Type: proto.ColumnType_STRING, Description: "Time the device last enrolled in Iru."},
			{Name: "last_check_in", Type: proto.ColumnType_STRING, Description: "Time the device last checked in."},

			// Assigned user (populated from the "user" field in the device list API response)
			{Name: "assigned_user_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("User.ID"), Description: "ID of the user assigned to this device."},
			{Name: "assigned_user_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("User.Name"), Description: "Name of the user assigned to this device."},
			{Name: "assigned_user_email", Type: proto.ColumnType_STRING, Transform: transform.FromField("User.Email"), Description: "Email of the user assigned to this device."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listDevices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	limitInt := 300
	limit := fmt.Sprintf("%d", limitInt)
	offset := 0

	for {
		params := map[string]string{
			"limit":  limit,
			"offset": fmt.Sprintf("%d", offset),
		}
		if v := d.EqualsQualString("blueprint_id"); v != "" {
			params["blueprint_id"] = v
		}
		if v := d.EqualsQualString("platform"); v != "" {
			params["platform"] = v
		}

		var result devicesResponse
		if err := client.get(ctx, "/api/v1/devices", params, &result); err != nil {
			return nil, fmt.Errorf("listing devices: %w", err)
		}

		for _, device := range result {
			d.StreamListItem(ctx, device)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if len(result) < limitInt {
			break
		}
		offset += len(result)
	}

	return nil, nil
}

func getDevice(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	deviceID := d.EqualsQualString("device_id")
	if deviceID == "" {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// GET /api/v1/devices/{device_id} returns device details as a JSON object.
	var result map[string]json.RawMessage
	if err := client.get(ctx, fmt.Sprintf("/api/v1/devices/%s", deviceID), nil, &result); err != nil {
		return nil, fmt.Errorf("getting device %s: %w", deviceID, err)
	}

	// The device details endpoint returns device info alongside grouped sub-objects.
	// Extract the flat "device" top-level fields and re-unmarshal into a Device.
	deviceRaw, ok := result["device"]
	if !ok {
		// If no "device" wrapper, the whole payload is the device.
		var dev Device
		full, _ := json.Marshal(result)
		if err := json.Unmarshal(full, &dev); err != nil {
			return nil, fmt.Errorf("parsing device response: %w", err)
		}
		return dev, nil
	}

	var dev Device
	if err := json.Unmarshal(deviceRaw, &dev); err != nil {
		return nil, fmt.Errorf("parsing device response: %w", err)
	}
	return dev, nil
}
