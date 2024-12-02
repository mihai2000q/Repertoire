package section

import (
	"os"
	"repertoire/server/test/integration/test/core"
	"testing"
)

func TestMain(m *testing.M) {
	ts := core.Start("../../../../.env")

	code := m.Run()

	core.Stop(ts)
	os.Exit(code)
}
