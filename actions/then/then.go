package then

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/byuoitav/common/nerr"
	"go.uber.org/zap"
)

// Then represents something to be done as a result of all of an action's If checks passing
type Then struct {
	Do   string          `json:"do"`
	With json.RawMessage `json:"with"`
}

// Func .
type Func func(ctx context.Context, with []byte, log *zap.SugaredLogger) *nerr.E

var (
	thens struct {
		sync.RWMutex
		m map[string]Func
	}
)

func init() {
	thens.Lock()
	thens.m = make(map[string]Func)

	thens.m["skip"] = func(ctx context.Context, with []byte, log *zap.SugaredLogger) *nerr.E {
		return nil
	}

	// declare then's here
	thens.m["send-email"] = SendEmail
	thens.m["send-slack"] = SendSlack
	thens.m["reboot-pi"] = RebootPi

	thens.Unlock()
}

// Add .
func Add(name string, f Func) {
	// TODO check if the function already exists
	thens.Lock()
	thens.m[name] = f
	thens.Unlock()
}

// Get .
func Get(name string) Func {
	thens.RLock()
	defer thens.RUnlock()

	return thens.m[name]
}

// Execute executes a then
func (t *Then) Execute(ctx context.Context, log *zap.SugaredLogger) *nerr.E {
	// get the function from the thens
	f := Get(t.Do)
	if f == nil {
		return nerr.Createf("not-found", "no then function found with the name '%s'. make sure it's been added to the then map", t.Do)
	}

	err := f(ctx, t.With, log)
	if err != nil {
		return err.Addf("something went wrong running then '%s'", t.Do)
	}

	return nil
}
