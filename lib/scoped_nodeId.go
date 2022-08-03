package lib

var UndefinedScopedNodeId = ScopedNodeId{
	NodeId:      UndefinedNodeId,
	FabricIndex: UndefinedFabricIndex,
}

type ScopedNodeId struct {
	NodeId      NodeId
	FabricIndex FabricIndex
}

func (s ScopedNodeId) GetNodeId() NodeId {
	return s.NodeId
}

func (s ScopedNodeId) GetFabricIndex() FabricIndex {
	return s.FabricIndex
}

func (s ScopedNodeId) IsOperational() bool {
	return s.FabricIndex != UndefinedFabricIndex && IsOperationalNodeId(s.NodeId)
}
