package staging

import "github.com/pauldotknopf/darch/reference"

// RunAllHooks Run the hooks on every image.
func (session *Session) RunAllHooks() error {
	allAssoications, err := session.imageStore.AllImages()
	if err != nil {
		return err
	}
	for _, association := range allAssoications {
		err = session.runHookForAssociation(association)
		if err != nil {
			return err
		}
	}
	return nil
}

// RunHookForImage Run hooks for a single image.
func (session *Session) RunHookForImage(imageRef reference.ImageRef) error {
	association, err := session.imageStore.Get(imageRef)
	if err != nil {
		return err
	}
	return session.runHookForAssociation(association)
}

func (session *Session) runHookForAssociation(association reference.Association) error {
	return nil
}
