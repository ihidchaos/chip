package device

import (
	"github.com/galenliu/chip"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/crypto"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

const (
	KMaxDiscriminatorValue     uint16 = 0xFFF
	KSpake2pMinPbkdfIterations uint32 = 1000
	KMinSetupPasscode          uint32 = 1
	KMaxSetupPasscode          uint32 = 0x5F5E0FE
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

type CommissionableData struct {
	mIsInitialized          bool
	mSerializedPaseVerifier []byte
	mPaseSalt               []byte
	mPaseIterationCount     uint32
	mSetupPasscode          uint32
	mDiscriminator          uint16
}

var defaultDataProvider atomic.Value

func init() {
	instance := &CommissionableData{}
	defaultDataProvider.Store(instance)
}

func DefaultCommissionableDateProvider() *CommissionableData {
	_instance := defaultDataProvider.Load().(*CommissionableData)
	return _instance
}

func SetDefaultCommissionableDataProvider(provider *CommissionableData) {
	defaultDataProvider.Store(provider)
}

func (c *CommissionableData) Init(options *config.DeviceOptions) error {
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
		options.Payload.Discriminator.SetLongValue(options.Discriminator)
	} else {
		var testOnlyCommissionableDataProvider = TestOnlyCommissionableDataProvider{}
		defaultTestDiscriminator, err := testOnlyCommissionableDataProvider.GetSetupDiscriminator()
		if err != nil {
			log.Panic(err.Error())
		}
		options.Payload.Discriminator.SetLongValue(defaultTestDiscriminator)
	}
	spake2pIterationCount := KSpake2pMinPbkdfIterations
	if options.Spake2pIterations != 0 {
		spake2pIterationCount = options.Spake2pIterations
	}
	log.Printf("PASE PBKDF iterations set to %d\n", spake2pIterationCount)

	err := c.initCommissionableData(options.Spake2pVerifier, options.Spake2pSalt, spake2pIterationCount, setupPasscode, options.Discriminator)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommissionableData) initCommissionableData(serializedSpake2pVerifier, spake2pSalt []byte,
	spake2pIterationCount, setupPasscode uint32,
	discriminator uint16) error {

	if c.mIsInitialized {
		return chip.ErrorIncorrectState
	}
	if discriminator > KMaxDiscriminatorValue {
		log.Infof("Discriminator value invalid: %d", discriminator)
		return chip.ErrorInvalidArgument
	}
	if spake2pIterationCount < KSpake2pMinPbkdfIterations || spake2pIterationCount > crypto.Spake2pMaxPBKDFIterations {
		log.Printf("PASE Iteration count invalid: %d", spake2pIterationCount)
		return chip.ErrorInvalidArgument
	}

	spake2pVerifier := crypto.Spake2pVerifier{}
	havePaseVerifier := serializedSpake2pVerifier != nil && len(serializedSpake2pVerifier) > 0
	var finalSerializedVerifier []byte
	if havePaseVerifier {
		if len(serializedSpake2pVerifier) != crypto.Spake2pVerifierSerializedLength {
			log.Error("PASE verifier size invalid: %d", len(serializedSpake2pVerifier))
			return chip.ErrorInvalidArgument
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
		return chip.ErrorInvalidArgument
	}

	spake2pSaltLength := len(spake2pSalt)
	if havePaseSalt && ((spake2pSaltLength < crypto.Spake2pMinPBKDFSaltLength) || (spake2pSaltLength > crypto.Spake2pMaxPBKDFSaltLength)) {
		log.Infof("PASE salt length invalid: %d", spake2pSaltLength)
		return chip.ErrorInvalidArgument
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
		//if (serializedPasscodeVerifier != serializedSpake2pVerifier.Raw())
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

func (c *CommissionableData) GetSetupDiscriminator() (uint16, error) {
	if !c.mIsInitialized {
		return 0, chip.ErrorIncorrectState
	}
	return c.mDiscriminator, nil
}

func (c *CommissionableData) SetSetupDiscriminator(uint16) error {
	return chip.ErrorNotImplemented
}

func (c *CommissionableData) GetSpake2pIterationCount() (uint32, error) {
	if !c.mIsInitialized {
		return 0, chip.ErrorIncorrectState
	}
	return c.mPaseIterationCount, nil
}

func (c *CommissionableData) GetSpake2pSalt() (bytes []byte, err error) {
	if !c.mIsInitialized {
		return nil, err.IncorrectState
	}
	return c.mPaseSalt, nil
}

func (c *CommissionableData) GetSpake2pVerifier() ([]byte, error) {
	if !c.mIsInitialized {
		return nil, chip.ErrorIncorrectState
	}
	if len(c.mSerializedPaseVerifier) != crypto.Spake2pVerifierSerializedLength {
		return nil, chip.ErrorInternal
	}
	return c.mSerializedPaseVerifier, nil
}

func (c *CommissionableData) GetSetupPasscode() (uint32, error) {
	if !c.mIsInitialized {
		return 0, chip.ErrorIncorrectState
	}
	if c.mSetupPasscode == 0 {
		return 0, chip.ErrorNotImplemented
	}
	return c.mSetupPasscode, nil
}

func (c *CommissionableData) SetSetupPasscode(uint322 uint32) error {
	return chip.ErrorNotImplemented
}

func GeneratePaseSalt() ([]byte, error) {
	return []byte("Pase Salt 2022"), nil
}
