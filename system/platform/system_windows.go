//go:build windows || (js && wasm)

package platform

import (
	"os"
	"path/filepath"
)

func GetHomeDir() string {
	p, _ := os.UserHomeDir()
	return p
}

func GetFatConFile() string {
	// return strings.ReplaceAll(strings.ReplaceAll(path.Join(os.TempDir(), FatConfDirFile), "/", string(os.PathSeparator)), "\\", string(os.PathSeparator))
	return filepath.Join(os.TempDir(), FatConfDirFile)
}

func GetSysConFile() string {

	//return strings.ReplaceAll(strings.ReplaceAll(path.Join(os.TempDir(), SysConfDirFile), "/", string(os.PathSeparator)), "\\", string(os.PathSeparator))
	return filepath.Join(os.TempDir(), SysConfDirFile)
}

func GetLocalConFile() string {
	//return strings.ReplaceAll(strings.ReplaceAll(path.Join(os.TempDir(), LocalStatedDirFile), "/", string(os.PathSeparator)), "\\", string(os.PathSeparator))
	return filepath.Join(os.TempDir(), LocalStatedDirFile)
}

func GetConfigsKvsPath() string {
	return filepath.Join("tmp", "chip_kvs")
}

var (
	ConfigKvsPath      = filepath.Join(os.TempDir(), "chip.ini")
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

func SystemLayer() system.Layer {
	return nil
}
