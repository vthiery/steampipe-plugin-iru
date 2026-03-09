package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/vthiery/steampipe-plugin-iru/iru"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: iru.Plugin})
}
