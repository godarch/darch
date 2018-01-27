package repository

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/remotes"

	"github.com/containerd/containerd/namespaces"
	"github.com/pauldotknopf/darch/pkg/reference"
)

// Pull Pulls an image locally.
func (session *Session) Pull(ctx context.Context, imageRef reference.ImageRef, resolver remotes.Resolver) error {
	_, err := session.client.Pull(namespaces.WithNamespace(ctx, "darch"),
		imageRef.FullName(),
		containerd.WithResolver(resolver),
		containerd.WithPullUnpack)
	return err
}
