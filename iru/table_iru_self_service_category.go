package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// SelfServiceCategory represents a category in the Iru Self Service catalog.
type SelfServiceCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// selfServiceCategoriesResponse is the paginated response from GET /api/v1/self-service/categories.
type selfServiceCategoriesResponse struct {
	Next    *string               `json:"next"`
	Results []SelfServiceCategory `json:"results"`
}

//// TABLE DEFINITION

func tableIruSelfServiceCategory() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_self_service_category",
		Description: "Lists categories defined in the Iru Self Service catalog.",
		List: &plugin.ListConfig{
			Hydrate: listSelfServiceCategories,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getSelfServiceCategory,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the category."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Category display name."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Category description."},
			{Name: "sort_order", Type: proto.ColumnType_INT, Description: "Display sort order."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listSelfServiceCategories(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{"limit": "300"}
	for {
		var resp selfServiceCategoriesResponse
		if err := client.get(ctx, "/api/v1/self-service/categories", params, &resp); err != nil {
			return nil, fmt.Errorf("listing self-service categories: %w", err)
		}
		for _, cat := range resp.Results {
			d.StreamListItem(ctx, cat)
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

func getSelfServiceCategory(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}
	var cat SelfServiceCategory
	if err := client.get(ctx, fmt.Sprintf("/api/v1/self-service/categories/%s", id), nil, &cat); err != nil {
		return nil, fmt.Errorf("getting self-service category %s: %w", id, err)
	}
	return cat, nil
}
