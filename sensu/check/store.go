package check

import stdCheck "github.com/upfluence/sensu-go/sensu/check"

var Store = make(map[string]Check)

type Check interface {
	Execute() stdCheck.CheckOutput
}
