// Brightbox Cloud Driver for Docker Machine
package brightbox

import (
	"fmt"

	"github.com/brightbox/gobrightbox"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	//	"github.com/docker/machine/libmachine/ssh"
	//"github.com/docker/machine/libmachine/state"
)

const (
	// Docker Machine application client credentials
	defaultClientID     = "app-dkmch"
	defaultClientSecret = "uogoelzgt0nwawb"

	defaultSSHPort = 22
	driverName     = "brightbox"
)

type Driver struct {
	drivers.BaseDriver
	authdetails
	brightbox.ServerOptions
	IPv6       bool
	liveClient *brightbox.Client
}

//Backward compatible Driver factory method
//Using new(brightbox.Driver) is preferred
func NewDriver(hostName, storePath string) Driver {
	return Driver{
		BaseDriver: drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "BRIGHTBOX_CLIENT",
			Name:   "brightbox-client",
			Usage:  "Brightbox Cloud API Client",
			Value:  defaultClientID,
		},
		mcnflag.StringFlag{
			EnvVar: "BRIGHTBOX_CLIENT_SECRET",
			Name:   "brightbox-client-secret",
			Usage:  "Brightbox Cloud API Client Secret",
			Value:  defaultClientSecret,
		},
		mcnflag.StringFlag{
			EnvVar: "BRIGHTBOX_USER_NAME",
			Name:   "brightbox-user-name",
			Usage:  "Brightbox Cloud User Name",
		},
		mcnflag.StringFlag{
			EnvVar: "BRIGHTBOX_PASSWORD",
			Name:   "brightbox-password",
			Usage:  "Brightbox Cloud Password for User Name",
		},
		mcnflag.StringFlag{
			EnvVar: "BRIGHTBOX_ACCOUNT",
			Name:   "brightbox-account",
			Usage:  "Brightbox Cloud Account to operate on",
		},
		mcnflag.StringFlag{
			EnvVar: "BRIGHTBOX_API_URL",
			Name:   "brightbox-api-url",
			Usage:  "Brightbox Cloud Api URL for selected Region",
			Value:  brightbox.DefaultRegionApiURL,
		},
		mcnflag.BoolFlag{
			EnvVar: "BRIGHTBOX_IPV6",
			Name:   "brightbox-ipv6",
			Usage:  "Access server directly over IPv6",
		},
		mcnflag.StringFlag{
			EnvVar: "BRIGHTBOX_ZONE",
			Name:   "brightbox-zone",
			Usage:  "Brightbox Cloud Availability Zone ID",
		},
		mcnflag.StringFlag{
			EnvVar: "BRIGHTBOX_IMAGE",
			Name:   "brightbox-image",
			Usage:  "Brightbox Cloud Image ID",
		},
		mcnflag.StringSliceFlag{
			EnvVar: "BRIGHTBOX_GROUP",
			Name:   "brightbox-group",
			Usage:  "Brightbox Cloud Security Group",
		},
		mcnflag.StringFlag{
			EnvVar: "BRIGHTBOX_TYPE",
			Name:   "brightbox-type",
			Usage:  "Brightbox Cloud Server Type",
		},
	}
}

func (d *Driver) DriverName() string {
	return driverName
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.APIClient = flags.String("brightbox-client")
	d.apiSecret = flags.String("brightbox-client-secret")
	d.UserName = flags.String("brightbox-user-name")
	d.password = flags.String("brightbox-password")
	d.Account = flags.String("brightbox-account")
	d.Image = flags.String("brightbox-image")
	d.ApiURL = flags.String("brightbox-api-url")
	d.ServerType = flags.String("brightbox-type")
	d.IPv6 = flags.Bool("brightbox-ipv6")
	group_list := flags.StringSlice("brightbox-security-group")
	if group_list != nil {
		d.ServerGroups = &group_list
	}
	d.Zone = flags.String("brightbox-zone")
	d.SSHPort = defaultSSHPort
	return d.checkConfig()
}

// Try and avoid authenticating more than once
// Store the authenticated api client in the driver for future use
func (d *Driver) getClient() (*brightbox.Client, error) {
	if d.liveClient != nil {
		log.Debug("Reusing authenticated Brightbox client")
		return d.liveClient, nil
	}
	log.Debug("Authenticating Credentials against Brightbox API")
	client, err := d.authenticatedClient()
	if err == nil {
		d.liveClient = client
		log.Debug("Using authenticated Brightbox client")
	}
	return client, err
}

const (
	errorMandatoryEnvOrOption = "%s must be specified either using the environment variable %s or the CLI option %s"
)

//Statically sanity check flag settings.
func (d *Driver) checkConfig() error {
	switch {
	case d.UserName != "" || d.password != "":
		switch {
		case d.UserName == "":
			return fmt.Errorf(errorMandatoryEnvOrOption, "Username", "BRIGHTBOX_USER_NAME", "--brightbox-user-name")
		case d.password == "":
			return fmt.Errorf(errorMandatoryEnvOrOption, "Password", "BRIGHTBOX_PASSWORD", "--brightbox-password")
		}
	case d.APIClient == defaultClientID:
		return fmt.Errorf(errorMandatoryEnvOrOption, "API Client", "BRIGHTBOX_CLIENT", "--brightbox-client")
	}
	return nil
}

// Make sure that the image details are complete
func (d *Driver) PreCreateCheck() error {
	switch {
	case d.Image == "":
		log.Info("No image specified. Looking for default image")
		client, err := d.getClient()
		if err != nil {
			return err
		}
		images, err := client.Images()
		if err != nil {
			return err
		}
		selectedImage, err := GetDefaultImage(*images)
		if err != nil {
			return err
		}
		d.Image = selectedImage.Id
		d.SSHUser = selectedImage.Username
	case d.SSHUser == "":
		log.Debugf("Looking for Username for Image %s", d.Image)
		client, err := d.getClient()
		if err != nil {
			return err
		}
		image, err := client.Image(d.Image)
		if err != nil {
			return err
		}
		d.SSHUser = image.Username
	}
	log.Debugf("Image %s selected. SSH user is %s", d.Image, d.SSHUser)
	return nil
}
