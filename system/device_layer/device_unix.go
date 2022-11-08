//go:build unix || (js && wasm)

package DeviceLayer

import "github.com/galenliu/chip/system"

func SystemLayer() system.Layer {
	return nil
}
