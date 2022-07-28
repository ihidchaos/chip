package chip

import (
	"github.com/galenliu/chip/access"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials"
	storage2 "github.com/galenliu/chip/crypto/persistent_storage"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/server"
	"github.com/galenliu/chip/storage"
	"net"
)

type CommonCaseDeviceServerInitParams struct {
	InitParams
}

type IgnoreCertificateValidityPolicy struct {
	credentials.CertificateValidityPolicy
}

func NewIgnoreCertificateValidityPolicy() *IgnoreCertificateValidityPolicy {
	return &IgnoreCertificateValidityPolicy{}
}

type InitParams struct {
	OperationalServicePort        uint16
	UserDirectedCommissioningPort uint16
	InterfaceId                   net.Interface
	AppDelegate                   AppDelegate

	// Persistent storage delegate: MUST be injected. Used to maintain storage by much common code.
	// Must be initialized before being provided.
	PersistentStorageDelegate storage.PersistentStorageDelegate
	// Session resumption storage: Optional. Support session resumption when provided.
	// Must be initialized before being provided.
	SessionResumptionStorage lib.SessionResumptionStorage
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
	AclStorage server.AclStorage
	// Network native params can be injected depending on the
	// selected Endpoint implementation

	// Network native params can be injected depending on the
	// selected Endpoint implementation
	EndpointNativeParams func()

	// Optional. Support test event triggers when provided. Must be initialized before being
	// provided.
	TestEventTriggerDelegate server.TestEventTriggerDelegate
	// Operational keystore with access to the operational keys: MUST be injected.
	OperationalKeystore storage2.PersistentStorageOperationalKeystore
	// Operational certificate store with access to the operational certs in persisted storage:
	// must not be null at timne of Server::initCommissionableData().
	OpCertStore credentials.PersistentStorageOpCertStore
}

func NewServerInitParams() *InitParams {
	return &InitParams{}
}

func (this *InitParams) Init(options *config.DeviceOptions) (*InitParams, error) {
	this.OperationalServicePort = options.SecuredDevicePort
	this.UserDirectedCommissioningPort = options.UnsecuredCommissionerPort
	this.InterfaceId = options.InterfaceId
	return this, nil
}

func NewCommonCaseDeviceServerInitParams() *CommonCaseDeviceServerInitParams {
	c := &CommonCaseDeviceServerInitParams{
		InitParams: InitParams{
			OperationalKeystore:           nil,
			OperationalServicePort:        config.GetDeviceOptionsInstance().SecuredDevicePort,
			UserDirectedCommissioningPort: config.GetDeviceOptionsInstance().UnsecuredCommissionerPort,
			InterfaceId:                   config.GetDeviceOptionsInstance().InterfaceId,
		},
	}
	return c
}

func (this *InitParams) InitializeStaticResourcesBeforeServerInit() error {

	var sKvsPersistentStorageDelegate storage.PersistentStorageDelegate
	var sPersistentStorageOperationalKeystore = storage2.NewPersistentStorageOperationalKeystoreImpl()
	var sPersistentStorageOpCertStore = credentials.NewPersistentStorageOpCertStoreImpl()
	var sGroupDataProvider = credentials.NewGroupDataProviderImpl()
	var sDefaultCertValidityPolicy = NewIgnoreCertificateValidityPolicy()

	var sSessionResumptionStorage = lib.NewSimpleSessionResumptionStorage()

	if this.PersistentStorageDelegate == nil {
		sKvsPersistentStorageDelegate = storage.KeyValueStoreMgr()
		this.PersistentStorageDelegate = sKvsPersistentStorageDelegate
	}

	if this.OperationalKeystore == nil {
		sPersistentStorageOperationalKeystore.Init(this.PersistentStorageDelegate)
		this.OperationalKeystore = sPersistentStorageOperationalKeystore
	}

	if this.OpCertStore == nil {
		sPersistentStorageOpCertStore.Init(this.PersistentStorageDelegate)
		this.OpCertStore = sPersistentStorageOpCertStore
	}

	sGroupDataProvider.SetStorageDelegate(this.PersistentStorageDelegate)
	err := sGroupDataProvider.Init()
	if err != nil {
		return err
	}
	this.GroupDataProvider = sGroupDataProvider

	{
		if config.ChipConfigEnableSessionResumption {
			err := sSessionResumptionStorage.Init(this.PersistentStorageDelegate)
			if err != nil {
				return err
			}
			this.SessionResumptionStorage = sSessionResumptionStorage
		} else {
			this.SessionResumptionStorage = nil
		}

	}

	this.AccessDelegate = access.GetAccessControlDelegate()

	{
		//TODO 未实现
		this.AclStorage = server.NewAclStorageImpl()
	}

	this.CertificateValidityPolicy = sDefaultCertValidityPolicy

	return nil
}

func (p *CommonCaseDeviceServerInitParams) InitializeStaticResourcesBeforeServerInit() error {

	var sKvsPersistentStorageDelegate storage.PersistentStorageDelegate
	var sPersistentStorageOperationalKeystore storage2.PersistentStorageOperationalKeystore
	var sPersistentStorageOpCertStore credentials.PersistentStorageOpCertStore
	var sGroupDataProvider credentials.GroupDataProvider
	var sDefaultCertValidityPolicy = NewIgnoreCertificateValidityPolicy()

	if p.PersistentStorageDelegate == nil {
		sKvsPersistentStorageDelegate = storage.KeyValueStoreMgr()
		p.PersistentStorageDelegate = sKvsPersistentStorageDelegate
	}
	if p.OperationalKeystore == nil {
		sPersistentStorageOperationalKeystore = storage2.NewPersistentStorageOperationalKeystoreImpl()
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
		p.AclStorage = server.NewAclStorageImpl()
	}

	p.CertificateValidityPolicy = sDefaultCertValidityPolicy

	return nil
}
