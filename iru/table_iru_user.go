package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)


//// TYPES

// User represents a user in an Iru tenant.
type User struct {
	Active           bool   `json:"active"`
	Archived         bool   `json:"archived"`
	CreatedAt        string `json:"created_at"`
	Department       string `json:"department"`
	DeprecatedUserID string `json:"deprecated_user_id"`
	Email            string `json:"email"`
	ID               string `json:"id"`
	JobTitle         string `json:"job_title"`
	Name             string `json:"name"`
	DeviceCount      *int   `json:"device_count"`
	UpdatedAt        string `json:"updated_at"`
}

// usersResponse is the cursor-paginated envelope returned by GET /api/v1/users.
type usersResponse struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []User  `json:"results"`
}

//// TABLE DEFINITION

func tableIruUser() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_user",
		Description: "List users in your Iru tenant.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getUser,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Unique identifier for the user."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Full name of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Email address of the user."},
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Whether the user account is active."},
			{Name: "archived", Type: proto.ColumnType_BOOL, Description: "Whether the user is archived."},
			{Name: "department", Type: proto.ColumnType_STRING, Description: "Department the user belongs to."},
			{Name: "job_title", Type: proto.ColumnType_STRING, Description: "Job title of the user."},
					{Name: "device_count", Type: proto.ColumnType_INT, Transform: transform.FromField("DeviceCount"), Description: "Number of devices assigned to this user."},
			{Name: "deprecated_user_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("DeprecatedUserID"), Description: "Legacy numeric user ID."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the user was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the user was last updated."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"limit": "300",
	}

	for {
		var result usersResponse
		if err := client.get(ctx, "/api/v1/users", params, &result); err != nil {
			return nil, fmt.Errorf("listing users: %w", err)
		}

		for _, user := range result.Results {
			d.StreamListItem(ctx, user)
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

func getUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result User
	if err := client.get(ctx, fmt.Sprintf("/api/v1/users/%s", id), nil, &result); err != nil {
		return nil, fmt.Errorf("getting user %s: %w", id, err)
	}

	return result, nil
}
