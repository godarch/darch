package reference

import (
	dockerref "github.com/docker/distribution/reference"
	"strings"
)

const (
	// DefaultDomain The default domain, if none was specified.
	DefaultDomain = "docker.io"
)

// ImageRef An image reference.
type ImageRef interface {
	Tag() string
	Name() string
	Domain() string
	FullName() string
	WithTag(tag string) (ImageRef, error)
	WithName(name string) (ImageRef, error)
	WithDomain(domain string) (ImageRef, error)
}

type imageRef struct {
	fullname string
}

// ParseImage Parses a string for image:tag.
func ParseImage(val string) (ImageRef, error) {
	return ParseImageWithDefaultTag(val, "latest")
}

// ParseImageWithDefaultTag Parse an image name. If not tag is given in the image name, use the optionalTag as the tag.
func ParseImageWithDefaultTag(val, optionalTag string) (ImageRef, error) {
	ref, err := dockerref.Parse(val)
	if err != nil {
		return nil, err
	}

	if _, ok := ref.(dockerref.Tagged); !ok {
		// There was no tag
		val = val + ":" + optionalTag
		ref, err = dockerref.Parse(val)
	}

	return &imageRef{
		fullname: val,
	}, nil
}

func (image *imageRef) Tag() string {
	ref, _ := dockerref.Parse(image.fullname)
	tagged, _ := ref.(dockerref.Tagged)
	return tagged.Tag()
}

func (image *imageRef) Name() string {
	ref, _ := dockerref.Parse(image.fullname)
	named, _ := ref.(dockerref.Named)
	return named.Name()
}

func (image *imageRef) Domain() string {
	ref, _ := dockerref.Parse(image.fullname)
	named, _ := ref.(dockerref.Named)
	domain, _ := splitDomain(named.Name())
	return domain
}

func (image *imageRef) FullName() string {
	return image.fullname
}

func (image *imageRef) WithTag(tag string) (ImageRef, error) {
	ref, _ := dockerref.Parse(image.fullname)
	named, _ := ref.(dockerref.Named)
	namedTagged, err := dockerref.WithTag(named, tag)
	if err != nil {
		return nil, err
	}
	return &imageRef{
		fullname: namedTagged.String(),
	}, nil
}

func (image *imageRef) WithName(name string) (ImageRef, error) {
	ref, _ := dockerref.Parse(image.fullname)
	tagged, _ := ref.(dockerref.Tagged)

	ref, err := dockerref.Parse(name + ":" + tagged.Tag())
	if err != nil {
		return nil, err
	}

	return &imageRef{
		fullname: ref.String(),
	}, nil
}

func (image *imageRef) WithDomain(domain string) (ImageRef, error) {
	ref, _ := dockerref.Parse(image.fullname)
	named, _ := ref.(dockerref.Named)
	_, name := splitDomain(named.Name())
	if len(domain) == 0 {
		return image.WithName(name)
	}
	return image.WithName(domain + "/" + name)
}

func splitDomain(name string) (domain, remainder string) {
	i := strings.IndexRune(name, '/')
	if i == -1 || (!strings.ContainsAny(name[:i], ".:") && name[:i] != "localhost") {
		domain, remainder = "", name
	} else {
		domain, remainder = name[:i], name[i+1:]
	}
	return domain, remainder
}
