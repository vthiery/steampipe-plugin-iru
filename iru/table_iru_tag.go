package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// Tag represents a tag in an Iru tenant.
type Tag struct {
	Color string `json:"color"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

// tagsResponse is the envelope returned by GET /api/v1/tags.
type tagsResponse struct {
	Count   int     `json:"count"`
	Next    *string `json:"next"`
	Results []Tag   `json:"results"`
}

//// TABLE DEFINITION

func tableIruTag() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_tag",
		Description: "List tags configured in your Iru tenant.",
		List: &plugin.ListConfig{
			Hydrate: listTags,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the tag."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the tag."},
			{Name: "color", Type: proto.ColumnType_STRING, Description: "Colour associated with the tag."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

		var result tagsResponse
		if err := client.get(ctx, "/api/v1/tags", params, &result); err != nil {
			return nil, fmt.Errorf("listing tags: %w", err)
		}

		for _, tag := range result.Results {
			d.StreamListItem(ctx, tag)
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
