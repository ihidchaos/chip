package DeviceLayer

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	KMaxDiscriminatorValue uint16 = 0xFFF

	KMinSetupPasscode uint32 = 1
	KMaxSetupPasscode uint32 = 0x5F5E0FE
)

type CommissionableDataProvider interface {
	GetSetupDiscriminator() (uint16, error)
	SetSetupDiscriminator(uint16) error
	GetSpake2pIterationCount() (uint32, error)
	GetSpake2pSalt() ([]byte, error)
	GetSpake2pVerifier() ([]byte, error)
	GetSetupPasscode() (uint32, error)
	SetSetupPasscode(uint32) error
}

type CommissionableDataImpl struct {
	mIsInitialized          bool
	mSerializedPaseVerifier []byte
	mPaseSalt               []byte
	mPaseIterationCount     uint32
	mSetupPasscode          uint32
	mDiscriminator          uint16
}

var _instance *CommissionableDataImpl
var _once sync.Once

func GetCommissionableDateProviderInstance() *CommissionableDataImpl {
	_once.Do(func() {
		if _instance == nil {
			_instance = &CommissionableDataImpl{}
		}
	})
	return _instance
}

func NewCommissionableDataImpl() *CommissionableDataImpl {
	return GetCommissionableDateProviderInstance()
}

func (c *CommissionableDataImpl) Init(options *config.DeviceOptions) (*CommissionableDataImpl, error) {
	var setupPasscode uint32
	if options.Payload.SetUpPINCode != 0 {
		setupPasscode = options.Payload.SetUpPINCode
	}
	if options.Spake2pVerifier == nil {
		var testOnlyCommissionableDataProvider = TestOnlyCommissionableDataProvider{}
		defaultTestPasscode, err := testOnlyCommissionableDataProvider.GetSetupPasscode()
		if err != nil {
			log.Panic(err.Error())
		}
		setupPasscode = defaultTestPasscode
		options.Payload.SetUpPINCode = defaultTestPasscode
	}

	if options.Discriminator != 0 {
		options.Payload.Discriminator = options.Discriminator
	} else {
		var testOnlyCommissionableDataProvider = TestOnlyCommissionableDataProvider{}
		defaultTestDiscriminator, err := testOnlyCommissionableDataProvider.GetSetupDiscriminator()
		if err != nil {
			log.Panic(err.Error())
		}
		options.Payload.Discriminator = defaultTestDiscriminator
	}
	spake2pIterationCount := crypto.KSpake2p_Min_PBKDF_Iterations
	if options.Spake2pIterations != 0 {
		spake2pIterationCount = options.Spake2pIterations
	}
	log.Printf("PASE PBKDF iterations set to %d\n", spake2pIterationCount)

	err := c.initCommissionableData(options.Spake2pVerifier, options.Spake2pSalt, spake2pIterationCount, setupPasscode, options.Discriminator)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *CommissionableDataImpl) initCommissionableData(serializedSpake2pVerifier, spake2pSalt []byte,
	spake2pIterationCount, setupPasscode uint32,
	discriminator uint16) error {

	if c.mIsInitialized {
		return lib.ChipErrorIncorrectState
	}
	if discriminator > KMaxDiscriminatorValue {
		log.Infof("Discriminator value invalid: %d", discriminator)
		return lib.ChipErrorInvalidArgument
	}
	if spake2pIterationCount < crypto.KSpake2p_Min_PBKDF_Iterations || spake2pIterationCount > crypto.KSpake2p_Max_PBKDF_Iterations {
		log.Printf("PASE Iteration count invalid: %d", spake2pIterationCount)
		return lib.ChipErrorInvalidArgument
	}

	spake2pVerifier := crypto.Spake2pVerifier{}
	havePaseVerifier := serializedSpake2pVerifier != nil && len(serializedSpake2pVerifier) > 0
	var finalSerializedVerifier []byte
	if havePaseVerifier {
		if len(serializedSpake2pVerifier) != crypto.KSpake2p_VerifierSerialized_Length {
			log.Error("PASE verifier size invalid: %d", len(serializedSpake2pVerifier))
			return lib.ChipErrorInvalidArgument
		}
		err := spake2pVerifier.Deserialize(serializedSpake2pVerifier)
		if err != nil {
			log.Infof("Failed to deserialized PASE verifier: %s", err.Error())
			return err
		}
		log.Print("Got externally provided verifier, using it.")
	}
	havePaseSalt := spake2pSalt != nil && len(spake2pSalt) > 0
	if havePaseVerifier && !havePaseSalt {
		log.Infof("CommissionableDataProvider didn't get a PASE salt, but got a verifier: ambiguous data")
		return lib.ChipErrorInvalidArgument
	}

	spake2pSaltLength := len(spake2pSalt)
	if havePaseSalt && ((spake2pSaltLength < crypto.KSpake2p_Min_PBKDF_Salt_Length) || (spake2pSaltLength > crypto.KSpake2p_Max_PBKDF_Salt_Length)) {
		log.Infof("PASE salt length invalid: %d", spake2pSaltLength)
		return lib.ChipErrorInvalidArgument
	}

	if !havePaseSalt {
		log.Infof("CommissionableDataProvider didn't get a PASE salt, generating one.")
		spake2pSaltBytes, err := GeneratePaseSalt()
		if err != nil {
			log.Infof("Failed to generate PASE salt: %s.", err.Error())
			return err
		}
		spake2pSalt = spake2pSaltBytes
	}

	havePasscode := setupPasscode != 0
	passcodeVerifier := crypto.Spake2pVerifier{}
	var serializedPasscodeVerifier []byte
	if havePasscode {
		err := passcodeVerifier.Generate(spake2pIterationCount, spake2pSalt, setupPasscode)
		if err != nil {
			log.Infof("Failed to generate PASE verifier from passcode: %s", err.Error())
			return err
		}
		//TODO 这里需要确认
		_, err = passcodeVerifier.Serialize()
		if err != nil {
			log.Infof("Failed to serialize PASE verifier from passcode: %s", err.Error())
			return err
		}
	}
	// Make sure we actually have a verifier
	if !havePasscode && !havePaseVerifier {
		log.Infof("Missing both externally provided verifier and passcode: cannot produce final verifier")
		//return lib.CHIP_ERROR_INVALID_ARGUMENT
	}

	if havePasscode && havePaseVerifier {
		//if (serializedPasscodeVerifier != serializedSpake2pVerifier.Value())
		//{
		//	ChipLogError(Support, "Mismatching verifier between passcode and external verifier. Validate inputs.");
		//	return CHIP_ERROR_INVALID_ARGUMENT;
		//}
		//ChipLogProgress(Support, "Validated externally provided passcode matches the one generated from provided passcode.");
	}

	if havePaseVerifier {
		finalSerializedVerifier = serializedSpake2pVerifier
	} else {
		finalSerializedVerifier = serializedPasscodeVerifier
	}
	c.mDiscriminator = discriminator
	c.mSerializedPaseVerifier = finalSerializedVerifier
	c.mPaseSalt = spake2pSalt
	c.mPaseIterationCount = spake2pIterationCount
	if havePasscode {
		c.mSetupPasscode = setupPasscode
	}
	c.mIsInitialized = true
	return nil
}

func (c *CommissionableDataImpl) GetSetupDiscriminator() (uint16, error) {
	if !c.mIsInitialized {
		return 0, lib.ChipErrorIncorrectState
	}
	return c.mDiscriminator, nil
}

func (c *CommissionableDataImpl) SetSetupDiscriminator(uint16) error {
	return lib.ChipErrorNotImplemented
}

func (c *CommissionableDataImpl) GetSpake2pIterationCount() (uint32, error) {
	if !c.mIsInitialized {
		return 0, lib.ChipErrorIncorrectState
	}
	return c.mPaseIterationCount, nil
}

func (c *CommissionableDataImpl) GetSpake2pSalt() (bytes []byte, err error) {
	if !c.mIsInitialized {
		return nil, lib.ChipErrorIncorrectState
	}
	return c.mPaseSalt, nil
}

func (c *CommissionableDataImpl) GetSpake2pVerifier() ([]byte, error) {
	if !c.mIsInitialized {
		return nil, lib.ChipErrorIncorrectState
	}
	if len(c.mSerializedPaseVerifier) != crypto.KSpake2p_VerifierSerialized_Length {
		return nil, lib.ChipErrorInternal
	}
	return c.mSerializedPaseVerifier, nil
}

func (c CommissionableDataImpl) GetSetupPasscode() (uint32, error) {
	if !c.mIsInitialized {
		return 0, lib.ChipErrorIncorrectState
	}
	if c.mSetupPasscode == 0 {
		return 0, lib.ChipErrorNotImplemented
	}
	return c.mSetupPasscode, nil
}

func (c CommissionableDataImpl) SetSetupPasscode(uint322 uint32) error {
	return lib.ChipErrorNotImplemented
}

func GeneratePaseSalt() ([]byte, error) {
	return []byte("Pase Salt 2022"), nil
}
