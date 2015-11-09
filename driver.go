//Package brightbox is the Brightbox Cloud Driver for Docker Machine
package brightbox

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/brightbox/gobrightbox"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
)

const (
	// Docker Machine application client credentials
	defaultClientID     = "app-dkmch"
	defaultClientSecret = "uogoelzgt0nwawb"

	// Server creation defaults
	defaultSSHPort    = 22
	defaultIPV6       = true
	defaultServerType = "1gb.ssd"

	driverName     = "brightbox"
	passwordEnvVar = "BRIGHTBOX_PASSWORD"
)

//Driver contains the details necessary for docker-machine to access
//the Brightbox Cloud
type Driver struct {
	drivers.BaseDriver
	authdetails
	brightbox.ServerOptions
	MachineID    string
	IPv6         bool
	mu           sync.Mutex //guards activeClient
	activeClient *brightbox.Client
}

//NewDriver is a backward compatible Driver factory method.  Using
//new(brightbox.Driver) is preferred
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
			EnvVar: passwordEnvVar,
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
			EnvVar: "BRIGHTBOX_IPV4",
			Name:   "brightbox-ipv4",
			Usage:  "Access server over IPv4 rather than IPv6",
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
			Value:  defaultServerType,
		},
	}
}

func (d *Driver) DriverName() string {
	return driverName
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.APIClient = flags.String("brightbox-client")
	d.APISecret = flags.String("brightbox-client-secret")
	d.UserName = flags.String("brightbox-user-name")
	d.password = flags.String("brightbox-password")
	d.Account = flags.String("brightbox-account")
	d.Image = flags.String("brightbox-image")
	d.APIURL = flags.String("brightbox-api-url")
	d.ServerType = flags.String("brightbox-type")
	d.IPv6 = !flags.Bool("brightbox-ipv4")
	groupList := flags.StringSlice("brightbox-group")
	if groupList != nil {
		d.ServerGroups = &groupList
	}
	d.Zone = flags.String("brightbox-zone")
	d.SSHPort = defaultSSHPort
	serverName := d.GetMachineName() + " (docker-machine)"
	d.Name = &serverName
	return d.checkConfig()
}

// Use a mutex to try and avoid authenticating more than once - even
// with parallel calls. Cache the authenticated api client in the driver
// for future use within this run.
func (d *Driver) getClient() (*brightbox.Client, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.activeClient != nil {
		log.Debug("Reusing authenticated Brightbox client")
		return d.activeClient, nil
	}
	log.Debug("Authenticating Credentials against Brightbox API")
	client, err := d.authenticatedClient()
	if err != nil {
		return nil, err
	}
	d.activeClient = client
	if client.AccountId == "" {
		if err := d.setDefaultAccount(); err != nil {
			return nil, err
		}
		log.Debugf("Client Account is %s, Driver Account is %s", client.AccountId, d.Account)
	}
	log.Debug("Using authenticated Brightbox client")
	return client, nil
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

func (d *Driver) checkImage() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}
	var image *brightbox.Image
	if d.Image == "" {
		log.Info("No image specified. Looking for default image")
		log.Debugf("Brightbox API Call: List of Images")
		images, err := client.Images()
		if err != nil {
			return err
		}
		image, err = GetDefaultImage(images)
		if err != nil {
			return err
		}
		d.Image = image.Id
	} else {
		log.Debugf("Brightbox API Call: Image Details for %s", d.Image)
		image, err = client.Image(d.Image)
		if err != nil {
			return err
		}
	}
	if image.Arch != "x86_64" {
		return fmt.Errorf("Docker requires a 64 bit image. Image %s not suitable", d.Image)
	}
	if d.SSHUser == "" {
		log.Debug("Setting SSH Username from image details")
		d.SSHUser = image.Username
	}
	log.Debugf("Image %s selected. SSH user is %s", d.Image, d.SSHUser)
	return nil
}

func (d *Driver) setDefaultAccount() error {
	log.Debug("Looking for default account")
	client := d.activeClient
	log.Debug("Brightbox API Call: List of Accounts")
	accounts, err := client.Accounts()
	switch {
	case err != nil:
		return err
	case len(accounts) == 1:
		log.Debug("Setting default account")
		d.Account = accounts[0].Id
		client.AccountId = d.Account
		return nil
	default:
		return fmt.Errorf(errorMandatoryEnvOrOption, "Account", "BRIGHTBOX_ACCOUNT", "--brightbox-account")
	}
}

// PreCreateCheck makes sure that the image details are complete
func (d *Driver) PreCreateCheck() error {
	return d.checkImage()
}

func (d *Driver) createSSHkey() error {
	return ssh.GenerateSSHKey(d.GetSSHKeyPath())
}

func (d *Driver) getCloudInit() ([]byte, error) {

	data := []byte(`#cloud-config
ssh_authorized_keys:
  - `)
	publickey, err := ioutil.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return nil, err
	}
	return append(data, publickey...), nil
}

func (d *Driver) Create() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}
	log.Infof("Creating SSH key...")
	err = d.createSSHkey()
	if err != nil {
		return err
	}
	userdata, err := d.getCloudInit()
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(userdata)
	d.UserData = &encoded
	log.Infof("Creating Brightbox Server...")
	log.Debugf("with the following Userdata")
	log.Debugf("%s", string(userdata))
	log.Debugf("Brightbox API Call: Create Server using image %s", d.Image)
	server, err := client.CreateServer(&d.ServerOptions)
	if err != nil {
		return err
	}
	d.MachineID = server.Id
	return nil
}

func (d *Driver) getServerDetails() (*brightbox.Server, error) {
	client, err := d.getClient()
	if err != nil {
		return nil, err
	}
	log.Debugf("Brightbox API Call: Server Details for %s", d.MachineID)
	return client.Server(d.MachineID)
}

func (d *Driver) GetIP() (string, error) {
	server, err := d.getServerDetails()
	if err != nil {
		return "", err
	}
	switch {
	case d.IPv6:
		tmp := ipv6Fqdn(server)
		log.Debugf("Returning IPV6 address %s", tmp)
		return tmp, nil
	case len(server.CloudIPs) > 0:
		tmp := publicFqdn(server)
		log.Debugf("Returning public address %s", tmp)
		return tmp, nil
	default:
		log.Debugf("Returning private address %s", server.Fqdn)
		return server.Fqdn, nil
	}
}

func ipv6Fqdn(server *brightbox.Server) string {
	return "ipv6." + server.Fqdn
}

func publicFqdn(server *brightbox.Server) string {
	return "public." + server.Fqdn
}

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetURL() (string, error) {
	fqdn, err := d.GetIP()
	if err != nil {
		return "", err
	}
	return "tcp://" + fqdn + ":2376", nil
}

func (d *Driver) GetState() (state.State, error) {
	server, err := d.getServerDetails()
	if err != nil {
		return state.Error, err
	}
	switch server.Status {
	case "creating":
		log.Debugf("Server %s is Starting", d.MachineID)
		return state.Starting, nil
	case "active":
		log.Debugf("Server %s is Running", d.MachineID)
		return state.Running, nil
	case "deleting":
		log.Debugf("Server %s is Stopping", d.MachineID)
		return state.Stopping, nil
	case "deleted", "inactive":
		log.Debugf("Server %s is Stopped", d.MachineID)
		return state.Stopped, nil
	case "failed", "unavailable":
		log.Debugf("Server %s is Errored", d.MachineID)
		return state.Error, nil
	}
	log.Debugf("Server %s: state not detected", d.MachineID)
	return state.None, nil
}

func (d *Driver) Start() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}
	log.Debugf("Brightbox API Call: Start Server %s", d.MachineID)
	if err := client.StartServer(d.MachineID); err != nil {
		return err
	}
	return nil
}

func (d *Driver) Stop() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}
	log.Debugf("Brightbox API Call: Shutdown Server %s", d.MachineID)
	if err := client.ShutdownServer(d.MachineID); err != nil {
		return err
	}
	return nil
}

func (d *Driver) Restart() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}
	log.Debugf("Brightbox API Call: Reboot Server %s", d.MachineID)
	if err := client.RebootServer(d.MachineID); err != nil {
		return err
	}
	return nil
}

func (d *Driver) Kill() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}
	log.Debugf("Brightbox API Call: Stop Server %s", d.MachineID)
	if err := client.StopServer(d.MachineID); err != nil {
		return err
	}
	return nil
}

func (d *Driver) Remove() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}
	log.Debugf("Brightbox API Call: Destroy Server %s", d.MachineID)
	if err := client.DestroyServer(d.MachineID); err != nil {
		return err
	}
	return nil
}
