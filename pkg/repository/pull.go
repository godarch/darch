package repository

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/remotes"

	"github.com/containerd/containerd/namespaces"
	"github.com/godarch/darch/pkg/reference"
)

// Pull Pulls an image locally.
func (session *Session) Pull(ctx context.Context, imageRef reference.ImageRef, resolver remotes.Resolver) error {
	originalRef := imageRef
	newRef := imageRef
	if len(newRef.Domain()) == 0 {
		parsedRef, err := newRef.WithDomain(reference.DefaultDomain)
		if err != nil {
			return err
		}
		newRef = parsedRef
	}
	_, err := session.client.Pull(namespaces.WithNamespace(ctx, "darch"),
		newRef.FullName(),
		containerd.WithResolver(resolver),
		containerd.WithPullUnpack)

	_ = originalRef

	return err
}
