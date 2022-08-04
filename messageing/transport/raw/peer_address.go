package raw

import (
	"github.com/galenliu/chip/lib"
	"net/netip"
)

const IPv6MulticastFlagTransient uint8 = 0x1
const IPv6MulticastFlagPrefix uint8 = 0x2

//--------IPv6------------------------
// 2001::/16  全球单播地址 Global Unicast
// FE80::/10  本地链路地址 Link Local Unicast
// FFxx::/8   组播地址    Multicast （IPv6无广播）

func Multicast(id lib.FabricId, gid lib.GroupId) netip.Addr {

	var scope uint8 = 0x05

	//lFlagsAndScope = (Ox01 << 4  =  0x10)  | 050  = 0x15
	var lFlagsAndScope uint8 = ((IPv6MulticastFlagTransient & 0xF) << 4) | (scope & 0xF)
	var lReserved uint8 = 0x0
	var prefixLength uint8 = 0x40

	var prefix uint64 = 0xfd00000000000000 | (uint64(id) >> 8 & 0x00ffffffffffffff)

	var groupId uint32 = ((uint32(id) << 24) & 0xff000000) | uint32(gid)

	ipV6 := netip.AddrFrom16([16]byte{
		0xFF, lFlagsAndScope, lReserved, prefixLength,
		byte(prefix >> 56), byte(prefix >> 48), byte(prefix >> 40), byte(prefix >> 32),
		byte(prefix >> 24), byte(prefix >> 16), byte(prefix >> 8), byte(prefix >> 0),
		byte(groupId >> 24), byte(groupId >> 16), byte(groupId >> 8), byte(groupId >> 0),
	})
	return ipV6
}
