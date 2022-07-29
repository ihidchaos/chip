//go:build unix || (js && wasm)

package platform

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

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
	ConfigKvsPath = "/tmp/chip_kvs"

	DefaultFactoryPath = filepath.Join(os.TempDir(), "chip_factory.ini")
	DefaultConfigPath  = filepath.Join(os.TempDir(), "chip_config.ini")
	DefaultDataPath    = filepath.Join(os.TempDir(), "chip_counters.ini")
)

const (
	FatConfDir     = "tmp"
	SysConfDir     = "tmp"
	LocalStatedDir = "tmp"

	FatConfDirFile     = "chip_factory.ini"
	SysConfDirFile     = "chip_config.ini"
	LocalStatedDirFile = "chip_counters.ini"
)
