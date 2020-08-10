package controller

import (
	"testing"

	"github.com/Asutorufa/yuhaiin/subscr"
)

func TestLatency(t *testing.T) {
	x, err := subscr.GetNowNode()
	if err != nil {
		t.Error(err)
	}
	switch x.(type) {
	case *subscr.Shadowsocks:
		t.Log(Latency(x.(*subscr.Shadowsocks).NGroup, x.(*subscr.Shadowsocks).NName))
	}
}