package brightbox

import (
	"github.com/brightbox/gobrightbox"
	"testing"
)

type DriverOptionsMock struct {
	Data map[string]interface{}
}

func (d DriverOptionsMock) String(key string) string {
	return d.Data[key].(string)
}

func (d DriverOptionsMock) StringSlice(key string) []string {
	return d.Data[key].([]string)
}

func (d DriverOptionsMock) Int(key string) int {
	return d.Data[key].(int)
}

func (d DriverOptionsMock) Bool(key string) bool {
	return d.Data[key].(bool)
}

func getDefaultTestDriverFlags(d *Driver) *DriverOptionsMock {
	data := make(map[string]interface{})
	mcnflags := d.GetCreateFlags()
	for _, f := range mcnflags {
		data[f.String()] = f.Default()
		if f.Default() == nil {
			data[f.String()] = false
		}
	}
	data["brightbox-client"] = "xyz"
	data["brightbox-client-secret"] = "abcdefg"
	data["brightbox-image"] = "img-freda"
	return &DriverOptionsMock{
		Data: data,
	}
}

func TestPasswordCredentialsValidation(t *testing.T) {
	drive := new(Driver)
	flags := getDefaultTestDriverFlags(drive)
	flags.Data["brightbox-user-name"] = "testuser"
	if err := drive.SetConfigFromFlags(flags); err == nil {
		t.Error("Missing password not picked up when Username present")
	}
	flags.Data["brightbox-password"] = "password"
	if err := drive.SetConfigFromFlags(flags); err != nil {
		t.Error("Username and password rejected")
	}
	flags.Data["brightbox-user-name"] = ""
	if err := drive.SetConfigFromFlags(flags); err == nil {
		t.Error("Missing Username not picked up when password present")
	}
}

func TestClientValidation(t *testing.T) {
	drive := new(Driver)
	flags := getDefaultTestDriverFlags(drive)
	flags.Data["brightbox-client"] = defaultClientID
	flags.Data["brightbox-client-secret"] = defaultClientSecret
	if err := drive.SetConfigFromFlags(flags); err == nil {
		t.Error("Missing Client ID not picked up")
	}
}

func TestDriverName(t *testing.T) {
	drive := new(Driver)
	if drive.DriverName() != "brightbox" {
		t.Error("Driver Name should be brightbox")
	}
}

func TestDefaultValues(t *testing.T) {
	driver := new(Driver)
	flags := getDefaultTestDriverFlags(driver)
	if err := driver.SetConfigFromFlags(flags); err != nil {
		t.Fatal("Unexpected set config failure testing defaults")
	}
	if driver.Account != "" {
		t.Errorf("Incorrect default Account: %s", driver.Account)
	}
	if driver.ApiURL != brightbox.DefaultRegionApiURL {
		t.Errorf("Incorrect default API URL: %s", driver.ApiURL)
	}
	if driver.ServerType != defaultServerType {
		t.Errorf("Incorrect default ServerType: %s", driver.ServerType)
	}
	if driver.IPv6 != defaultIPV6 {
		t.Errorf("Incorrect default IPV6: %s", driver.IPv6)
	}
	if driver.ServerGroups != nil {
		t.Errorf("Incorrect default ServerGroups, %v", driver.ServerGroups)
	}
	if driver.Zone != "" {
		t.Errorf("Incorrect default Zone: %s", driver.Zone)
	}
	if driver.SSHPort != defaultSSHPort {
		t.Errorf("Incorrect default SSHPort: %s", driver.SSHPort)
	}
	if *driver.Name != " (docker-machine)" {
		t.Errorf("Incorrect default Name: %s", *driver.Name)
	}
}
