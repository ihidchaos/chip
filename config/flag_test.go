package config

import (
	"github.com/galenliu/chip/system"
	"testing"
)

func TestFlags(t *testing.T) {

	t.Log(system.GetFatConFile())
	t.Log(system.GetSysConFile())
	t.Log(system.GetLocalConFile())
}
