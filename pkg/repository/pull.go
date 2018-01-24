package repository

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
)

// Pull Pulls an image locally.
func (session *Session) Pull(ctx context.Context, image string) (containerd.Image, error) {
	return session.client.Pull(namespaces.WithNamespace(ctx, "darch"),
		image,
		containerd.WithPullUnpack)
}
