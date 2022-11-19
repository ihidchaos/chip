package lib

type ScopedNodeId struct {
	mNodeId      NodeId
	mFabricIndex FabricIndex
}

func UndefinedScopedNodeId() *ScopedNodeId {
	return &ScopedNodeId{
		mNodeId:      UndefinedNodeId(),
		mFabricIndex: UndefinedFabricIndex(),
	}
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
	return s.mFabricIndex != UndefinedFabricIndex() && s.mNodeId.IsOperationalNodeId()
}
