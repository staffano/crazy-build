package artifacts

import (
	"github.com/staffano/crazy-build/artifact"
	"github.com/staffano/crazy-build/examples/example2/build/services"
)

// PrintArtifact is the artifact built from this directory
type PrintArtifact struct {
	artifact.BaseArtifact
	S1 services.Service1API `requirement:"target=x86_64-pc-linux-gnu, host=mipsel-unknown-linux,"`
	S2 services.Service2API `requirement:"target=x86_64-pc-linux-gnu, host=mipsel-unknown-linux,"`
}

// Print ...
func (a *PrintArtifact) Print() {
	a.S1.PrintHello()

}

// Print2 ...
func (a *PrintArtifact) Print2() {
	a.S2.PrintWorld()
}

func init() {
	artifact.Depends("PrintArtifact.Print2", "PrintArtifact.Print")
}
