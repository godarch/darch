package reference

import (
	"os"
	"path"
	"testing"

	"github.com/pauldotknopf/darch/utils"
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

	if retrievedID != newID {
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
