package playlist

import (
	"os"
	"repertoire/server/test/integration/test/core"
	"testing"
)

func TestMain(m *testing.M) {
	ts := &core.TestServer{
		WithMeili: true,
	}
	ts.Start()

	code := m.Run()

	ts.Stop()
	os.Exit(code)
}
