package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// AuditEvent represents a single entry in the Iru audit log.
type AuditEvent struct {
	Action       string      `json:"action"`
	ActionUID    string      `json:"action_uid"`
	ActorID      string      `json:"actor_id"`
	ActorName    string      `json:"actor_name"`
	CreatedAt    string      `json:"created_at"`
	EventType    string      `json:"event_type"`
	ID           string      `json:"id"`
	IPAddress    string      `json:"ip_address"`
	ResourceID   string      `json:"resource_id"`
	ResourceName string      `json:"resource_name"`
	ResourceType string      `json:"resource_type"`
	UID          string      `json:"uid"`
}

// auditLogResponse is the cursor-paginated envelope returned by GET /api/v1/audit/events.
type auditLogResponse struct {
	Next     *string      `json:"next"`
	Previous *string      `json:"previous"`
	Results  []AuditEvent `json:"results"`
}

//// TABLE DEFINITION

func tableIruAuditLog() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_audit_log",
		Description: "List audit log events from your Iru tenant.",
		List: &plugin.ListConfig{
			Hydrate: listAuditLog,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "actor_id", Require: plugin.Optional},
				{Name: "event_type", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the audit log entry."},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "External UID for the audit log entry."},
			{Name: "action", Type: proto.ColumnType_STRING, Description: "Human-readable description of the action taken."},
			{Name: "action_uid", Type: proto.ColumnType_STRING, Description: "Machine-readable action identifier."},
			{Name: "event_type", Type: proto.ColumnType_STRING, Description: "Type of event (e.g. blueprint.updated, device.enrolled)."},
			{Name: "actor_id", Type: proto.ColumnType_STRING, Description: "ID of the user or service that performed the action."},
			{Name: "actor_name", Type: proto.ColumnType_STRING, Description: "Name of the user or service that performed the action."},
			{Name: "ip_address", Type: proto.ColumnType_STRING, Description: "IP address from which the action was performed."},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Description: "ID of the resource that was acted upon."},
			{Name: "resource_name", Type: proto.ColumnType_STRING, Description: "Name of the resource that was acted upon."},
			{Name: "resource_type", Type: proto.ColumnType_STRING, Description: "Type of resource that was acted upon (device, blueprint, user, etc.)."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the audit event was recorded."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listAuditLog(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"limit": "300",
	}

	for {
		if v := d.EqualsQualString("actor_id"); v != "" {
			params["actor_id"] = v
		}
		if v := d.EqualsQualString("event_type"); v != "" {
			params["event_type"] = v
		}

		var result auditLogResponse
		if err := client.get(ctx, "/api/v1/audit/events", params, &result); err != nil {
			return nil, fmt.Errorf("listing audit log: %w", err)
		}

		for _, event := range result.Results {
			d.StreamListItem(ctx, event)
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
