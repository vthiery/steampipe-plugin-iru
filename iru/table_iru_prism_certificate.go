package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// PrismCertificate represents a certificate record from the Prism certificates endpoint.
type PrismCertificate struct {
	DeviceID    string `json:"device_id"`
	DeviceName  string `json:"device_name"`
	CommonName  string `json:"common_name"`
	Issuer      string `json:"issuer"`
	Subject     string `json:"subject"`
	ExpiryDate  string `json:"expiry_date"`
	IsExpired   bool   `json:"is_expired"`
	KeyUsage    string `json:"key_usage"`
	SerialNum   string `json:"serial_number"`
	UpdatedAt   string `json:"updated_at"`
}

// prismCertificatesResponse is the paginated response from GET /api/v1/prism/certificates.
type prismCertificatesResponse struct {
	Next    *string            `json:"next"`
	Results []PrismCertificate `json:"results"`
}

//// TABLE DEFINITION

func tableIruPrismCertificate() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_prism_certificate",
		Description: "Cross-fleet certificate inventory from the Prism telemetry endpoint. Requires Prism to be enabled on the API token.",
		List: &plugin.ListConfig{
			Hydrate: listPrismCertificates,
		},
		Columns: []*plugin.Column{
			{Name: "device_id", Type: proto.ColumnType_STRING, Description: "The device ID."},
			{Name: "device_name", Type: proto.ColumnType_STRING, Description: "The device name."},
			{Name: "common_name", Type: proto.ColumnType_STRING, Description: "Certificate common name (CN)."},
			{Name: "issuer", Type: proto.ColumnType_STRING, Description: "Certificate issuer."},
			{Name: "subject", Type: proto.ColumnType_STRING, Description: "Certificate subject."},
			{Name: "expiry_date", Type: proto.ColumnType_STRING, Description: "Certificate expiry date."},
			{Name: "is_expired", Type: proto.ColumnType_BOOL, Description: "Whether the certificate is expired."},
			{Name: "key_usage", Type: proto.ColumnType_STRING, Description: "Key usage flags."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "Certificate serial number."},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Description: "Last updated timestamp."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listPrismCertificates(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{"limit": "300"}
	for {
		var resp prismCertificatesResponse
		if err := client.get(ctx, "/api/v1/prism/certificates", params, &resp); err != nil {
			return nil, fmt.Errorf("listing prism certificates: %w", err)
		}
		for _, item := range resp.Results {
			d.StreamListItem(ctx, item)
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
