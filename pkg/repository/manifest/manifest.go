package manifest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images"
	digest "github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const containerdUncompressed = "containerd.io/uncompressed"

// Manifest The manifest that can be mutated.
type Manifest interface {
	AddLayer(ctx context.Context, contentStore content.Store, layer ocispec.Descriptor) error
	Descriptor() ocispec.Descriptor
}

type manifestImpl struct {
	d    map[string]json.RawMessage
	desc ocispec.Descriptor
}

// LoadManifest Load a manifest in-memory for easy interaction.
func LoadManifest(ctx context.Context, contentStore content.Store, desc ocispec.Descriptor) (Manifest, error) {
	p, err := content.ReadBlob(ctx, contentStore, desc.Digest)
	if err != nil {
		return nil, err
	}

	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(p, &m); err != nil {
		return nil, err
	}

	return &manifestImpl{
		d:    m,
		desc: desc,
	}, nil
}

func (m *manifestImpl) AddLayer(ctx context.Context, contentStore content.Store, layer ocispec.Descriptor) error {
	d := m.d

	// These builds can be done on docker images, or OCI image.
	// Let's make sure the new layer uses the same content type as the manifest expects.
	switch m.desc.MediaType {
	case images.MediaTypeDockerSchema2Manifest:
		layer.MediaType = images.MediaTypeDockerSchema2LayerGzip
		break
	case ocispec.MediaTypeImageManifest:
		layer.MediaType = ocispec.MediaTypeImageLayerGzip
		break
	default:
		return fmt.Errorf("unknown parent image manifest type: %s", m.desc.MediaType)
	}

	// Get the diffId for the diff descriptor.
	info, err := contentStore.Info(ctx, layer.Digest)
	if err != nil {
		return err
	}
	diffIDStr, ok := info.Labels[containerdUncompressed]
	if !ok {
		return fmt.Errorf("invalid differ response with no diffID")
	}
	diffIDDigest, err := digest.Parse(diffIDStr)
	if err != nil {
		return err
	}

	// Deserialize the image config
	imageConfigDesc, err := getDescriptor(d["config"])
	if err != nil {
		return err
	}

	// Patch the config and store it in the content store.
	imageConfigDesc, err = patchImageConfig(ctx, contentStore, imageConfigDesc, diffIDDigest)
	if err != nil {
		return err
	}

	// Store the image config back into our json object.
	imageConfigJSON, err := json.Marshal(imageConfigDesc)
	if err != nil {
		return err
	}
	d["config"] = imageConfigJSON

	// Update the layers on the manifest.
	layers := []ocispec.Descriptor{}
	layersJSON, err := d["layers"].MarshalJSON()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(layersJSON, &layers); err != nil {
		return err
	}
	layers = append(layers, layer)
	layersJSON, err = json.Marshal(layers)
	if err != nil {
		return err
	}
	d["layers"] = layersJSON

	// Now that we have all the data ready, let's try to store it.
	// Prepare the labels that will tell the garbage collector
	// to NOT delete the content this manifest references.
	labels := map[string]string{
		"containerd.io/gc.ref.content.0": imageConfigDesc.Digest.String(),
	}
	for i, layer := range layers {
		labels[fmt.Sprintf("containerd.io/gc.ref.content.%d", i+1)] = layer.Digest.String()
	}

	// Save our new image manifest, which now hows our new layer,
	// and a patched image config with a reference to the new layer.
	newDesc := m.desc
	manifestBytes, err := json.Marshal(d)
	if err != nil {
		return err
	}
	newDesc.Digest = digest.FromBytes(manifestBytes)
	newDesc.Size = int64(len(manifestBytes))
	if err := content.WriteBlob(ctx,
		contentStore,
		"custom-ref",
		bytes.NewReader(manifestBytes),
		newDesc.Size,
		newDesc.Digest,
		content.WithLabels(labels)); err != nil {
		return err
	}

	m.desc = newDesc
	m.d = d

	return nil
}

func (m *manifestImpl) Descriptor() ocispec.Descriptor {
	return m.desc
}

func patchImageConfig(ctx context.Context, contentStore content.Store, imageConfig ocispec.Descriptor, newLayer digest.Digest) (ocispec.Descriptor, error) {
	result := imageConfig

	// Get the current image configuration.
	p, err := content.ReadBlob(ctx, contentStore, imageConfig.Digest)
	if err != nil {
		return result, err
	}

	// Deserialize the image configuration to a generic json object.
	// We do this so that we can patch it, without requiring knowledge
	// of the entire schema.
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(p, &m); err != nil {
		return result, err
	}

	// Pull the rootfs section out, so that we can append a layer to the diff_ids array.
	var rootFS ocispec.RootFS
	p, err = m["rootfs"].MarshalJSON()
	if err != nil {
		return result, err
	}
	if err = json.Unmarshal(p, &rootFS); err != nil {
		return result, err
	}
	rootFS.DiffIDs = append(rootFS.DiffIDs, newLayer)
	p, err = json.Marshal(rootFS)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	m["rootfs"] = p

	// Convert our entire image configuration back to bytes, and write it to the content store.
	p, err = json.Marshal(m)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	result.Digest = digest.FromBytes(p)
	result.Size = int64(len(p))
	err = content.WriteBlob(ctx, contentStore,
		"custom-ref",
		bytes.NewReader(p),
		result.Size,
		result.Digest,
	)
	if err != nil {
		return result, err
	}

	return result, nil
}

func getDescriptor(m json.RawMessage) (ocispec.Descriptor, error) {
	var r ocispec.Descriptor
	p, err := m.MarshalJSON()
	if err != nil {
		return r, err
	}
	if err = json.Unmarshal(p, &r); err != nil {
		return r, err
	}
	return r, nil
}
