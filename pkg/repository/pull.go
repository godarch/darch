package repository

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/remotes"

	"github.com/containerd/containerd/namespaces"
	"github.com/godarch/darch/pkg/reference"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type overrideNameResolve struct {
	RealResolver remotes.Resolver
	Name         string
	FullRef      string
}

func (r *overrideNameResolve) Resolve(ctx context.Context, ref string) (name string, desc ocispec.Descriptor, err error) {
	_, d, e := r.RealResolver.Resolve(ctx, ref)
	return r.Name, d, e
}
func (r *overrideNameResolve) Fetcher(ctx context.Context, ref string) (remotes.Fetcher, error) {
	return r.RealResolver.Fetcher(ctx, r.FullRef)
}
func (r *overrideNameResolve) Pusher(ctx context.Context, ref string) (remotes.Pusher, error) {
	return r.RealResolver.Pusher(ctx, r.FullRef)
}

// Pull Pulls an image locally.
func (session *Session) Pull(ctx context.Context, imageRef reference.ImageRef, resolver remotes.Resolver) (containerd.Image, error) {
	pullRef := imageRef
	if len(pullRef.Domain()) == 0 {
		parsedRef, err := pullRef.WithDomain(reference.DefaultDomain)
		if err != nil {
			return nil, err
		}
		pullRef = parsedRef
	}

	img, err := session.client.Pull(namespaces.WithNamespace(ctx, "darch"),
		pullRef.FullName(),
		containerd.WithResolver(&overrideNameResolve{
			RealResolver: resolver,
			Name:         imageRef.FullName(),
			FullRef:      pullRef.FullName(),
		}),
		containerd.WithPullUnpack)

	return img, err
}
