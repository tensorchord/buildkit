package main

import (
	"context"
	"path/filepath"

	"github.com/docker/docker/api/types"
	_ "github.com/docker/docker/daemon/graphdriver/overlay2"
	"github.com/docker/docker/daemon/images"
	dmetadata "github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/idtools"
	refstore "github.com/docker/docker/reference"
)

func main() {
	root := "/var/lib/docker"
	graphDriver := "overlay2"
	layerStore, err := layer.NewStoreFromOptions(layer.StoreOptions{
		Root:                      root,
		MetadataStorePathTemplate: filepath.Join(root, "image", "%s", "layerdb"),
		GraphDriver:               graphDriver,
		GraphDriverOptions:        []string{},
		IDMapping:                 idtools.IdentityMapping{},
		ExperimentalEnabled:       false,
	})
	if err != nil {
		panic(err)
	}
	m := layerStore.Map()
	for k, v := range m {
		println(k, v)
	}
	imageRoot := filepath.Join(root, "image", graphDriver)
	ifs, err := image.NewFSStoreBackend(filepath.Join(imageRoot, "imagedb"))
	if err != nil {
		panic(err)
	}

	imageStore, err := image.NewImageStore(ifs, layerStore)
	if err != nil {
		panic(err)
	}
	im := imageStore.Map()
	for k, v := range im {
		println(k, v.Size)
	}

	refStoreLocation := filepath.Join(imageRoot, `repositories.json`)
	rs, err := refstore.NewReferenceStore(refStoreLocation)
	if err != nil {
		panic(err)
	}
	_ = rs
	distributionMetadataStore, err := dmetadata.NewFSMetadataStore(filepath.Join(imageRoot, "distribution"))
	if err != nil {
		panic(err)
	}
	_ = distributionMetadataStore
	imgSvcConfig := images.ImageServiceConfig{
		DistributionMetadataStore: distributionMetadataStore,
		ImageStore:                imageStore,
		LayerStore:                layerStore,
		ReferenceStore:            rs,
	}
	imageService := images.NewImageService(imgSvcConfig)
	is, err := imageService.Images(context.TODO(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _, i := range is {
		println(i.ID)
	}
}
