package lib

type PasscodeId uint16

const DefaultCommissioningPasscodeId PasscodeId = 0

func (p PasscodeId) NodeId() NodeId {
	return NodeId(uint64(kMinPAKEKeyId) | uint64(p))
}
