//go:build unix || (js && wasm)

package platform

func GetFatConDir() string {
	return FatConfDir
}
func GetFatConFile() string {
	return strings.ReplaceAll(strings.ReplaceAll(path.Join(os.TempDir(), FatConfDirFile), "/", string(os.PathSeparator)), "\\", string(os.PathSeparator))
}

func GetSysConFile() string {

	return strings.ReplaceAll(strings.ReplaceAll(path.Join(os.TempDir(), SysConfDirFile), "/", string(os.PathSeparator)), "\\", string(os.PathSeparator))
}

func GetLocalConFile() string {
	return strings.ReplaceAll(strings.ReplaceAll(path.Join(os.TempDir(), LocalStatedDirFile), "/", string(os.PathSeparator)), "\\", string(os.PathSeparator))
}

var (
	ChipConfigKvsPath = "/tmp/chip_kvs"
)

const (
	FatConfDir     = "tmp"
	SysConfDir     = "tmp"
	LocalStatedDir = "tmp"

	FatConfDirFile     = "chip_factory.ini"
	SysConfDirFile     = "chip_config.ini"
	LocalStatedDirFile = "chip_counters.ini"
)
