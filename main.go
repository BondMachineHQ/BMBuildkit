package main

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func main() {
	// Image references
	base := empty.Index
	image1Ref := "docker.io/dciangot/test:v102"
	image2Ref := "docker.io/dciangot/test:v101"

	// // Annotation for image 1
	// image1Annotations := map[string]string{
	// 	"author":  "John Doe",
	// 	"purpose": "Backend Service",
	// }

	// // Annotation for image 2
	// image2Annotations := map[string]string{
	// 	"author":  "Jane Smith",
	// 	"purpose": "Frontend Service",
	// }

	img1ref, err := name.ParseReference(image1Ref)
	if err != nil {
		panic(err)
	}

	img1, err := remote.Get(img1ref)
	if err != nil {
		panic(err)
	}

	img1GOOD, err := img1.Image()
	if err != nil {
		panic(err)
	}

	img2ref, err := name.ParseReference(image2Ref)
	if err != nil {
		panic(err)
	}

	img2, err := remote.Get(img2ref)
	if err != nil {
		panic(err)
	}

	img2GOOD, err := img2.Image()
	if err != nil {
		panic(err)
	}

	// Create a new index image
	idx := mutate.AppendManifests(base,
		mutate.IndexAddendum{Add: img2GOOD, Descriptor: v1.Descriptor{
			Annotations: map[string]string{
				"foo": "bar",
			},
		}},
		mutate.IndexAddendum{Add: img1GOOD},
	)

	// Push the index image to a registry
	dest := "dciangot/index_image:latest"
	tag, err := name.NewTag(dest, name.WeakValidation)
	if err != nil {
		panic(err)
	}

	err = remote.WriteIndex(tag, idx, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Index image %s created and pushed successfully.\n", dest)
}
