package config

const (
	FATCONFDIR    = "/tmp"
	SYSCONFDIR    = "/tmp"
	LOCALSTATEDIR = "/tmp"

	CHIP_DEFAULT_FACTORY_PATH = FATCONFDIR + "/chip_factory.ini"
	CHIP_DEFAULT_CONFIG_PATH  = SYSCONFDIR + "/chip_config.ini"
	CHIP_DEFAULT_DATA_PATH    = LOCALSTATEDIR + "/chip_counters.ini"
)
