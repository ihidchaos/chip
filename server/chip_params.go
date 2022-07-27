package server

import (
	"github.com/galenliu/chip/access"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/storage"
	"net"
)

type SessionResumptionStorage interface {
}

type CommonCaseDeviceServerInitParams struct {
	ServerInitParams
}

type IgnoreCertificateValidityPolicy struct {
	credentials.CertificateValidityPolicy
}

type ServerInitParams struct {
	OperationalServicePort        uint16
	UserDirectedCommissioningPort uint16
	InterfaceId                   net.Interface
	AppDelegate                   AppDelegate

	// Persistent storage delegate: MUST be injected. Used to maintain storage by much common code.
	// Must be initialized before being provided.
	PersistentStorageDelegate storage.PersistentStorageDelegate
	// Session resumption storage: Optional. Support session resumption when provided.
	// Must be initialized before being provided.
	SessionResumptionStorage SessionResumptionStorage
	// Certificate validity policy: Optional. If none is injected, CHIPCert
	// enforces a default policy.

	CertificateValidityPolicy credentials.CertificateValidityPolicy

	// Group data provider: MUST be injected. Used to maintain critical keys such as the Identity
	// Protection Key (IPK) for CASE. Must be initialized before being provided.
	GroupDataProvider credentials.GroupDataProvider
	// Access control delegate: MUST be injected. Used to look up access control rules. Must be
	// initialized before being provided.
	AccessDelegate access.Delegate
	// ACL storage: MUST be injected. Used to store ACL entries in persistent storage. Must NOT
	// be initialized before being provided.
	//aclStorage app::AclStorage * aclStorage = nullptr;
	AclStorage AclStorage
	// Network native params can be injected depending on the
	// selected Endpoint implementation

	// Network native params can be injected depending on the
	// selected Endpoint implementation
	EndpointNativeParams func()

	// Optional. Support test event triggers when provided. Must be initialized before being
	// provided.
	TestEventTriggerDelegate TestEventTriggerDelegate
	// Operational keystore with access to the operational keys: MUST be injected.
	OperationalKeystore crypto.PersistentStorageOperationalKeystore
	// Operational certificate store with access to the operational certs in persisted storage:
	// must not be null at timne of Server::initCommissionableData().
	OpCertStore credentials.PersistentStorageOpCertStore
}

func NewServerInitParams() *ServerInitParams {
	return &ServerInitParams{}
}

func (p *ServerInitParams) Init(options *config.DeviceOptions) (*ServerInitParams, error) {
	p.OperationalServicePort = options.SecuredDevicePort
	p.UserDirectedCommissioningPort = options.UnsecuredCommissionerPort
	p.InterfaceId = options.InterfaceId
	return p, nil
}

func NewCommonCaseDeviceServerInitParams() *CommonCaseDeviceServerInitParams {
	c := &CommonCaseDeviceServerInitParams{
		ServerInitParams: ServerInitParams{
			OperationalKeystore:           nil,
			OperationalServicePort:        config.GetDeviceOptionsInstance().SecuredDevicePort,
			UserDirectedCommissioningPort: config.GetDeviceOptionsInstance().UnsecuredCommissionerPort,
			InterfaceId:                   config.GetDeviceOptionsInstance().InterfaceId,
		},
	}
	return c
}

func (p *ServerInitParams) InitializeStaticResourcesBeforeServerInit() error {

	var sKvsPersistentStorageDelegate storage.PersistentStorageDelegate
	var sPersistentStorageOperationalKeystore crypto.PersistentStorageOperationalKeystore
	var sPersistentStorageOpCertStore credentials.PersistentStorageOpCertStore
	var sGroupDataProvider credentials.GroupDataProvider
	var sDefaultCertValidityPolicy credentials.CertificateValidityPolicy

	if p.PersistentStorageDelegate == nil {
		sKvsPersistentStorageDelegate = storage.KeyValueStoreMgr()
		p.PersistentStorageDelegate = sKvsPersistentStorageDelegate
	}

	if p.OperationalKeystore == nil {
		sPersistentStorageOperationalKeystore = crypto.NewPersistentStorageOperationalKeystoreImpl()
		sPersistentStorageOperationalKeystore.Init(p.PersistentStorageDelegate)
	}
	if p.OpCertStore == nil {
		sPersistentStorageOpCertStore = credentials.NewPersistentStorageOpCertStoreImpl()
		sPersistentStorageOpCertStore.Init(p.PersistentStorageDelegate)
		p.OpCertStore = sPersistentStorageOpCertStore
	}

	sGroupDataProvider = credentials.NewGroupDataProviderImpl()
	sGroupDataProvider.SetStorageDelegate(p.PersistentStorageDelegate)
	err := sGroupDataProvider.Init()
	if err != nil {
		return err
	}
	p.GroupDataProvider = sGroupDataProvider

	{
		//TODO 根据配置 CHIP_CONFIG_ENABLE_SESSION_RESUMPTION 初始化
		p.SessionResumptionStorage = nil
	}

	p.AccessDelegate = access.GetAccessControlDelegate()

	{
		//TODO 未实现
		p.AclStorage = NewAclStorageImpl()
	}

	sDefaultCertValidityPolicy = IgnoreCertificateValidityPolicy{}
	p.CertificateValidityPolicy = sDefaultCertValidityPolicy

	return nil
}

func (p *CommonCaseDeviceServerInitParams) InitializeStaticResourcesBeforeServerInit() error {

	var sKvsPersistentStorageDelegate storage.PersistentStorageDelegate
	var sPersistentStorageOperationalKeystore crypto.PersistentStorageOperationalKeystore
	var sPersistentStorageOpCertStore credentials.PersistentStorageOpCertStore
	var sGroupDataProvider credentials.GroupDataProvider
	var sDefaultCertValidityPolicy credentials.CertificateValidityPolicy

	if p.PersistentStorageDelegate == nil {
		sKvsPersistentStorageDelegate = storage.KeyValueStoreMgr()
		p.PersistentStorageDelegate = sKvsPersistentStorageDelegate
	}

	if p.OperationalKeystore == nil {
		sPersistentStorageOperationalKeystore = crypto.NewPersistentStorageOperationalKeystoreImpl()
		sPersistentStorageOperationalKeystore.Init(p.PersistentStorageDelegate)
	}
	if p.OpCertStore == nil {
		sPersistentStorageOpCertStore = credentials.NewPersistentStorageOpCertStoreImpl()
		sPersistentStorageOpCertStore.Init(p.PersistentStorageDelegate)
		p.OpCertStore = sPersistentStorageOpCertStore
	}

	sGroupDataProvider = credentials.NewGroupDataProviderImpl()
	sGroupDataProvider.SetStorageDelegate(p.PersistentStorageDelegate)
	err := sGroupDataProvider.Init()
	if err != nil {
		return err
	}
	p.GroupDataProvider = sGroupDataProvider

	{
		//TODO 根据配置 CHIP_CONFIG_ENABLE_SESSION_RESUMPTION 初始化
		p.SessionResumptionStorage = nil
	}

	p.AccessDelegate = access.GetAccessControlDelegate()

	{
		//TODO 未实现
		p.AclStorage = NewAclStorageImpl()
	}

	sDefaultCertValidityPolicy = IgnoreCertificateValidityPolicy{}
	p.CertificateValidityPolicy = sDefaultCertValidityPolicy

	return nil
}
