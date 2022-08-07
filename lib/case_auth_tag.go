package lib

type CASEAuthTag uint32

const kUndefinedCAT CASEAuthTag = 0
const KTagIdentifierMask NodeId = 0x0000_0000_FFFF_0000
const KTagIdentifierShift uint32 = 16

const KTagVersionMask NodeId = 0x0000_0000_0000_FFFF

type CATValues struct {
}
