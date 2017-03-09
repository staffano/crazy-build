package artifact

import (
	"flag"
	"log"
)

// An Artifact is something we want out of the build system.
type Artifact interface {
	SetArgs([]string)
	GetArgs() []string
	Configure(args ...string)
	Build(args ...string)
	Clean(args ...string)
	Install(args ...string)
	Test(args ...string)
	Debug(args ...string)

	// The flagset this artifact uses. First the specific
	// artifact flags are parsed, then the general command
	// flags.
	Flags() *flag.FlagSet

	// Usage
	Usage() string
	ID() string
}

// A BaseArtifact the implements empty methods
type BaseArtifact struct {
	id   string
	args []string
}

// SetID sets the id
func (b *BaseArtifact) SetID(str string) {
	b.id = str
}

// SetArgs sets the specified arguments for the operation on the
// artifact
func (b *BaseArtifact) SetArgs(a []string) {
	b.args = a
}

// AddArg adds an argument
func (b *BaseArtifact) AddArg(a string) {
	b.args = append(b.args, a)
}

// GetArgs returns the argumenst set on operations on the
// artifact
func (b BaseArtifact) GetArgs() []string {
	return b.args
}

// Configure artifact
func (b BaseArtifact) Configure(args ...string) {
	log.Print("Configure not implemented.")
}

// Build artifact
func (b BaseArtifact) Build(args ...string) {
	log.Print("Build not implemented.")
}

// Clean artifact
func (b BaseArtifact) Clean(args ...string) {
	log.Print("Clean not implemented.")
}

// Install artifact
func (b BaseArtifact) Install(args ...string) {
	log.Print("Install not implemented.")
}

// Test artifact
func (b BaseArtifact) Test(args ...string) {
	log.Print("Test not implemented.")
}

// Debug artifact
func (b BaseArtifact) Debug(args ...string) {
	log.Print("Debug not implemented.")
}

// Flags to be used
func (b BaseArtifact) Flags() *flag.FlagSet {
	log.Print("Flags not implemented.")
	return nil
}

// Usage of the artifact
func (b BaseArtifact) Usage() string {
	log.Print("Usage not implemented.")
	return ""
}

// ID is the identity of the artifact
func (b BaseArtifact) ID() string {
	return b.id
}
