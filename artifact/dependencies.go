package artifact

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"strings"

	"github.com/staffano/crazy-build/workspace"
)

// IgnoreStamps if true stamps folder will be ignored
var IgnoreStamps bool

// dependencies between commands
var dependencies map[string][]string

// isDone checks if a command already is executed
func isDone(cmd string) bool {
	if IgnoreStamps {
		return false
	}
	path := filepath.Join(workspace.GetStampDirPath(), cmd)
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// Mark a command as done in the stamp dir
func markDone(cmd string) {
	if IgnoreStamps {
		return
	}
	sd := workspace.GetStampDirPath()
	f := filepath.Join(sd, cmd)
	ff, err := os.OpenFile(f, os.O_RDONLY|os.O_CREATE, 0666)
	defer ff.Close()
	if err != nil {
		log.Fatalf("Error creating stamp file %s : %v", f, err)
	}
}

// Depends declares dependencies for arty (artifact:cmd)
// Depends("AMBuild.Compile", "AMBuild.Configure", "AMBuild.Verify")
func Depends(arty string, deps ...string) {
	_, ok := dependencies[arty]
	if !ok {
		dependencies[arty] = make([]string, 0, 20)
	}
	for _, s := range deps {
		dependencies[arty] = append(dependencies[arty], s)
	}
}

// Call the command by using introspection.
// Example artifact.Call("AMBuilder:Instantiate") will first call all dependencies
// then make sure used services are started and then call AMBuilder.Instantiate()
// The services are associated with the artifact using dependency injection
func Call(cmd string) {

	if isDone(cmd) {
		log.Printf("%s already done, skipping...", cmd)
		return
	}
	cs := strings.Split(cmd, ".")
	a := Find(cs[0])
	if len(a) == 0 {
		log.Fatal("Artifact not found")
	}
	cmds := GetCommands(*a[0])
	if len(cmds) == 0 {
		log.Fatal("No commands for artifact")
	}

	// First recursively call any non-complete dependencies
	deps, ok := dependencies[cmd]
	if ok {
		for _, dep := range deps {
			Call(dep)
		}
	}

	// Resolve Service dependencies
	artifactValue := reflect.ValueOf(*a[0]).Elem()
	t := artifactValue.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if req, ok := field.Tag.Lookup("requirement"); ok {
			for _, s := range services {
				// The pointer to the instance could Implement
				// the service interface. Why?
				serviceType := reflect.ValueOf(*s).Type()
				if serviceType.Implements(field.Type) {
					if (*s).Satisfies(req) && (*s).IsAvailable() {
						artifactValue.Field(i).Set(reflect.ValueOf(*s))
					}
				}
			}
		}
	}

	// Call cmd
	CallCmd(a[0], cs[1])

	// Mark cmd done
	markDone(cmd)

}

func init() {
	dependencies = make(map[string][]string)
	flag.BoolVar(&IgnoreStamps, "ignore-stamps", false, "Ignore stamps and force execution")
}
