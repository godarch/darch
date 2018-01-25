package reference

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/pauldotknopf/darch/pkg/utils"
)

var (
	// ErrDoesNotExist is returned if a reference is not found in the
	// store.
	ErrDoesNotExist notFoundError = "reference does not exist"
)

// Store The store holding info about the live images.
type Store interface {
	References(id string) ([]ImageRef, error)
	AddTag(ref ImageRef, id string, force bool) error
	Delete(ref ImageRef) (bool, error)
	Get(ref ImageRef) (Association, error)
	AllImages() ([]Association, error)
}

// Association An association between an id and an image.
type Association struct {
	ID  string
	Ref ImageRef
}

type store struct {
	// TODO: make this object thread safe
	//mu sync.RWMutex
	// jsonPath is the path to the file where the serialized tag data is
	// stored.
	jsonPath string
	// Images is a map of digests, mapped to image names
	Images map[string][]string
}

// NewReferenceStore Create a new store.
func NewReferenceStore(jsonPath string) (Store, error) {
	abspath, err := filepath.Abs(jsonPath)
	if err != nil {
		return nil, err
	}

	store := &store{
		jsonPath: abspath,
		Images:   make(map[string][]string),
	}

	// Load the json file if it exists, otherwise create it.
	if err := store.reload(); os.IsNotExist(err) {
		// TODO: create parent directory, if doesn't exist.
		parentDir := path.Join(path.Dir(jsonPath))
		if !utils.DirectoryExists(parentDir) {
			err = os.MkdirAll(parentDir, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
		if err := store.save(); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return store, nil
}

func (store *store) References(id string) ([]ImageRef, error) {
	result := []ImageRef{}

	images, exists := store.Images[id]
	if !exists || images == nil {
		return result, nil
	}

	for _, image := range images {
		parsed, err := ParseImage(image)
		if err != nil {
			return result, err
		}
		result = append(result, parsed)
	}

	return result, nil
}

func (store *store) AddTag(ref ImageRef, id string, force bool) error {
	// First, make sure it doesn't exist
	existing, err := store.Get(ref)
	if err == ErrDoesNotExist {
		// This is fine
	} else if err != nil {
		return err
	} else if existing.ID == id {
		// Already added
		return nil
	} else if existing.ID != id {
		if !force {
			return fmt.Errorf("tag already added")
		}
		// Delete the current reference, so we can overwrite it.
		result, err := store.Delete(ref)
		if err != nil {
			return err
		}
		if !result {
			return fmt.Errorf("couldn't delete existing reference for image")
		}
	}

	images, exists := store.Images[id]

	if !exists || images == nil {
		images = []string{}
	}

	images = append(images, ref.FullName())
	store.Images[id] = images

	return store.save()
}

func (store *store) Delete(ref ImageRef) (bool, error) {
	outerUpdated := false
	for id, images := range store.Images {
		updated := false
		result := []string{}
		for _, image := range images {
			if image == ref.FullName() {
				updated = true
			} else {
				result = append(result, image)
			}
		}
		if updated {
			outerUpdated = true
			if len(result) == 0 {
				delete(store.Images, id)
			} else {
				store.Images[id] = result
			}
		}
	}

	if !outerUpdated {
		return false, ErrDoesNotExist
	}

	return true, store.save()
}

func (store *store) Get(ref ImageRef) (Association, error) {
	for id, images := range store.Images {
		for _, image := range images {
			if image == ref.FullName() {
				return Association{
					ID:  id,
					Ref: ref,
				}, nil
			}
		}
	}
	return Association{}, ErrDoesNotExist
}

func (store *store) AllImages() ([]Association, error) {
	result := []Association{}
	for id, images := range store.Images {
		for _, image := range images {
			imageRef, err := ParseImage(image)
			if err != nil {
				return result, err
			}
			result = append(result, Association{
				ID:  id,
				Ref: imageRef,
			})
		}
	}
	return result, nil
}

func (store *store) save() error {
	// Store the json
	jsonData, err := json.Marshal(store)
	if err != nil {
		return err
	}
	return ioutils.AtomicWriteFile(store.jsonPath, jsonData, 0600)
}

func (store *store) reload() error {
	f, err := os.Open(store.jsonPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&store)
}
