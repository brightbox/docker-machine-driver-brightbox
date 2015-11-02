package brightbox

import (
	"fmt"
	"sort"
	"strings"

	"github.com/brightbox/gobrightbox"
)

const (
	DefaultArch     = "x86_64"
	//The tag is looked for in the name of the image
	DefaultImageTag = "CoreOS"
)

/*
Searches the supplied Image List for official Images and selects one
according to the constant definitions in this file. Returns a reference
to that Image.  If no Image matches the default settings you will get
an error.
*/
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
	return image.Official && image.Arch == DefaultArch && strings.Contains(image.Name, DefaultImageTag)
}

type imagesByAgeDescending []*brightbox.Image

func (a imagesByAgeDescending) Len() int      { return len(a) }
func (a imagesByAgeDescending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imagesByAgeDescending) Less(i, j int) bool {
	return strings.ToLower(a[i].Name) > strings.ToLower(a[j].Name)
}
