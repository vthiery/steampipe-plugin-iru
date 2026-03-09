package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// Blueprint represents a blueprint in an Iru tenant.
type Blueprint struct {
	ActiveLibraryItems      int      `json:"active_library_items"`
	ActiveDevices           int      `json:"active_devices"`
	ComputerCount           int      `json:"computer_count"`
	Description             string   `json:"description"`
	EnrollmentCodeEnabled   bool     `json:"enrollment_code_is_enabled"`
	ID                      string   `json:"id"`
	IsDEPEnrollmentEnabled  bool     `json:"is_dep_enrollment_enabled"`
	IsDefaultBlueprint      bool     `json:"is_default_blueprint"`
	IsUserPortalEnabled     bool     `json:"is_user_portal_enabled"`
	Name                    string   `json:"name"`
	Params                  string   `json:"params"`
	SourceBlueprint         *string  `json:"source_blueprint"`
	Tags                    []string `json:"tags"`
}

// blueprintsResponse is the envelope returned by GET /api/v1/blueprints.
type blueprintsResponse struct {
	Count   int         `json:"count"`
	Next    *string     `json:"next"`
	Results []Blueprint `json:"results"`
}

//// TABLE DEFINITION

func tableIruBlueprint() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_blueprint",
		Description: "List blueprints configured in your Iru tenant.",
		List: &plugin.ListConfig{
			Hydrate: listBlueprints,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getBlueprint,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the blueprint."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the blueprint."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the blueprint."},
			{Name: "is_default_blueprint", Type: proto.ColumnType_BOOL, Description: "Whether this is the default blueprint for new device enrollments."},
			{Name: "is_dep_enrollment_enabled", Type: proto.ColumnType_BOOL, Description: "Whether ADE (Automated Device Enrollment) is enabled for this blueprint."},
			{Name: "is_user_portal_enabled", Type: proto.ColumnType_BOOL, Description: "Whether the self-service user portal is enabled for this blueprint."},
			{Name: "enrollment_code_is_enabled", Type: proto.ColumnType_BOOL, Description: "Whether enrollment code enrollment is enabled for this blueprint."},
			{Name: "active_library_items", Type: proto.ColumnType_INT, Description: "Number of active library items assigned to this blueprint."},
			{Name: "active_devices", Type: proto.ColumnType_INT, Description: "Number of active devices assigned to this blueprint."},
			{Name: "computer_count", Type: proto.ColumnType_INT, Description: "Total number of computers (Mac devices) in this blueprint."},
			{Name: "source_blueprint", Type: proto.ColumnType_STRING, Description: "ID of the blueprint this one was cloned from, if any."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "Tags associated with this blueprint."},
			{Name: "params", Type: proto.ColumnType_STRING, Description: "Additional blueprint parameters."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listBlueprints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	limit := "300"
	offset := 0

	for {
		params := map[string]string{
			"limit":  limit,
			"offset": fmt.Sprintf("%d", offset),
		}

		var result blueprintsResponse
		if err := client.get(ctx, "/api/v1/blueprints", params, &result); err != nil {
			return nil, fmt.Errorf("listing blueprints: %w", err)
		}

		for _, bp := range result.Results {
			d.StreamListItem(ctx, bp)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.Next == nil || len(result.Results) == 0 {
			break
		}
		offset += len(result.Results)
	}

	return nil, nil
}

func getBlueprint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result Blueprint
	if err := client.get(ctx, fmt.Sprintf("/api/v1/blueprints/%s", id), nil, &result); err != nil {
		return nil, fmt.Errorf("getting blueprint %s: %w", id, err)
	}

	return result, nil
}
