//go:build unix || (js && wasm)

package DeviceLayer

import (
	"github.com/galenliu/chip/platform/system"
)

func SystemLayer() system.Layer {
	return nil
}
