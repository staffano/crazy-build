package main

import (
	"flag"
	"log"
	"os"

	"github.com/staffano/crazy-build/artifact"
	"github.com/staffano/crazy-build/cmd"
	"github.com/staffano/crazy-build/examples/example2/build/artifacts"
	"github.com/staffano/crazy-build/workspace"
)

// Handling of configurations really needs two passes
// First pass to register interest, second pass to
// execute. An artifact that has registered intereset
// but is excluded should not be allowed to modify the
// configuration? Or?

// Instantiation is handled as a separate cmd

func main() {
	log.Printf("%v", os.Args)
	flag.Parse()
	workspace.Init()
	artifact.Add(new(artifacts.PrintArtifact))
	artifact.RegisterConfigurationInterest()
	cmd.Execute()
}
