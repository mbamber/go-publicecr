package publicecr_test

import (
	"context"
	"log"

	"github.com/mbamber/go-publicecr"
)

func Example() {
	r := publicecr.New()
	images, err := r.ListTags(context.TODO(), "appmesh", "aws-appmesh-envoy")
	if err != nil {
		log.Fatal(err)
	}

	// Print all the image names
	for _, image := range images {
		log.Println(image.Name)
	}
}
