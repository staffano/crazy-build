package artifact

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var arties []Artifact

var configurations map[string][]*Artifact

// IWantToConfigure tells the system that the self artifact
// has an interest in the shape of this configuration
func IWantToConfigure(conf string, self *Artifact) error {
	configurations[conf] = append(configurations[conf], self)
	return nil
}

// WhoWantToConfigure returns who is interested in the
// specified configuration
func WhoWantToConfigure(conf string) []*Artifact {
	return configurations[conf]
}

// GetAll artifacts
func GetAll() []Artifact {
	return arties
}

func lowName(a Artifact) string {
	ts := strings.ToLower(reflect.ValueOf(a).Type().String())
	return strings.Split(ts, ".")[1]
}

// Version ...
type Version struct {
	Major int
	Minor int
	Micro int
	Git   string
}

// NullVersion ...
var NullVersion = Version{Major: 0, Minor: 0, Micro: 0, Git: ""}

var versionRegExp = regexp.MustCompile(`^[vV]([0-9]+)[vV]*([0-9]*)[vV]*([0-9]*)`)
var gitRegExp = regexp.MustCompile(`^[vV]git[vV]([a-z0-9]*)[vV]`)

func getVersion(str string) Version {
	vre := versionRegExp.FindStringSubmatch(str)
	if vre != nil {
		minor := 0
		micro := 0
		major, _ := strconv.Atoi(vre[1])
		if len(vre) > 2 {
			minor, _ = strconv.Atoi(vre[2])
		}
		if len(vre) > 3 {
			micro, _ = strconv.Atoi(vre[3])
		}
		return Version{Major: major, Minor: minor, Micro: micro, Git: ""}
	}
	gre := gitRegExp.FindStringSubmatch(str)
	if gre != nil {
		return Version{Major: 0, Minor: 0, Micro: 0, Git: gre[1]}
	}
	return NullVersion
}

// Find artifacts given a search criteria
// gcc would match GccV1, gcc, GCC, GCCV2
// but not gccser
func Find(cond string) []*Artifact {
	condLow := strings.ToLower(cond)
	var res []*Artifact
	for _, a := range arties {
		aType := lowName(a)
		// Check if prefix match
		if !strings.HasPrefix(aType, condLow) {
			continue
		}
		// Check for exact match
		if aType == condLow {
			res = append(res, &a)
			continue
		}
		// Check if versioned
		rest := strings.TrimPrefix(aType, condLow)
		if getVersion(rest) != NullVersion {
			res = append(res, &a)
		}
	}
	return res
}

// Add one or more Artifacts
func Add(a ...Artifact) {
	arties = append(arties, a...)
}
