package brightbox

import (
	"fmt"
	"sort"
	"strings"

	"github.com/brightbox/gobrightbox"
)

const (
	defaultArch     = "x86_64"
	defaultImageTag = "ubuntu"
)

//Pass in a list of Images obtained from the API and a reference to the default
//image within that list selected according to the constants in this file
func GetDefaultImage(images []brightbox.Image) (*brightbox.Image, error) {
	filteredImages := filterImages(images, defaultImage)
	switch len(filteredImages) {
	case 0:
		return nil, fmt.Errorf("Unable to find a default Image")
	}
	sort.Sort(imagesByAgeDescending(filteredImages))
	return filteredImages[0], nil
}

func filterImages(images []brightbox.Image, selector func(*brightbox.Image) bool) []*brightbox.Image {
	result := make([]*brightbox.Image, 0)
	for index := range images {
		if selector(&images[index]) {
			result = append(result, &images[index])
		}
	}
	return result
}

func defaultImage(image *brightbox.Image) bool {
	return image.Official && image.Arch == defaultArch && strings.Contains(image.Name, defaultImageTag)
}

type imagesByAgeDescending []*brightbox.Image

func (a imagesByAgeDescending) Len() int      { return len(a) }
func (a imagesByAgeDescending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imagesByAgeDescending) Less(i, j int) bool {
	return strings.ToLower(a[i].Name) > strings.ToLower(a[j].Name)
}
