package cmd

import (
	"flag"
	"log"
	"os"

	"github.com/staffano/crazy-build/artifact"
)

// Flags

// VerboseFlag ...
var VerboseFlag bool

func init() {

	flag.BoolVar(&VerboseFlag, "verbose", false, "Set to true for more verbose output.")
}

// Command represents a command line command
type Command struct {
	ID    string
	Short string
	Long  string
	Cmd   func(args ...string)
}

var nativeCmds = []Command{
	{ID: "ls", Short: "List available artifacts", Cmd: func(args ...string) {}},
	{ID: "conf", Short: "Configure the build system", Cmd: func(args ...string) {}},
	{ID: "help", Short: "Show help", Cmd: func(args ...string) {}}}

// a == nil => glbal help
func showHelp(a *artifact.Artifact) {
	log.Printf("Help!")
}

// Execute ...
// cbt --workspace=232323 --verbose gcc build
// cbt [flags] [artifact] cmd [aux args]
// or...
// cbt [flags] [native command] aux args
func Execute() {

	if len(flag.Args()) == 0 {
		showHelp(nil)
		os.Exit(1)
	}

	// Check if the first argument is a native command
	for _, nc := range nativeCmds {
		if nc.ID == flag.Arg(0) {
			nc.Cmd(flag.Args()[1:]...)
			os.Exit(0)
		}
	}
	for _, a := range flag.Args() {
		artifact.Call(a)
	}
}

// build stages:
// 1. Declare - scope of artifacts (Compile time)
// 2. Configure - the artifacts registers their configuration interests
//    I want to configure the kernel (IKernelConfig)
// 3. Execute - Starting from the top artifact.
//    If an artifact owns a configuration, then it will call all artifacts
//    with a declared interest in configuring the artifact. The way to configure
//    is defined by the interface
// We need an artifact registry and a configuration registry.

// cbt --workspace=erer --wes=sd  build extra args
