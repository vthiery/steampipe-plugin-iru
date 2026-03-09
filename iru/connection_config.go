package iru

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// iruConfig stores the connection configuration for the plugin.
type iruConfig struct {
	// URL is the full Iru tenant API URL, e.g. https://acme.api.kandji.io
	URL   *string `hcl:"url"`
	Token *string `hcl:"token"`
}

// ConfigInstance returns a new instance of iruConfig (used by the SDK).
func ConfigInstance() interface{} {
	return &iruConfig{}
}

// GetConfig retrieves and casts the connection config from the plugin query data.
func GetConfig(connection *plugin.Connection) iruConfig {
	if connection == nil || connection.Config == nil {
		return iruConfig{}
	}
	config, _ := connection.Config.(iruConfig)
	return config
}
