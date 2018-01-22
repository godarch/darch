package repository

import (
	"context"
	"time"

	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"github.com/pauldotknopf/darch/reference"
)

type Image struct {
	Name      string
	Tag       string
	CreatedAt time.Time
}

// GetImages Get all the built images.
func (session *Session) GetImages(ctx context.Context) ([]Image, error) {
	ctx = namespaces.WithNamespace(ctx, "darch")

	imgs, err := session.client.ImageService().List(ctx)
	if err != nil {
		return nil, err
	}
	result := []Image{}

	for _, img := range imgs {
		ref, err := reference.ParseImage(img.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, Image{
			Name:      ref.Name,
			Tag:       ref.Tag,
			CreatedAt: img.CreatedAt,
		})
	}

	return result, nil
}

// TagImage Tag an image.
func (session *Session) TagImage(ctx context.Context, source, destination string) error {
	ctx = namespaces.WithNamespace(ctx, "darch")

	sourceRef, err := reference.ParseImage(source)
	if err != nil {
		return err
	}

	destinationRef, err := reference.ParseImage(destination)
	if err != nil {
		return err
	}

	sourceImage, err := session.client.GetImage(ctx, sourceRef.FullName())
	if err != nil {
		return err
	}

	sourceImageTarget := sourceImage.Target()

	_, err = session.client.ImageService().Create(ctx,
		images.Image{
			Name:   destinationRef.FullName(),
			Target: sourceImageTarget,
		})
	if err != nil {
		return err
	}

	return nil
}
