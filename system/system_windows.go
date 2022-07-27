//go:build windows || (js && wasm)

package system

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

const (
	FatConfDir     = "tmp"
	SysConfDir     = "tmp"
	LocalStatedDir = "tmp"

	FatConfDirFile     = "chip_factory.ini"
	SysConfDirFile     = "chip_config.ini"
	LocalStatedDirFile = "chip_counters.ini"
)
