package lib

type ScopedNodeId struct {
	mNodeId      NodeId
	mFabricIndex FabricIndex
}

var UndefinedScopedNodeId = &ScopedNodeId{
	mNodeId:      UndefinedNodeId(),
	mFabricIndex: FabricIndexUndefined,
}

func NewScopedNodeId(id NodeId, index FabricIndex) ScopedNodeId {
	return ScopedNodeId{
		mNodeId:      id,
		mFabricIndex: index,
	}
}

func (s ScopedNodeId) NodeId() NodeId {
	return s.mNodeId
}

func (s ScopedNodeId) FabricIndex() FabricIndex {
	return s.mFabricIndex
}

func (s ScopedNodeId) IsOperational() bool {
	return s.mFabricIndex != FabricIndexUndefined && s.mNodeId.IsOperationalNodeId()
}
