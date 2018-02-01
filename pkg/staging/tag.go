package staging

import "github.com/godarch/darch/pkg/reference"

// Tag Tag a staged image as something else.
func (session *Session) Tag(sourceImageRef, destinationImageRef reference.ImageRef, force bool) error {
	sourceID, err := session.imageStore.Get(sourceImageRef)
	if err != nil {
		return err
	}

	return session.imageStore.AddTag(destinationImageRef, sourceID.ID, force)
}
