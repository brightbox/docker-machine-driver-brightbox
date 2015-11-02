package brightbox

import (
	"github.com/brightbox/gobrightbox"
	"testing"
)

func TestEmptyImages(t *testing.T) {
	emptyImages := []brightbox.Image{}
	if _, err := GetDefaultImage(emptyImages); err == nil {
		t.Error("Missing default image not detected in empty list")
	}
}

func TestSingleImageFound(t *testing.T) {
	singleImage := []brightbox.Image{
		{
			Resource: brightbox.Resource{
				Id: "img-upwxc",
			},
			Name:              "CoreOS 766.4.0",
			Owner:             "brightbox",
			Arch:              "x86_64",
			Description:       "ID: com.brightbox:test/net.core-os.release:amd64-usr/766.4.0/disk1.img, Release: stable",
			Username:          "core",
			Official:          true,
			Public:            true,
			CompatibilityMode: false,
		},
	}
	image, err := GetDefaultImage(singleImage)
	if err != nil {
		t.Fatal(err)
	}
	if image.Id != "img-upwxc" {
		t.Error("Failed to select correct image")
	}
}

func TestSingleImageNotFound(t *testing.T) {
	singleImage := []brightbox.Image{
		{
			Resource: brightbox.Resource{
				Id: "img-abcde",
			},
			Name:              "ubuntu-wily-daily-amd64-server",
			Owner:             "brightbox",
			Arch:              "x86_64",
			Description:       "ID: com.ubuntu.cloud:daily:download/com.ubuntu.cloud.daily:server:15.10:amd64/20151026/disk1.img, Release: daily",
			Username:          "ubuntu",
			Official:          true,
			Public:            true,
			CompatibilityMode: false,
		},
	}
	image, err := GetDefaultImage(singleImage)
	if err == nil {
		t.Error("Expected no image")
	}
	if image != nil {
		t.Errorf("Received image reference %s when not expected", image.Id)
	}
}

func TestFilterAndSort(t *testing.T) {
	multipleImages := []brightbox.Image{
		{
			Resource: brightbox.Resource{
				Id: "img-upwxc",
			},
			Name:              "CoreOS 766.4.0",
			Owner:             "brightbox",
			Arch:              "x86_64",
			Description:       "ID: com.brightbox:test/net.core-os.release:amd64-usr/766.4.0/disk1.img, Release: stable",
			Username:          "core",
			Official:          true,
			Public:            true,
			CompatibilityMode: false,
		},
		{
			Resource: brightbox.Resource{
				Id: "img-gnhsz",
			},
			Name:              "CoreOS 845.0.0",
			Owner:             "brightbox",
			Arch:              "x86_64",
			Description:       "ID: com.brightbox:test/net.core-os.release:amd64-usr/845.0.0/disk1.img, Release: alpha",
			Username:          "core",
			Official:          true,
			Public:            true,
			CompatibilityMode: false,
		},
		{
			Resource: brightbox.Resource{
				Id: "img-77dmp",
			},
			Name:              "ubuntu-wily-15.10-amd64-server-uefi1",
			Owner:             "brightbox",
			Arch:              "x86_64",
			Description:       "ID: com.ubuntu.cloud:released:download/com.ubuntu.cloud:server:15.10:amd64/20151021/uefi1.img, Release: release",
			Username:          "ubuntu",
			Official:          true,
			Public:            true,
			CompatibilityMode: false,
		},
		{
			Resource: brightbox.Resource{
				Id: "img-b0ieg",
			},
			Name:              "ubuntu-wily-15.10-i386-server",
			Owner:             "brightbox",
			Arch:              "i686",
			Description:       "ID: com.ubuntu.cloud:released:download/com.ubuntu.cloud:server:15.10:i386/20151026/disk1.img, Release: release",
			Username:          "ubuntu",
			Official:          true,
			Public:            true,
			CompatibilityMode: false,
		},
		{
			Resource: brightbox.Resource{
				Id: "img-5atge",
			},
			Name:              "ubuntu-wily-15.10-amd64-server",
			Owner:             "brightbox",
			Arch:              "x86_64",
			Description:       "ID: com.ubuntu.cloud:released:download/com.ubuntu.cloud:server:15.10:amd64/20151021/disk1.img, Release: release",
			Username:          "ubuntu",
			Official:          true,
			Public:            true,
			CompatibilityMode: false,
		},
		{
			Resource: brightbox.Resource{
				Id: "img-abcde",
			},
			Name:              "ubuntu-wily-daily-amd64-server",
			Owner:             "acc-7wy80",
			Arch:              "x86_64",
			Description:       "ID: com.ubuntu.cloud:daily:download/com.ubuntu.cloud.daily:server:15.10:amd64/20151026/disk1.img, Release: daily",
			Username:          "ubuntu",
			Official:          false,
			Public:            true,
			CompatibilityMode: false,
		},
	}
	image, err := GetDefaultImage(multipleImages)
	if err != nil {
		t.Fatal(err)
	}
	if image.Id != "img-gnhsz" {
		t.Errorf("Received image reference %s - expecting img-77dmp", image.Id)
	}
}
