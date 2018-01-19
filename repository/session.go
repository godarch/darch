package repository

import (
	"github.com/containerd/containerd"
)

var (
	// DefaultContainerdSocketLocation The location to the containerd socket.
	DefaultContainerdSocketLocation = "/var/run/containerd/containerd.sock"
	//DefaultContainerdSocketLocation = "/run/containerd/debug.sock"
)

// Session An object that represent a session to a containerd runtime.
type Session struct {
	client *containerd.Client
}

// NewSession creates a new session
func NewSession(containerdSocket string) (*Session, error) {
	client, err := containerd.New(containerdSocket)
	if err != nil {
		return nil, err
	}

	return &Session{
		client: client,
	}, nil
}

// Close Closes the session.
func (session *Session) Close() error {
	return session.client.Close()
}
