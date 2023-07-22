package reality

import (
	"context"
	"testing"

	proxy "github.com/Asutorufa/yuhaiin/pkg/net/interfaces"
	"github.com/Asutorufa/yuhaiin/pkg/net/proxy/simple"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/protocol"
	"github.com/Asutorufa/yuhaiin/pkg/utils/assert"
)

func TestClient(t *testing.T) {
	sm := simple.New(&protocol.Protocol_Simple{
		Simple: &protocol.Simple{
			Host: "127.0.0.1",
			Port: 8888,
		},
	})

	c := NewRealityClient(
		&protocol.Protocol_Reality{
			Reality: &protocol.Reality{
				ServerName: "www.baidu.com",
				ShortId:    "",
				PublicKey:  "SOW7P-17ibm_-kz-QUQwGGyitSbsa5wOmRGAigGvDH8",
			},
		},
	)

	pp, err := sm(nil)
	assert.NoError(t, err)
	pp, err = c(pp)
	assert.NoError(t, err)

	conn, err := pp.Conn(context.Background(), proxy.EmptyAddr)
	assert.NoError(t, err)
	defer conn.Close()

	conn.Write([]byte("aaa"))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	assert.NoError(t, err)

	t.Log(string(buf[:n]))
}