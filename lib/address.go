package lib

import (
	"net/netip"
)

const IPv6MulticastFlagTransient uint8 = 0x30
const IPv6MulticastFlagPrefix uint8 = 0x20
const IPv6MulticastPort = 5540

//--------IPv6------------------------
// 2001::/16  全球单播地址 Global Unicast
// FE80::/10  本地链路地址 Link Local Unicast
// FFxx::/8   组播地址    Multicast （IPv6无广播）

func Multicast(id FabricId, gid GroupId) netip.Addr {
	return MakeIPv6TransientMulticast(id, gid)
}

func MakeIPv6TransientMulticast(fabricId FabricId, gid GroupId) netip.Addr {
	var siteLocal uint8 = 0x05
	return IPV6Multicast(IPv6MulticastFlagTransient, siteLocal, fabricId, gid)
}

func IPV6Multicast(flag, scope uint8, fid FabricId, gid GroupId) netip.Addr {

	var lFlagsAndScope = flag & scope
	var lReserved uint8 = 0x00

	var prefix = 0xfd00000000000000 | (uint64(fid) >> 8 & 0x00ffffffffffffff)

	ipV6 := netip.AddrFrom16([16]byte{
		0xFF, lFlagsAndScope, lReserved, 0x40,
		byte(prefix >> 56), byte(prefix >> 48), byte(prefix >> 40), byte(prefix >> 32),
		byte(prefix >> 24), byte(prefix >> 16), byte(prefix >> 8), byte(prefix >> 0),
		byte((fid & 0x000000FF) >> 0), 0x00, byte((gid & 0xFF00) >> 8), byte((gid & 0xFF) >> 0),
	})
	return ipV6
}

func MakeIPv6PrefixMulticast(scope uint8, prefixLength uint8, prefix uint64, gid GroupId) netip.Addr {
	return netip.AddrFrom16([16]byte{
		0xFF,
	})
}
