package dnssd

import (
	"github.com/galenliu/chip/device"
)

const (
	DiscoveryFilterType_kNone = iota
	DiscoveryFilterType_kShortDiscriminator
	DiscoveryFilterType_kLongDiscriminator
	DiscoveryFilterType_kVendorId
	DiscoveryFilterType_kDeviceType
	DiscoveryFilterType_kCommissioningMode
	DiscoveryFilterType_kInstanceName
	DiscoveryFilterType_kCommissioner
	DiscoveryFilterType_kCompressedFabricId
)

type CommissioningResolveDelegate interface {
}

type OperationalResolveDelegate interface {
}

type DiscoveryFilter interface {
}

type Resolver interface {
	Init() error
	Shutdown()

	ResolveNodeId(peerId device.PeerId, isIpV6 bool)

	SetOperationalDelegate(delegate OperationalResolveDelegate)
	SetCommissioningDelegate(delegate CommissioningResolveDelegate)

	DiscoverCommissionableNodes(filter DiscoveryFilter)
	DiscoverCommissioners(filter DiscoveryFilter)
}

type MinMdnsResolver struct {
	mOperationalDelegate   OperationalResolveDelegate
	mCommissioningDelegate CommissioningResolveDelegate
	//ActiveResolveAttempts mActiveResolves;
	//PacketParser mPacketParser;
}

func (m MinMdnsResolver) Init() error {
	//TODO implement me
	panic("implement me")
}

func (m MinMdnsResolver) Shutdown() {
	//TODO implement me
	panic("implement me")
}

func (m MinMdnsResolver) ResolveNodeId(peerId device.PeerId, isIpV6 bool) {
	//TODO implement me
	panic("implement me")
}

func (m MinMdnsResolver) SetOperationalDelegate(delegate OperationalResolveDelegate) {
	m.mOperationalDelegate = delegate
}

func (m MinMdnsResolver) SetCommissioningDelegate(delegate CommissioningResolveDelegate) {
	m.mCommissioningDelegate = delegate
}

func (m MinMdnsResolver) DiscoverCommissionableNodes(filter DiscoveryFilter) {
	//TODO implement me
	panic("implement me")
}

func (m MinMdnsResolver) DiscoverCommissioners(filter DiscoveryFilter) {
	//TODO implement me
	panic("implement me")
}
