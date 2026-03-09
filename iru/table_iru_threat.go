package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// Threat represents a threat event recorded by Iru's Endpoint Detection & Response.
type Threat struct {
	BehavioralDetectionID string `json:"behavioral_detection_id"`
	Classification        string `json:"classification"`
	DateTime              string `json:"date_time"`
	DeviceID              string `json:"device_id"`
	DeviceName            string `json:"device_name"`
	FilePath              string `json:"file_path"`
	FileName              string `json:"file_name"`
	InitiatingProcess     string `json:"initiating_process"`
	MalwareFamily         string `json:"malware_family"`
	ProcessOwner          string `json:"process_owner"`
	SerialNumber          string `json:"serial_number"`
	SHA256                string `json:"sha256"`
	Status                string `json:"status"`
	ThreatID              string `json:"threat_id"`
}

// threatsResponse is the cursor-paginated envelope returned by GET /api/v1/threat-details.
type threatsResponse struct {
	Next     *string  `json:"next"`
	Previous *string  `json:"previous"`
	Results  []Threat `json:"results"`
}

//// TABLE DEFINITION

func tableIruThreat() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_threat",
		Description: "List threat events detected on devices in your Iru tenant.",
		List: &plugin.ListConfig{
			Hydrate: listThreats,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "device_id", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "classification", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "threat_id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the threat event."},
			{Name: "device_id", Type: proto.ColumnType_STRING, Description: "ID of the device on which the threat was detected."},
			{Name: "device_name", Type: proto.ColumnType_STRING, Description: "Name of the device on which the threat was detected."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "Serial number of the affected device."},
			{Name: "file_name", Type: proto.ColumnType_STRING, Description: "Name of the malicious or suspicious file."},
			{Name: "file_path", Type: proto.ColumnType_STRING, Description: "Full path to the malicious or suspicious file."},
			{Name: "sha256", Type: proto.ColumnType_STRING, Description: "SHA-256 hash of the detected file."},
			{Name: "classification", Type: proto.ColumnType_STRING, Description: "Threat classification (malware, pup, etc.)."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Current status of the threat (blocked, quarantined, resolved, etc.)."},
			{Name: "malware_family", Type: proto.ColumnType_STRING, Description: "Malware family name, if identified."},
			{Name: "behavioral_detection_id", Type: proto.ColumnType_STRING, Description: "ID of the behavioral detection rule that triggered this event."},
			{Name: "initiating_process", Type: proto.ColumnType_STRING, Description: "Name of the process that initiated the threat."},
			{Name: "process_owner", Type: proto.ColumnType_STRING, Description: "Owner of the process that triggered the detection."},
			{Name: "date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time the threat was detected."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listThreats(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

params := map[string]string{
		"limit": "300",
	}
	if v := d.EqualsQualString("device_id"); v != "" {
		params["device_id"] = v
	}
	if v := d.EqualsQualString("status"); v != "" {
		params["status"] = v
	}
	if v := d.EqualsQualString("classification"); v != "" {
		params["classification"] = v
	}

	for {
		var result threatsResponse
		if err := client.get(ctx, "/api/v1/threat-details", params, &result); err != nil {
			return nil, fmt.Errorf("listing threats: %w", err)
		}

		for _, threat := range result.Results {
			d.StreamListItem(ctx, threat)
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
