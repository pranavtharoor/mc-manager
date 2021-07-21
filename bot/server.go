package bot

import (
	"github.com/pranavtharoor/mc-manager/azure"
	"github.com/pranavtharoor/mc-manager/config"
)

func serverStart(c config.ServerConfiguration) string {
	err := azure.VMStart(c.ResourceGroup, c.Name)
	if err != nil {
		return err.Error()
	}
	return "Starting up server..."
}

func serverStop(c config.ServerConfiguration) string {
	err := azure.VMDeallocate(c.ResourceGroup, c.Name)
	if err != nil {
		return err.Error()
	}
	return "Shutting down server..."
}
