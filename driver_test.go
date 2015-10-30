package brightbox

import (
	"testing"
	"github.com/brightbox/gobrightbox"
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

func getDefaultTestDriverFlags() *DriverOptionsMock {
	return &DriverOptionsMock{
		Data: map[string]interface{}{
			"name":                     "test",
			"url":                      "unix:///var/run/docker.sock",
			"swarm":                    false,
			"swarm-host":               "",
			"swarm-master":             false,
			"swarm-discovery":          "",
			"brightbox-client":         "xyz",
			"brightbox-client-secret":  "abcdefg",
			"brightbox-user-name":      "",
			"brightbox-password":       "",
			"brightbox-account":        "",
			"brightbox-api-url":        brightbox.DefaultRegionApiURL,
			"brightbox-ipv6":           false,
			"brightbox-zone":           "",
			"brightbox-image":          "",
			"brightbox-type":           "",
			"brightbox-security-group": []string(nil),
		},
	}
}

func TestPasswordCredentialsValidation(t *testing.T) {
	drive := new(Driver)
	flags := getDefaultTestDriverFlags()
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
	flags := getDefaultTestDriverFlags()
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
