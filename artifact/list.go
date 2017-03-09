package artifact

import (
	"log"
	"reflect"
	"strings"
)

var arties []Artifact

// GetAll artifacts
func GetAll() []Artifact {
	return arties
}

// Find artifacts given a search criteria
// type:id
// type
// id
func Find(cond string) []Artifact {
	ss := strings.Split(cond, ":")
	var res []Artifact
	switch parts := len(ss); parts {
	case 1:
		// first search on id then on type
		for _, a := range arties {
			if a.ID() == ss[0] {
				res = append(res, a)
			}
		}
		if len(res) > 0 {
			return res
		}
		for _, a := range arties {
			aType := reflect.TypeOf(a).String()
			if strings.Contains(aType, ss[0]) {
				res = append(res, a)
			}
		}
		return res
	case 2:
		for _, a := range arties {
			if reflect.TypeOf(a).String() == ss[0] &&
				a.ID() == ss[1] {
				res = append(res, a)
			}
		}
	default:
		log.Fatalf("Could not parse artifact search condition: %q", cond)
	}
	return nil
}

// Add one or more Artifacts
func Add(a ...Artifact) {
	arties = append(arties, a...)
}
