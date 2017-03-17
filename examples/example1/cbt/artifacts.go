package main

import (
	"github.com/staffano/crazy-build/artifact"
	"github.com/staffano/crazy-build/examples/example1/cbt/builder"
)

// LoadArtifacts will load all artifacts into the database
// This is where any new artifacts is entered.
func LoadArtifacts() {
 
	artifact.Add(new(builder.AMBuilder))
}
