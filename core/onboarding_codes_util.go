package core

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/device"
	log "github.com/sirupsen/logrus"
)

func GetPayloadContents(payload *config.PayloadContents, aRendezvousFlags uint8) error {
	payload.Version = 0
	payload.RendezvousInformation = aRendezvousFlags

	var err error
	payload.SetUpPINCode, err = device.GetCommissionableDateProvider().GetSetupPasscode()
	if err != nil {
		log.Infof("*** Using default EXAMPLE passcode %d ***", config.UseTestSetupPinCode)
		payload.SetUpPINCode = config.UseTestSetupPinCode
	}
	discriminator, err := device.GetCommissionableDateProvider().GetSetupDiscriminator()
	if err != nil {
		log.Infof("GetCommissionableDataProvider()->GetSetupDiscriminator() failed: %s", err.Error())
		return err
	}
	payload.Discriminator.SetLongValue(discriminator)

	payload.ProductID, err = device.DefaultInstanceInfo().GetProductId()
	if err != nil {
		log.Printf("GetDeviceInstanceInfoProvider()->GetProductId() failed: %s", err.Error())
	}
	return err
}

func PrintOnboardingCodes(contents config.PayloadContents) {

}
