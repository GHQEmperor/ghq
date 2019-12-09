package watchdogs

import (
	"net/rpc"
	"testing"
)

func TestAddContainer(t *testing.T) {
	if err := rpc.Register(&ContMap); err != nil {
		panic(err)
	}
}
