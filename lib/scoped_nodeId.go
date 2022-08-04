package lib

type ScopedNodeId struct {
	NodeId      NodeId
	FabricIndex FabricIndex
}

func NewScopedNodeId() *ScopedNodeId {
	return &ScopedNodeId{
		NodeId:      KUndefinedNodeId,
		FabricIndex: UndefinedFabricIndex,
	}
}

func (s ScopedNodeId) GetNodeId() NodeId {
	return s.NodeId
}

func (s ScopedNodeId) GetFabricIndex() FabricIndex {
	return s.FabricIndex
}

func (s ScopedNodeId) IsOperational() bool {
	return s.FabricIndex != UndefinedFabricIndex && s.NodeId.IsOperationalNodeId()
}
