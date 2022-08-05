package core

import (
	"github.com/galenliu/chip/access"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials"
	storage2 "github.com/galenliu/chip/crypto/operational_storage"
	"github.com/galenliu/chip/lib"
	storage3 "github.com/galenliu/chip/pkg/storage"
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
	PersistentStorageDelegate storage3.PersistentStorageDelegate
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
	OperationalKeystore storage2.OperationalKeystore
	// Operational certificate store with access to the operational certs in persisted storage:
	// must not be null at timne of Server::initCommissionableData().
	OpCertStore credentials.PersistentStorageOpCertStore
}

func NewServerInitParams() *InitParams {
	return &InitParams{}
}

func (params *InitParams) Init(options *config.DeviceOptions) (*InitParams, error) {
	params.OperationalServicePort = options.SecuredDevicePort
	params.UserDirectedCommissioningPort = options.UnsecuredCommissionerPort
	params.InterfaceId = options.InterfaceId
	return params, nil
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

func (params *InitParams) InitializeStaticResourcesBeforeServerInit() error {

	var sKvsPersistentStorageDelegate storage3.PersistentStorageDelegate
	var sPersistentStorageOperationalKeystore = storage2.NewPersistentStorageOperationalKeystoreImpl()
	var sPersistentStorageOpCertStore = credentials.NewPersistentStorageOpCertStoreImpl()
	var sGroupDataProvider = credentials.NewGroupDataProviderImpl()
	var sDefaultCertValidityPolicy = NewIgnoreCertificateValidityPolicy()

	var sSessionResumptionStorage = lib.NewSimpleSessionResumptionStorage()

	if params.PersistentStorageDelegate == nil {
		sKvsPersistentStorageDelegate = storage3.KeyValueStoreMgr()
		params.PersistentStorageDelegate = sKvsPersistentStorageDelegate
	}

	if params.OperationalKeystore == nil {
		sPersistentStorageOperationalKeystore.Init(params.PersistentStorageDelegate)
		params.OperationalKeystore = sPersistentStorageOperationalKeystore
	}

	if params.OpCertStore == nil {
		sPersistentStorageOpCertStore.Init(params.PersistentStorageDelegate)
		params.OpCertStore = sPersistentStorageOpCertStore
	}

	sGroupDataProvider.SetStorageDelegate(params.PersistentStorageDelegate)
	err := sGroupDataProvider.Init()
	if err != nil {
		return err
	}
	params.GroupDataProvider = sGroupDataProvider

	{
		if config.ChipConfigEnableSessionResumption {
			err := sSessionResumptionStorage.Init(params.PersistentStorageDelegate)
			if err != nil {
				return err
			}
			params.SessionResumptionStorage = sSessionResumptionStorage
		} else {
			params.SessionResumptionStorage = nil
		}

	}

	params.AccessDelegate = access.GetAccessControlDelegate()

	{
		//TODO 未实现
		params.AclStorage = NewAclStorageImpl()
	}

	params.CertificateValidityPolicy = sDefaultCertValidityPolicy

	return nil
}

func (p *CommonCaseDeviceServerInitParams) InitializeStaticResourcesBeforeServerInit() error {

	var sKvsPersistentStorageDelegate storage3.PersistentStorageDelegate
	var sPersistentStorageOperationalKeystore storage2.OperationalKeystore
	var sPersistentStorageOpCertStore credentials.PersistentStorageOpCertStore
	var sGroupDataProvider credentials.GroupDataProvider
	var sDefaultCertValidityPolicy = NewIgnoreCertificateValidityPolicy()

	if p.PersistentStorageDelegate == nil {
		sKvsPersistentStorageDelegate = storage3.KeyValueStoreMgr()
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
		p.AclStorage = NewAclStorageImpl()
	}

	p.CertificateValidityPolicy = sDefaultCertValidityPolicy

	return nil
}
