package reference

import (
	"os"
	"path"
	"testing"

	"github.com/godarch/darch/pkg/utils"
)

func TestCreatesOnNew(t *testing.T) {
	t.Parallel()

	jsonFile := path.Join(os.TempDir(), utils.NewID())
	defer os.RemoveAll(jsonFile)

	_, err := NewReferenceStore(jsonFile)
	if err != nil {
		t.Fatalf("error creating store %v", err)
	}

	if !utils.FileExists(jsonFile) {
		t.Fatal("file not created on load")
	}
}

func TestTag(t *testing.T) {
	jsonFile := path.Join(os.TempDir(), utils.NewID())
	defer os.RemoveAll(jsonFile)

	store, err := NewReferenceStore(jsonFile)
	if err != nil {
		t.Fatalf("error creating store %v", err)
	}

	newImage, _ := ParseImage("base:latest")
	newID := utils.NewID()

	err = store.AddTag(newImage, newID, false)
	if err != nil {
		t.Fatalf("error adding image %v", err)
	}

	retrievedID, err := store.Get(newImage)
	if err != nil {
		t.Fatalf("error getting id %v", err)
	}

	if retrievedID.ID != newID {
		t.Fatalf("invalid id retrieved")
	}
}

func TestTagTwiceSameId(t *testing.T) {
	jsonFile := path.Join(os.TempDir(), utils.NewID())
	defer os.RemoveAll(jsonFile)

	store, err := NewReferenceStore(jsonFile)
	if err != nil {
		t.Fatalf("error creating store %v", err)
	}

	newImage, _ := ParseImage("base:latest")
	newID := utils.NewID()

	err = store.AddTag(newImage, newID, false)
	if err != nil {
		t.Fatalf("error adding image %v", err)
	}

	err = store.AddTag(newImage, newID, false)
	if err != nil {
		t.Fatalf("error adding image %v", err)
	}
}

func TestTagTwiceDifferentIdError(t *testing.T) {
	jsonFile := path.Join(os.TempDir(), utils.NewID())
	defer os.RemoveAll(jsonFile)

	store, err := NewReferenceStore(jsonFile)
	if err != nil {
		t.Fatalf("error creating store %v", err)
	}

	newImage, _ := ParseImage("base:latest")
	newID1 := utils.NewID()
	newID2 := utils.NewID()

	err = store.AddTag(newImage, newID1, false)
	if err != nil {
		t.Fatalf("error adding image %v", err)
	}

	err = store.AddTag(newImage, newID2, false)
	if err == nil {
		t.Fatalf("should have thrown error")
	}
}

func TestTagTwiceDifferentIdForce(t *testing.T) {
	jsonFile := path.Join(os.TempDir(), utils.NewID())
	defer os.RemoveAll(jsonFile)

	store, err := NewReferenceStore(jsonFile)
	if err != nil {
		t.Fatalf("error creating store %v", err)
	}

	newImage, _ := ParseImage("base:latest")
	newID1 := utils.NewID()
	newID2 := utils.NewID()

	err = store.AddTag(newImage, newID1, false)
	if err != nil {
		t.Fatalf("error adding image %v", err)
	}

	err = store.AddTag(newImage, newID2, true)
	if err != nil {
		t.Fatalf("should have forced the update to new id")
	}
}

func TestErrorReturnDeleteNonExistingImage(t *testing.T) {
	jsonFile := path.Join(os.TempDir(), utils.NewID())
	defer os.RemoveAll(jsonFile)

	store, err := NewReferenceStore(jsonFile)
	if err != nil {
		t.Fatalf("error creating store %v", err)
	}

	nonExistantImage, _ := ParseImage("base:latest")

	result, err := store.Delete(nonExistantImage)
	if err != ErrDoesNotExist {
		t.Fatalf("wrong error returned")
	}
	if result {
		t.Fatalf("says we updated, when we shouldn't have")
	}
}

func TestGetImagesForId(t *testing.T) {
	jsonFile := path.Join(os.TempDir(), utils.NewID())
	defer os.RemoveAll(jsonFile)

	store, err := NewReferenceStore(jsonFile)
	if err != nil {
		t.Fatalf("error creating store %v", err)
	}

	id := utils.NewID()
	newImage1, _ := ParseImage("base:latest")
	newImage2, _ := ParseImage("base2:latest")

	err = store.AddTag(newImage1, id, false)
	if err != nil {
		t.Fatalf("couldn't add image")
	}
	err = store.AddTag(newImage2, id, false)
	if err != nil {
		t.Fatalf("couln't add image")
	}

	images, err := store.References(id)
	if len(images) != 2 {
		t.Fatal("invalid images count")
	}
	if images[0].FullName() != "base:latest" {
		t.Fatalf("invalid image")
	}
	if images[1].FullName() != "base2:latest" {
		t.Fatalf("invalid image")
	}
}
