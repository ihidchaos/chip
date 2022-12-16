package store

import (
	"fmt"
	gini "gopkg.in/ini.v1"
	"os"
	"time"
)

type iniFile struct {
	mConfigStore *gini.File
}

func newIniFile() *iniFile {
	s := &iniFile{}
	s.mConfigStore = gini.Empty()
	return s
}

func (i *iniFile) init() error {
	return i.removeAll()
}

func (i *iniFile) addEntry(key, value string) error {
	section, err := i.getDefaultSection()
	if err != nil {
		return err
	}
	_, err = section.NewKey(key, value)
	return err
}

func (i *iniFile) removeEntry(key string) error {
	section, err := i.getDefaultSection()
	if err != nil {
		return err
	}
	section.DeleteKey(key)
	return nil
}

func (i *iniFile) removeAll() error {
	i.mConfigStore = gini.Empty()
	return nil
}

func (i *iniFile) addConfig(configFile string) error {
	var err error
	_, err = os.ReadFile(configFile)
	if err != nil {
		_, err = os.Create(configFile)
		if err != nil {
			return err
		}
	}
	i.mConfigStore, err = gini.Load(configFile)
	if err != nil {
		return err
	}
	return nil
}

func (i *iniFile) commitConfig(configFile string) error {
	err := i.mConfigStore.SaveTo(configFile)
	return err
}

func (i *iniFile) readUInt16Value(key string) (uint16, error) {
	u, err := i.mConfigStore.Section("").Key(key).Uint()
	return uint16(u), err
}

func (i *iniFile) readUintValue(key string) (uint, error) {
	section, err := i.getDefaultSection()
	if err != nil {
		return 0, err
	}
	return section.Key(key).Uint()
}

func (i *iniFile) readUint64Value(key string) (uint64, error) {
	section, err := i.getDefaultSection()
	if err != nil {
		return 0, err
	}
	return section.Key(key).Uint64()
}

func (i *iniFile) readFloatValue(key string) (float64, error) {
	section, err := i.getDefaultSection()
	if err != nil {
		return 0, err
	}
	return section.Key(key).Float64()
}

func (i *iniFile) readStringValue(key string) (string, error) {
	section, err := i.getDefaultSection()
	if err != nil {
		return "", err
	}
	return section.Key(key).String(), nil
}

func (i *iniFile) readTimeValue(key string) (time.Time, error) {
	section, err := i.getDefaultSection()
	if err != nil {
		return time.Time{}, err
	}
	return section.Key(key).Time()
}

func (i *iniFile) readBinaryValue(key string) ([]byte, error) {
	section, err := i.getDefaultSection()
	if err != nil {
		return nil, err
	}
	data := section.Key(key).Uints(" ")
	if len(data) == 0 {
		return nil, fmt.Errorf("key vaild")
	}
	bts := make([]byte, 0)
	for _, i := range data {
		bts = append(bts, byte(i))
	}
	return bts, nil
}

func (i *iniFile) hasValue(key string) bool {
	section, err := i.getDefaultSection()
	if err != nil {
		return false
	}
	return section.HasValue(key)
}

func (i *iniFile) getDefaultSection() (*gini.Section, error) {
	var section *gini.Section
	var err error
	if section = i.mConfigStore.Section("DEFAULT"); section != nil {
		return section, nil
	} else {
		section, err = i.mConfigStore.NewSection("DEFAULT")
	}
	return section, err
}
