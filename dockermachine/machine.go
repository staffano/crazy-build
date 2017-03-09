package dockermachine

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	normalLog "log"

	"github.com/docker/docker/client"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/drivers/virtualbox"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/state"
	"github.com/staffano/crazy-build/workspace"
)

// Store tells which persistent store to use
// for known docker machines
type Store int

const (
	// Default uses the docker machine store at
	// ${HOME}/.docker/machine
	Default Store = iota

	// Workspace creates and stores the docker
	// machine inside ${WORKSPACE}/.crazy_build/docker/machine
	Workspace
)

// SelectedStore is the store selected for using with Docker MAchines
var SelectedStore = Default

// DockerMachine represents an instance of a Docker Machine
// that we use for our docker engine.
// We store our machine configurations at ./ducker_dan/docker/machines
//

// DefaultMachineName the name of the default machine
const DefaultMachineName = "default"

// VirtualBoxDockerMachineV1v0v0 is an artifact that encapsulates
// virtualbox docker machine
type VirtualBoxDockerMachineV1v0v0 struct {
	SelectedStore Store
}

// Init ...
func (a VirtualBoxDockerMachineV1v0v0) Init() {
	// Initialize docker machines so we have something to execute
	// our containers on
	if !a.Exists(DefaultMachineName) {
		err := a.CreateDefaultMachine()
		if err != nil {
			normalLog.Fatalf("Error creating default docker machine: %q", err.Error())
		}
	}
	a.StartDockerMachine()
}

// GetBaseDir returns the base directory for our docker machine store
func (a VirtualBoxDockerMachineV1v0v0) GetBaseDir() string {
	switch a.SelectedStore {
	case Default:
		return mcndirs.GetBaseDir()
	case Workspace:
		wd, _ := workspace.Get("WORKSPACE")
		return filepath.Join(wd, workspace.WspConfigFolder, "docker", "machine")
	}
	return ""
}

// GetMachineCertDir returns the directory where we store certificates
// used to access our docker machines
func (a VirtualBoxDockerMachineV1v0v0) GetMachineCertDir() string {
	return filepath.Join(a.GetBaseDir(), "certs")
}

// CreateVirtualBoxDockerMachine creates a virtual box docker machine
// shareFolder follows the same syntax as the virtualbox driver option ShareFolder
func (a VirtualBoxDockerMachineV1v0v0) CreateVirtualBoxDockerMachine(name string, cpus int, memory int, shareFolder string) error {
	client := libmachine.NewClient(a.GetBaseDir(), a.GetMachineCertDir())
	defer client.Close()

	driver := virtualbox.NewDriver(name, a.GetBaseDir())
	driver.CPU = cpus
	driver.Memory = memory

	if shareFolder != "" {
		driver.ShareFolder = shareFolder
	}

	data, err := json.Marshal(driver)
	if err != nil {
		return err
	}

	h, err := client.NewHost("virtualbox", data)
	if err != nil {
		return err
	}

	if err := client.Create(h); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// CreateDefaultMachine creates a default machine
func (a VirtualBoxDockerMachineV1v0v0) CreateDefaultMachine() error {
	// DefaultMachineName the name of the default machine
	return a.CreateVirtualBoxDockerMachine(DefaultMachineName, 6, 6000, "")
}

// StartDockerMachine starts an existing docker machine
func (a VirtualBoxDockerMachineV1v0v0) StartDockerMachine(name ...string) error {
	var n string
	if len(name) == 0 {
		// DefaultMachineName the name of the default machineme) == 0 {
		n = DefaultMachineName
	} else {
		n = name[0]
	}
	client := libmachine.NewClient(a.GetBaseDir(), a.GetMachineCertDir())
	defer client.Close()
	host, err := client.Load(n)
	if err != nil {
		return err
	}

	s, err := host.Driver.GetState()
	if err != nil {
		log.Errorf("error: %s", err)
		return err
	}

	if s != state.Running {
		err = host.Start()
		if err != nil {
			log.Errorf("error: %s", err)
			return err
		}
	}

	return nil
}

// StopDockerMachine stops a machine
func (a VirtualBoxDockerMachineV1v0v0) StopDockerMachine(name ...string) error {
	var n string
	if len(name) == 0 {
		// DefaultMachineName the name of the default machineme) == 0 {
		n = DefaultMachineName
	} else {
		n = name[0]
	}
	client := libmachine.NewClient(a.GetBaseDir(), a.GetMachineCertDir())
	defer client.Close()
	host, err := client.Load(n)
	if err != nil {
		return err
	}
	return host.Stop()
}

// Exists checks if a machine exists with this name
func (a VirtualBoxDockerMachineV1v0v0) Exists(name string) bool {
	client := libmachine.NewClient(a.GetBaseDir(), a.GetMachineCertDir())
	defer client.Close()
	ex, err := client.Filestore.Exists(name)
	if err != nil {
		panic(err)
	}
	return ex
}

// CreateClient returns a valid API client to the docker engine
// running within the machine
func (a VirtualBoxDockerMachineV1v0v0) CreateClient(name ...string) (*client.Client, error) {
	lmClient := libmachine.NewClient(a.GetBaseDir(), a.GetMachineCertDir())
	defer lmClient.Close()
	var n string
	if len(name) == 0 {
		n = DefaultMachineName
	} else {
		n = name[0]
	}
	host, err := lmClient.Load(n)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	url, err := host.URL()
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	version := "v1.26"

	var client *http.Client
	options := tlsconfig.Options{
		CAFile:   filepath.Join(a.GetMachineCertDir(), "ca.pem"),
		CertFile: filepath.Join(a.GetMachineCertDir(), "cert.pem"),
		KeyFile:  filepath.Join(a.GetMachineCertDir(), "key.pem"),
	}
	tlsc, err := tlsconfig.Client(options)
	if err != nil {
		return nil, err
	}

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsc,
		},
	}

	s, err := dockerClient.NewClient(url, version, client, nil)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// MachineDriver returns the name of the driver
func (a VirtualBoxDockerMachineV1v0v0) MachineDriver(name ...string) string {
	client := libmachine.NewClient(a.GetBaseDir(), a.GetMachineCertDir())
	defer client.Close()
	var n string
	if len(name) == 0 {
		n = DefaultMachineName
	} else {
		n = name[0]
	}
	host, err := client.Load(n)
	if err != nil {
		normalLog.Fatal(err)
	}
	return host.DriverName
}
