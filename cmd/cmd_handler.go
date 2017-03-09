// Copyright © 2017 Staffan Olsson <staffan@diversum.nu>

// Add artifacts to be controlled using
// artifact.Add
// Then call cmd.HandleCmd() to handle commands

package cmd

import (
	"log"
	"os"
	"reflect"

	"axis.com/dawn/go/docker_dawn/dockermachine"

	"github.com/staffano/crazy-build/artifact"
	"github.com/staffano/crazy-build/workspace"
)

// ListArtifacts formats and prints out a list of artifacts
func ListArtifacts(artis []artifact.Artifact) {
	log.Printf("listArtifacts count=%d", len(artis))
	for _, a := range artis {
		log.Printf("%q", reflect.TypeOf(a))
	}
}

// Init ialize workspace and environment
func Init() {
	workspace.Init()

	// Initialize docker machines so we have something to execute
	// our containers on
	if !dockermachine.Exists(dockermachine.DefaultMachineName) {
		err := dockermachine.CreateDefaultMachine()
		if err != nil {
			log.Fatalf("Error creating default docker machine: %q", err.Error())
		}
	}

	dockermachine.StartDockerMachine()

}

// HandleCmd handles commandline
func HandleCmd(printUsage func()) {
	var artySelection string
	var flagArgs []string

	switch argCount := len(os.Args); argCount {
	case 0:
		printUsage()
	case 1:
		printUsage()
	case 2:
		// No object specified, we have to use "default"
		artySelection = "defaultType:default"
		flagArgs = os.Args[1:]
	default:
		artySelection = os.Args[2]
		flagArgs = os.Args[2:]
	}

	arty := artifact.Find(artySelection)

	// Set the arguments
	for _, a := range arty {
		a.SetArgs(flagArgs)
	}

	switch os.Args[1] {
	case "configure":
		for _, a := range arty {
			a.Configure()
		}
	case "build":
		for _, a := range arty {
			a.Build()
		}
	case "clean":
		for _, a := range arty {
			a.Clean()
		}
	case "install":
		for _, a := range arty {
			a.Install()
		}
	case "test":
		for _, a := range arty {
			a.Test()
		}
	case "ls":
		if len(os.Args) == 2 {
			ListArtifacts(artifact.GetAll())
		} else {
			ListArtifacts(arty)
		}
	case "debug":
		for _, a := range arty {
			a.Debug()
		}
	case "init":
		wd, _ := os.Getwd()
		err := workspace.InitWorkspace(wd)
		if err != nil {
			log.Fatalf("Error: %q", err.Error())
		}
	default:
		printUsage()
		os.Exit(1)
	}
}
