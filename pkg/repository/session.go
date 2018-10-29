package repository

import (
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/snapshots"
)

var (
	// DefaultContainerdSocketLocation The location to the containerd socket.
	DefaultContainerdSocketLocation = "/var/run/containerd/containerd.sock"
)

// Session An object that represent a session to a containerd runtime.
type Session struct {
	client      *containerd.Client
	snapshotter snapshots.Snapshotter
	imagesStore images.Store
	differ      containerd.DiffService
	content     content.Store
}

// NewSession creates a new session
func NewSession(containerdSocket string) (*Session, error) {
	client, err := containerd.New(containerdSocket)
	if err != nil {
		return nil, err
	}

	client.ImageService()

	return &Session{
		client:      client,
		snapshotter: client.SnapshotService(containerd.DefaultSnapshotter),
		imagesStore: client.ImageService(),
		differ:      client.DiffService(),
		content:     client.ContentStore(),
	}, nil
}

// Close Closes the session.
func (session *Session) Close() error {
	return session.client.Close()
}
