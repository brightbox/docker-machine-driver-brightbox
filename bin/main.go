package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/NeilW/docker-machine-driver-brightbox"
)

func main() {
	plugin.RegisterDriver(new(brightbox.Driver))
}
