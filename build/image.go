package build

import (
	"regexp"
)

// ParseImageTag parse the given string to image name and tag name.
func ParseImageTag(s string) (image, tag string) {
	exp := regexp.MustCompile(`^([A-z0-9/\-_.]+)(:([A-z0-9\-_.]+))*$`)
	matched := exp.FindStringSubmatch(s)

	switch len(matched) {
	case 2:
		image = matched[1]
	case 4:
		image = matched[1]
		tag = matched[3]
	}
	return
}
