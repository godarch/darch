package staging

import "github.com/godarch/darch/pkg/reference"

// Remove Removes an image from the stage.
func (session *Session) Remove(imageRef reference.ImageRef) error {
	result, err := session.imageStore.Delete(imageRef)

	if result {
		// We deleted the image.
		// Let's do a clean up, which will delete the local data,
		// if it isn't referenced anymore.
		return session.Clean()
	}

	return err
}
