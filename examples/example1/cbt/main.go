package main

import (
	"log"

	"github.com/staffano/crazy-build/artifact"
	"github.com/staffano/crazy-build/cmd"
	"github.com/staffano/crazy-build/examples/example1/cbt/builder"
)

// PrintUsage formats and prints usage
func printUsage() {
	log.Print("Usage: cbt cmd AMBuilder [opts] \n Possible the most complicated way of building hello world...\n")
}

func main() {
	cmd.Init()
	artifact.Add(builder.NewAMBuilder())
	cmd.HandleCmd(printUsage)
}
