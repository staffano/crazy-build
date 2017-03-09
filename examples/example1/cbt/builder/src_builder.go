package builder

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/docker/go-connections/nat"
	"github.com/staffano/crazy-build/artifact"
	"github.com/staffano/crazy-build/workspace"
)

// AMBuilder is the artifact built from this directory
type AMBuilder struct {
	artifact.BaseArtifact
	dockerImage *artifact.DockerArtifact
}

// NewAMBuilder creates a new initialized instance of AMBuilder
func NewAMBuilder() *AMBuilder {
	builder := new(AMBuilder)
	builder.dockerImage = artifact.NewDockerArtifact()
	builder.dockerImage.SetID("ambuilder")

	// Build the docker image from the docker directory
	builder.dockerImage.ContextFolder = "${WORKSPACE}/cbt/docker"

	// Let the /src folder in the container hold the source code
	builder.dockerImage.Bindings = map[string]string{"${WORKSPACE}": "/src"}

	// Create a separate volume to hold build result
	builder.dockerImage.VolumeMap = []string{"build_vol:/build"}

	// Need a portmap for gdbserver
	builder.dockerImage.PortMap = nat.PortMap{
		"5555/tcp": []nat.PortBinding{
			{HostIP: "localhost",
				HostPort: "5555"}},
	}

	// Build the image
	builder.dockerImage.Build()

	return builder
}

// runCmd collects arguments and runs the command in the container
func (builder AMBuilder) runCmd(cmd ...string) {
	rargs := cmd
	for _, v := range builder.GetArgs()[1:] {
		rargs = append(rargs, v)
	}
	builder.dockerImage.Run(rargs...)
}

// Configure runs /src/configure [args] in the /build dir
func (builder AMBuilder) Configure(args ...string) {
	builder.dockerImage.WorkingDir = "/src"
	builder.runCmd("autoreconf", "--install")
	builder.dockerImage.WorkingDir = "/build"
	builder.runCmd("/src/configure", "--host=i686-w64-mingw32")
}

// Build ...
func (builder AMBuilder) Build(args ...string) {
	builder.runCmd("make", "-j8")
}

// Clean ...
func (builder AMBuilder) Clean(args ...string) {
	builder.runCmd("make", "distclean")
}

// Install ...
func (builder AMBuilder) Install(args ...string) {
	builder.runCmd("rm", "-rf", "/build/tmp/dist")
	builder.runCmd("make", "install", "DESTDIR=/build/tmp/dist")
	builder.runCmd("tar", "-C", "/build/tmp/dist", "-cvf", "hello_crazy_build-1.0.tar", ".")
	builder.runCmd("gzip", "-9f", "hello_crazy_build-1.0.tar")
	builder.runCmd("cp", "/build/hello_crazy_build-1.0.tar.gz", "/src/")
}

// Test ...
func (builder AMBuilder) Test(args ...string) {
	builder.Build()
	builder.runCmd("cp", "/build/src/hello.exe", "/src/")
	cmd := exec.Command(filepath.Clean(workspace.Resolve("${WORKSPACE}/hello.exe")))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n", out.String())
}

// Debug ...
func (builder AMBuilder) Debug(args ...string) {
	builder.dockerImage.SecOpts = []string{"seccomp=unconfined"}
	builder.runCmd("debug")
}
