// +build ignore

package artifact

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"archive/tar"

	"io/ioutil"

	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/staffano/crazy-build/dockermachine"
	"github.com/staffano/crazy-build/workspace"
)

// VolumeMap specifies a mapping between a Host folder and a Image folder
type VolumeMap struct {
	HostPath, ImagePath string
}

// A DockerArtifact is the docker container created
// from a docker image. The image and container is built
// with the Build() method. It is executed with the
// Run() method. Clean() will remove the image and the
// container. Install() will push the image to the
// predefined repository
// We only allow whats called 'bind-mounts' in docker.
type DockerArtifact struct {
	BaseArtifact
	DockerFile     string            // url to docker file
	ContextFolder  string            // Folder that is used to create the image
	ContainerID    string            // Docker container id
	ImageID        string            // The id of the docker image
	Bindings       map[string]string // Mappings between host and container paths
	VolumeMap      []string          // Volume mappings like "build_vol:/build:ro"
	PortMap        nat.PortMap       // Port maps like 2223/tcp:2323
	SecOpts        []string          // Security options like seccomp=unconfined
	WorkingDir     string            // The current working dir the command will be executed in
	SuppressOutput bool              // Print less or more?
	isBuilt        bool              // IS the image already?
}

// NewDockerArtifact returns a new instance
func NewDockerArtifact() *DockerArtifact {
	return new(DockerArtifact)
}

func addFile(tw *tar.Writer, basepath, path string) error {
	log.Printf("Adding file %q to docker context", path)
	file, err := os.Open(filepath.Join(basepath, path))
	if err != nil {
		return err
	}
	defer file.Close()
	if stat, err := file.Stat(); err == nil {
		header := new(tar.Header)
		header.Name = path
		header.Size = stat.Size()
		header.Mode = 0770
		header.ModTime = stat.ModTime()

		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// copy file data into tar writer
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
	}
	return nil
}

func addDir(tw *tar.Writer, basepath, path string) error {
	log.Printf("Adding directory %q to docker context", path)
	files, _ := ioutil.ReadDir(filepath.Join(basepath, path))
	for _, f := range files {
		if f.IsDir() {
			addDir(tw, basepath, filepath.Join(path, f.Name()))
		} else {
			addFile(tw, basepath, filepath.Join(path, f.Name()))
		}
	}
	return nil
}

func createDockerCtxt(path string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	// Create a new tar archive.
	tw := tar.NewWriter(buf)
	defer tw.Close()
	err := addDir(tw, path, "")
	if err != nil {
		panic(err)
	}
	return buf, nil
}

// Giveh a volume description, return a volume mount
func getVolume(volDescr string) mount.Mount {
	parts := strings.Split(volDescr, ":")
	src := parts[0]
	target := parts[1]
	readOnly := false
	if len(parts) == 3 {
		if parts[2] == "ro" {
			readOnly = true
		}
	}
	return mount.Mount{
		Type:     "volume",
		Source:   src,
		Target:   target,
		ReadOnly: readOnly,
	}
}

func fixPath(p string) string {
	// Fix bind mounts on windows with machine driver virtualbox
	if runtime.GOOS == "windows" && dockermachine.MachineDriver() == "virtualbox" {
		q := filepath.ToSlash(filepath.Clean(p))
		l := strings.Split(q, ":")
		switch {
		case len(l) == 1:
			return q
		case len(l) == 2:
			return "//" + strings.ToLower(l[0][0:1]) + l[1]
		case len(l) > 2:
			log.Fatalf("Invalid path used: %s", p)
		}
	}
	return ""
}

// Build a docker image or load it from repository
func (d *DockerArtifact) Build(args ...string) {
	if d.isBuilt {
		return
	}
	ctx := context.Background()
	cli, err := dockermachine.CreateClient()
	if err != nil {
		panic(err)
	}
	buildOptions := types.ImageBuildOptions{
		Tags:           []string{d.ID()},
		ForceRemove:    true,
		SuppressOutput: d.SuppressOutput,
	}
	buildCtx, _ := createDockerCtxt(workspace.Resolve(d.ContextFolder))
	buildImageResponse, err := cli.ImageBuild(ctx, buildCtx, buildOptions)
	if err != nil {
		log.Fatalf("ImageBuild error %s", err)
	}
	defer buildImageResponse.Body.Close()
	// The body contains an output stream of the build result, send it to stdout
	err = jsonmessage.DisplayJSONMessagesStream(buildImageResponse.Body, os.Stdout, os.Stdout.Fd(), true, nil)
	if err != nil {
		log.Fatalf("DisplayJSONMessagesStream error %s", err)
	}
	d.isBuilt = true
}

func getCid(c *client.Client, name string) (string, error) {
	ctx := context.Background()
	containers, err := c.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Fatalf("ContainerList error %s", err)
		return "", err
	}
	for _, cont := range containers {
		for _, n := range cont.Names {
			if n == "/"+name {
				return cont.ID, nil
			}
		}
	}
	return "", nil
}

func removeContainer(c *client.Client, name string) error {
	ctx := context.Background()
	cid, err := getCid(c, name)
	if cid != "" {
		// Remove the container that has our name
		timeout := 500 * time.Millisecond
		if err = c.ContainerStop(ctx, cid, &timeout); err != nil {
			log.Fatalf("ContainerStop returns %s", err)
		}

		if err = c.ContainerRemove(ctx, cid, types.ContainerRemoveOptions{Force: true}); err != nil {
			log.Fatalf("ContainerRemove returns %s", err)
		}
	}
	return nil
}

// Run executes the docker container that was created in the Build method
func (d *DockerArtifact) Run(args ...string) {
	d.Build()
	ctx := context.Background()
	cli, err := dockermachine.CreateClient()
	if err != nil {
		log.Fatalf("CreateClient error %s", err)
	}
	// Make sure container does not alread exist
	if err := removeContainer(cli, d.ID()); err != nil {
		log.Fatalf("removeContainer error %s", err)
	}

	config := container.Config{
		Cmd:        args,
		Tty:        true,
		WorkingDir: d.WorkingDir,
	}

	config.Image = d.ID()
	hostConfig := container.HostConfig{}

	for _, so := range d.SecOpts {
		hostConfig.SecurityOpt = append(hostConfig.SecurityOpt, so)
	}

	// Bind mounts
	for k, v := range d.Bindings {
		hostConfig.Binds = append(hostConfig.Binds, fixPath(workspace.Resolve(k))+":"+v)
	}

	// Set volume mounts
	for _, v := range d.VolumeMap {
		hostConfig.Mounts = append(hostConfig.Mounts, getVolume(v))
	}

	// Network port bindings
	hostConfig.PortBindings = d.PortMap
	networkConfig := network.NetworkingConfig{}
	hostConfig.LogConfig = container.LogConfig{Type: "json-file", Config: map[string]string{}}
	// Create container
	buildContainerResponse, err := cli.ContainerCreate(ctx, &config, &hostConfig, &networkConfig, d.ID())
	if err != nil {
		log.Fatalf("ContainerCreate error %s", err)
	}
	d.ContainerID = buildContainerResponse.ID
	log.Printf("Docker Container created %s: %s", d.ID(), d.ContainerID)

	if err := cli.ContainerStart(ctx, d.ContainerID, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("ContainerStart error %s", err)
	}

	out, err := cli.ContainerLogs(ctx, d.ContainerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	defer out.Close()

	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		log.Fatalf("io.Copy error %s", err)
	}

	if _, err = cli.ContainerWait(ctx, d.ContainerID); err != nil {
		log.Fatalf("ContainerWait error %s", err)
	}

	if err != nil {
		log.Fatal(err.Error())
	}
	timeout := 500 * time.Millisecond
	if err = cli.ContainerStop(ctx, d.ContainerID, &timeout); err != nil {
		log.Fatalf("ContainerStop returns %s", err)
	}

	if err = cli.ContainerRemove(ctx, d.ContainerID, types.ContainerRemoveOptions{Force: true}); err != nil {
		log.Fatalf("ContainerRemove returns %s", err)
	}
}
