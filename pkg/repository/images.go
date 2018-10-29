package repository

import (
	"context"
	"time"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"github.com/godarch/darch/pkg/reference"
	"github.com/pkg/errors"
)

// Image An image fetched from the repository.
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
			Name:      ref.Name(),
			Tag:       ref.Tag(),
			CreatedAt: img.CreatedAt,
		})
	}

	return result, nil
}

// TagImage Tag an image.
func (session *Session) TagImage(ctx context.Context, source, destination reference.ImageRef) error {
	ctx = namespaces.WithNamespace(ctx, "darch")

	sourceImage, err := session.client.GetImage(ctx, source.FullName())
	if err != nil {
		return err
	}

	// Prevent garbage collection while we work.
	ctx, done, err := session.client.WithLease(ctx)
	if err != nil {
		return err
	}
	defer done(ctx)

	// Make sure the destination image:tag doesn't already exist.
	destinationImage, err := session.client.GetImage(ctx, destination.FullName())
	if err != nil && errors.Cause(err) != errdefs.ErrNotFound {
		return err
	}
	if err == nil {
		err = session.imagesStore.Delete(ctx, destinationImage.Name(), images.SynchronousDelete())
		if err != nil {
			return err
		}
	}

	_, err = session.client.ImageService().Create(ctx,
		images.Image{
			Name:   destination.FullName(),
			Target: sourceImage.Target(),
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveImage Removes an image locally.
func (session *Session) RemoveImage(ctx context.Context, image string) error {
	ctx = namespaces.WithNamespace(ctx, "darch")
	ref, err := reference.ParseImage(image)
	if err != nil {
		return err
	}
	return session.client.ImageService().Delete(ctx, ref.FullName(), images.SynchronousDelete())
}
