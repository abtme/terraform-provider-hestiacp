package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/your-org/terraform-provider-hestiacp/internal/provider"
)

// Run "go generate" to regenerate docs.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	// Set by goreleaser at build time.
	version string = "dev"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "start provider in debug mode for use with delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/your-org/hestiacp",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
