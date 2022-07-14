package DeviceLayer

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/gateway/pkg/log"
)

const (
	kMaxDiscriminatorValue = 0xFFF
)

var _CommissionDataProvider CommissionableDataProvider

func InitCommissionableDataProvider(provider CommissionableDataProvider) CommissionableDataProvider {
	_CommissionDataProvider = provider
	return _CommissionDataProvider
}

func GetCommissionableDataProvider() CommissionableDataProvider {
	return _CommissionDataProvider
}

type CommissionableData struct {
	mDiscriminator          uint16
	mSerializedPaseVerifier []byte
	mPaseSalt               []byte
	mPaseIterationCount     uint32
	mSetupPasscode          uint32
	mIsInitialized          bool
}

func NewCommissionableData(options *config.DeviceOptions) (*CommissionableData, error) {

	var setupPasscode uint32 = 0
	if options.Payload.SetUpPINCode != 0 {
		setupPasscode = options.Payload.SetUpPINCode
	} else if options.Spake2pVerifier == nil {
		var defaultTestPasscode uint32 = 0
		testOnlyCommissionableDataProvider := TestOnlyCommissionableDataProvider{}
		defaultTestPasscode, err := testOnlyCommissionableDataProvider.GetSetupPasscode()
		if err != nil {
			log.Infof("*** WARNING: Using temporary passcode %u due to no neither --passcode or --spake2p-verifier-base64 "+
				"given on command line. This is temporary and will disappear. Please update your scripts "+
				"to explicitly configure onboarding credentials. ***", defaultTestPasscode)
		}
		setupPasscode = defaultTestPasscode
		options.Payload.SetUpPINCode = defaultTestPasscode
	}
	if options.Discriminator != 0 {
		options.Payload.Discriminator = options.Discriminator
	} else {
		var defaultTestDiscriminator uint16 = 0
		testOnlyCommissionableDataProvider := TestOnlyCommissionableDataProvider{}
		defaultTestDiscriminator, err := testOnlyCommissionableDataProvider.GetSetupDiscriminator()
		if err != nil {
			log.Infof("*** WARNING: Using temporary test discriminator %u due to --discriminator not "+
				"given on command line. This is temporary and will disappear. Please update your scripts "+
				"to explicitly configure discriminator. ***", defaultTestDiscriminator)
		}
		options.Payload.Discriminator = defaultTestDiscriminator
	}

	var spake2pIterationCount uint32 = crypto.KSpake2p_Min_PBKDF_Iterations
	if options.Spake2pIterations != 0 {
		spake2pIterationCount = options.Spake2pIterations
	}
	return newProvider(options.Spake2pVerifier, options.Spake2pSalt, spake2pIterationCount, setupPasscode, options.Payload.Discriminator)
}

func newProvider(serializedSpake2pVerifier []byte, spake2pSalt []byte, spake2pIterationCount uint32, setupPasscode uint32, discriminator uint16) (*CommissionableData, error) {

	if discriminator > kMaxDiscriminatorValue {
		return nil, fmt.Errorf("discriminator value invalid: %d", discriminator)
	}
	if spake2pIterationCount < crypto.KSpake2p_Min_PBKDF_Iterations || spake2pIterationCount > crypto.KSpake2p_Max_PBKDF_Iterations {
		return nil, fmt.Errorf("PASE Iteration count invalid: %d", spake2pIterationCount)
	}

	havePaseVerifier := len(serializedSpake2pVerifier) > 1

	providedVerifier := crypto.Spake2pVerifier{}

	finalSerializedVerifier := make([]byte, 0)

	if havePaseVerifier {
		if len(serializedSpake2pVerifier) != crypto.KSpake2p_VerifierSerialized_Length {
			return nil, fmt.Errorf("PASE verifier size invalid: %d", len(serializedSpake2pVerifier))
		}
		err := providedVerifier.Deserialize(serializedSpake2pVerifier)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialized PASE verifier: %s", err.Error())
		}
		log.Info("got externally provided verifier, using it.")
	}

	havePaseSalt := len(spake2pSalt) > 1
	if havePaseVerifier && !havePaseSalt {
		return nil, fmt.Errorf("commissionableDataProvider didn't get a PASE salt, but got a verifier: ambiguous data")
	}
	var spake2pSaltLength = len(spake2pSalt)
	if havePaseSalt && ((spake2pSaltLength < crypto.KSpake2p_Min_PBKDF_Salt_Length) || (spake2pSaltLength > crypto.KSpake2p_Max_PBKDF_Salt_Length)) {
		return nil, fmt.Errorf("PASE salt length invalid: %d", spake2pSaltLength)
	}

	if !havePaseSalt {
		spake2pSaltVector, err := GeneratePaseSalt()
		if err != nil {
			return nil, fmt.Errorf("failed to generate PASE salt: %s", err.Error())
		}
		spake2pSalt = spake2pSaltVector
	}

	havePasscode := setupPasscode != 0
	var passcodeVerifier = crypto.Spake2pVerifier{}
	serializedPasscodeVerifier := make([]byte, 0)
	saltSpan := spake2pSalt

	if havePasscode {
		err := passcodeVerifier.Generate(spake2pIterationCount, saltSpan, setupPasscode)
		if err != nil {
			return nil, fmt.Errorf("Failed to generate PASE verifier from passcode: %s", err.Error())
		}

		serializedPasscodeVerifier, err = passcodeVerifier.Serialize()
		if err != nil {
			return nil, fmt.Errorf("failed to serialize PASE verifier from passcode: %s", err.Error())
		}
	}
	if !havePasscode && !havePaseVerifier {
		return nil, fmt.Errorf("missing both externally provided verifier and passcode: cannot produce final verifier")
	}

	if havePasscode && havePaseVerifier {
		if string(serializedPasscodeVerifier) != string(serializedSpake2pVerifier) {
			return nil, fmt.Errorf("mismatching verifier between passcode and external verifier. Validate inputs")
		}
	}

	if havePaseVerifier {
		finalSerializedVerifier = serializedPasscodeVerifier
	}
	data := &CommissionableData{}
	data.mDiscriminator = discriminator
	data.mSerializedPaseVerifier = finalSerializedVerifier
	data.mPaseSalt = spake2pSalt
	data.mPaseIterationCount = spake2pIterationCount

	if havePasscode {
		data.mSetupPasscode = setupPasscode
	}
	data.mIsInitialized = true
	return data, nil
}

func GeneratePaseSalt() ([]byte, error) {
	return nil, nil
}

func (c CommissionableData) GetProductId() (uint16, error) {
	//TODO implement me
	panic("implement me")
}

func (c CommissionableData) GetSetupDiscriminator() (uint16, error) {
	//TODO implement me
	panic("implement me")
}

func (c CommissionableData) SetSetupDiscriminator(uint162 uint16) {
	//TODO implement me
	panic("implement me")
}

func (c CommissionableData) GetSpake2pIterationCount() (uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (c CommissionableData) GetSpake2pSalt() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (c CommissionableData) GetSpake2pVerifier() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (c CommissionableData) GetSetupPasscode() (uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (c CommissionableData) SetSetupPasscode(uint322 uint32) {
	//TODO implement me
	panic("implement me")
}
