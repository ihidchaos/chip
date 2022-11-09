package config

import (
	"github.com/galenliu/chip/platform/path"
	"testing"
)

func TestFlags(t *testing.T) {
	t.Log(path.GetFatConFile())
	t.Log(path.GetSysConFile())
	t.Log(path.GetLocalConFile())
}
