package check

import stdCheck "github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"

var Store = make(map[string]Check)

type Check interface {
	Execute() stdCheck.CheckOutput
}
