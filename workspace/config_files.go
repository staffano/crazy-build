package workspace

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config is the model of the configuration file
type Config struct {
	Vars map[string]string `json:"env,omitempty"`
}

var variables map[string]string

var configuration *Config

// WspConfigFolder is the workspace config folder that marks the
// root of the workspace
const WspConfigFolder string = ".crazy_build"

// ConfigFile is the filename of the config file holding workspace
// configuration
const ConfigFile string = "config.json"

// StampDirName is the directory containing all stamps,
// which is markers that something has been done successfully
const StampDirName string = "stamps"

// WorkspaceRoot should be set to the root folder of the
// workspace
var WorkspaceRoot string

// GetWorkspaceRoot returns the root path of the workspace
func GetWorkspaceRoot(d ...string) string {
	if v, exist := variables["WORKSPACE"]; exist {
		WorkspaceRoot = v
		return v
	}

	if WorkspaceRoot != "" {
		return WorkspaceRoot
	}

	var cwd string
	if len(d) == 0 {
		cwd, _ = os.Getwd()
	} else {
		cwd = d[0]
	}

	for true {
		if cwd == "." || strings.HasSuffix(cwd, string(filepath.Separator)) {
			return ""
		}

		if _, err := os.Stat(filepath.Join(cwd, WspConfigFolder)); err == nil {
			WorkspaceRoot = cwd
			return cwd
		}
		cwd = filepath.Dir(cwd)
	}
	return ""
}

// GetStampDirPath returns the path within workspace config folder that
// contains stamp files
func GetStampDirPath() string {
	dir := filepath.Join(GetWorkspaceRoot(), WspConfigFolder, StampDirName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatalf("Error when creating %s: %v", dir, err)
		}
	}
	return filepath.Join(GetWorkspaceRoot(), WspConfigFolder, StampDirName)
}

// GetConfigFilePath returns the path to the config file within
// the workspace
func GetConfigFilePath() string {
	// GetWorkspaceRoot returns the root of the workspacefigFilePath() string {
	wr := GetWorkspaceRoot()
	if wr == "." {
		return ""
	}

	return filepath.Join(wr, WspConfigFolder, ConfigFile)
}

// Init initializes the environment package by loading variables from
// the project.json file
func Init() {
	// GetWorkspaceRoot returns the root of the workspace
	wspRoot := GetWorkspaceRoot()
	if wspRoot == "" {

		log.Fatalf("No %s directory found.", WspConfigFolder)
	}
	configuration = new(Config)
	configuration.Vars = make(map[string]string)
	variables = make(map[string]string)

	projectFile := filepath.Join(wspRoot, WspConfigFolder, "config.json")
	raw, err := os.Open(projectFile)
	defer raw.Close()
	if err != nil {
		panic(err)
	}
	json.NewDecoder(raw).Decode(configuration)

	for k, v := range configuration.Vars {
		variables[k] = v
	}

	// Set environment variables from os
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		variables[pair[0]] = pair[1]
	}

	// Set automatic variables
	variables["WORKSPACE"], _ = filepath.Abs(wspRoot)

}

// Get variable
func Get(k string) (string, bool) {
	v, err := variables[k]
	return v, err
}

// Configuration contains the config file
func Configuration() *Config {
	return configuration
}

// DumpConfig dumps the config to the console
func DumpConfig() {
	res2B, _ := json.Marshal(Configuration())
	fmt.Println(string(res2B))
}

// SetVar sets a variable
func SetVar(key, val string, perm bool) {
	variables[key] = Resolve(val)
	if perm {
		Configuration().Vars[key] = val
	}
}

// Resolve a string using the specified variables.
func Resolve(str string) string {
	// Brute force...
	for k, v := range variables {
		replVal := "${" + k + "}"
		str = strings.Replace(str, replVal, v, -1)
	}
	return str
}

// InitWorkspace initizliases a workspace at the
// specified location
func InitWorkspace(p string) error {
	wspDir := filepath.Join(p, WspConfigFolder)
	if _, err := os.Stat(wspDir); err == nil {
		return nil
	}
	err := os.Mkdir(wspDir, 0777)
	if err != nil {
		return err
	}
	configuration = new(Config)
	configuration.Vars = make(map[string]string)
	variables = make(map[string]string)
	variables["WORKSPACE"] = p
	return SaveConfig()
}

// SaveConfig stores the file .crazy_build/config.json
func SaveConfig() error {
	cfgFile := GetConfigFilePath()
	if cfgFile == "" {
		return errors.New("couldn't find workspace filepath")
	}
	file, err := os.Create(cfgFile)
	defer file.Close()

	if err != nil {
		return err
	}
	enc := json.NewEncoder(file)
	return enc.Encode(configuration)
}

func init() {
	flag.StringVar(&WorkspaceRoot, "workspace", "", "Set the workspace root path.")
}
