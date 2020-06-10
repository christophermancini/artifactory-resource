package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	resource "github.com/digitalocean/artifactory-resource"
	rlog "github.com/digitalocean/concourse-resource-library/log"
	"github.com/digitalocean/concourse-resource-library/metadata"
	jlog "github.com/jfrog/jfrog-client-go/utils/log"
)

func main() {
	input := rlog.WriteStdin()
	defer rlog.Close()

	jlog.SetLogger(jlog.NewLogger(jlog.DEBUG, log.Writer()))

	var request resource.GetRequest
	err := request.Read(input)
	if err != nil {
		log.Fatalf("failed to read request input: %s", err)
	}

	err = request.Source.Validate()
	if err != nil {
		log.Fatalf("invalid source config: %s", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("missing arguments")
	}
	dir := os.Args[1]

	response, err := resource.Get(request, dir)
	if err != nil {
		log.Fatalf("failed to perform get: %s", err)
	}

	// write metadata to output dir
	err = writeMetadataFile(response.Metadata, dir)
	if err != nil {
		log.Fatalf("failed to write metadata.json: %s", err)
	}

	err = response.Write()
	if err != nil {
		log.Fatalf("failed to write response to stdout: %s", err)
	}

	log.Println("Get complete")
}

func writeMetadataFile(m metadata.Metadata, dir string) error {
	data, err := m.JSON()
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Join(dir, "resource"), os.ModePerm)
	err = ioutil.WriteFile(filepath.Join(dir, "resource", "metadata.json"), data, 0644)
	if err != nil {
		return err
	}

	return nil
}
