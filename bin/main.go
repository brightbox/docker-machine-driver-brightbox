package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/brightbox/docker-machine-driver-brightbox"
)

func main() {
	plugin.RegisterDriver(new(brightbox.Driver))
}
