package repository

import (
	"context"
	"github.com/containerd/containerd"

	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/remotes"
	"github.com/godarch/darch/pkg/reference"
)

// Push Push an image remotely.
func (session *Session) Push(ctx context.Context, imageRef reference.ImageRef, resolver remotes.Resolver) error {
	ctx = namespaces.WithNamespace(ctx, "darch")
	image, err := session.client.GetImage(ctx, imageRef.FullName())
	if err != nil {
		return err
	}

	pushRef := imageRef
	if len(pushRef.Domain()) == 0 {
		parsedRef, err := pushRef.WithDomain(reference.DefaultDomain)
		if err != nil {
			return err
		}
		pushRef = parsedRef
	}

	err = session.client.Push(ctx,
		pushRef.FullName(),
		image.Target(),
		containerd.WithResolver(&overrideNameResolve{
			RealResolver: resolver,
			Name:         imageRef.FullName(),
			FullRef:      pushRef.FullName(),
		}))
	return err
}
