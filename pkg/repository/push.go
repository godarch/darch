package repository

import (
	"context"
	"github.com/containerd/containerd"

	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/remotes"
	"github.com/pauldotknopf/darch/pkg/reference"
)

// Push Push an image remotely.
func (session *Session) Push(ctx context.Context, imageRef reference.ImageRef, resolver remotes.Resolver) error {
	ctx = namespaces.WithNamespace(ctx, "darch")
	image, err := session.client.GetImage(ctx, imageRef.FullName())
	if err != nil {
		return err
	}
	err = session.client.Push(ctx,
		image.Name(),
		image.Target(),
		containerd.WithResolver(resolver))
	return err
}
