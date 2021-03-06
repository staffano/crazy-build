// +build ignore

package main

import (
	"flag"

	"github.com/staffano/crazy-build/artifact"
	"github.com/staffano/crazy-build/cmd"
	"github.com/staffano/crazy-build/workspace"
)

// Handling of configurations really needs two passes
// First pass to register interest, second pass to
// execute. An artifact that has registered intereset
// but is excluded should not be allowed to modify the
// configuration? Or?

// Instantiation is handled as a separate cmd

func main() {
	flag.Parse()
	workspace.Init()
	LoadArtifacts()
	artifact.RegisterConfigurationInterest()
	cmd.Execute()
}
