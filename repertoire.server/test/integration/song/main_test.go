package song

import (
	"os"
	"repertoire/server/test/integration/test/core"
	"testing"
)

func TestMain(m *testing.M) {
	ts := core.Start()

	code := m.Run()

	core.Stop(ts)
	os.Exit(code)
}
