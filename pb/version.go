package pb

import "github.com/hashicorp/go-version"

var VERSION *version.Version
var CONSTRAINT version.Constraints

func init() {
	var err error
	VERSION, err = version.NewVersion("0.1.0")
	if err != nil {
		panic(err)
	}
	CONSTRAINT, err = version.NewConstraint(">= 0.1.0, < 1.0.0")
	if err != nil {
		panic(err)
	}
}
