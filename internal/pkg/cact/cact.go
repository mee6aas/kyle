package cact

import (
	"github.com/mee6aas/zeep/pkg/activity"
)

var (
	cact activity.Activity
)

// UnmarshalFromFile parses activity manifest from specified file.
func UnmarshalFromFile(act string) (e error) {
	cact, e = activity.UnmarshalFromFile(act)

	return
}

// HasDep checks if current activity has the dependency.
func HasDep() bool { return len(cact.Dependencies) > 0 }

// Dep returns dependency descriptor with specified name.
func Dep(name string) (dep activity.Dep, ok bool) {
	dep, ok = cact.Dependencies[name]

	return
}
