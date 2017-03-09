package artifact

import (
	"log"
	"reflect"
)

// An Artifact identifies a product of the build system. It can be
// in two states: Loaded and Instatiated.
// First all artifacts are loaded and then, when needed, they
// are instantiated. It is assumed loading is fast.
type Artifact interface {
	// Usage
	Usage() string
	ID() string

	CheckConfiguration()
}

// A BaseArtifact the implements empty methods
type BaseArtifact struct {
	id             string
	args           []string
	isInstantiated bool
}

// SetID sets the id
func (b *BaseArtifact) SetID(str string) {
	b.id = str
}

// Usage of the artifact
func (b *BaseArtifact) Usage() string {
	log.Print("Usage not implemented.")
	return ""
}

// ID is the identity of the artifact
func (b *BaseArtifact) ID() string {
	return b.id
}

// CheckConfiguration lets the artifact register interest
// in configurations
func (b *BaseArtifact) CheckConfiguration() {
}

// GetCommands returns a list of valid commands on the artifact
func GetCommands(a Artifact) []string {
	var res []string
	val := reflect.ValueOf(a).Type()
	for n := 0; n < val.NumMethod(); n++ {
		res = append(res, val.Method(n).Name)
	}
	return res
}

// CallCmd calls a.meth(args)
func CallCmd(a *Artifact, meth string, args ...string) {
	inputs := make([]reflect.Value, 1)
	inputs[0] = reflect.ValueOf(args)
	m := reflect.ValueOf(*a).MethodByName(meth)
	m.Call(nil)
	//	reflect.ValueOf(*a).MethodByName(meth).CallSlice(inputs)
}

// RegisterConfigurationInterest will loop through all artifacts and let them
// tell what configurations they are interested in.
func RegisterConfigurationInterest() {
	for _, a := range arties {
		a.CheckConfiguration()
	}
}
