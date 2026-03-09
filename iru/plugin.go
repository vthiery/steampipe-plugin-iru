// Package iru provides a Steampipe plugin for querying Iru (Kandji) resources using SQL.
package iru

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin returns the definition of the Iru Steampipe plugin.
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-iru",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			// Device inventory
			"iru_device":                tableIruDevice(),
			"iru_device_detail":         tableIruDeviceDetail(),
			"iru_device_app":            tableIruDeviceApp(),
			"iru_device_library_item":   tableIruDeviceLibraryItem(),
			"iru_device_parameter":      tableIruDeviceParameter(),
			// Users, blueprints, tags
			"iru_user":                  tableIruUser(),
			"iru_blueprint":             tableIruBlueprint(),
			"iru_tag":                   tableIruTag(),
			// Library items, threats, audit
			"iru_library_item":          tableIruLibraryItem(),
			"iru_threat":                tableIruThreat(),
			"iru_audit_log":             tableIruAuditLog(),
			// Prism telemetry (requires Prism-enabled API token)
			"iru_prism_app":             tableIruPrismApp(),
			"iru_prism_certificate":     tableIruPrismCertificate(),
			"iru_prism_filevault":       tableIruPrismFileVault(),
			// Self Service
			"iru_self_service_category": tableIruSelfServiceCategory(),
			// Vulnerability Management (requires Vulnerability Management add-on)
			"iru_vulnerability":        tableIruVulnerabilitySummary(),
			"iru_vulnerability_detail": tableIruSoftwareVulnerability(),
		},
	}
	return p
}
