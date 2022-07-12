package storage

import (
	"fmt"
	gini "gopkg.in/ini.v1"
	"os"
)

type IniStorage interface {
	AddConfig(configFile string) error
	CommitConfig(configFile string) error
	GetUInt16Value(key string) (uint16, error)
	GetUIntValue(key string) (uint, error)
	GetUInt64Value(key string) (uint64, error)
	GetStringValue(key string) (string, error)
	GetBinaryBlobValue(key string) ([]byte, error)
	HasValue(key string) bool

	AddEntry(key, value string) error
	RemoveEntry(key string) error
	RemoveAll() error
}

type Ini struct {
	mConfigStore *gini.File
}

func (i Ini) AddEntry(key, value string) error {
	section, err := i.getDefaultSection()
	if err != nil {
		return err
	}
	_, err = section.NewKey(key, value)
	return err
}

func (i Ini) RemoveEntry(key string) error {
	section, err := i.getDefaultSection()
	if err != nil {
		return err
	}
	section.DeleteKey(key)
	return nil
}

func (i Ini) RemoveAll() error {
	i.mConfigStore = gini.Empty()
	return nil
}

func (i *Ini) AddConfig(configFile string) error {
	var err error
	if i.mConfigStore == nil {
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
	} else {
		err := i.mConfigStore.Append(configFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i Ini) CommitConfig(configFile string) error {
	return i.mConfigStore.SaveTo(configFile)
}

func (i Ini) GetUInt16Value(key string) (uint16, error) {
	u, err := i.mConfigStore.Section("").Key(key).Uint()
	return uint16(u), err
}

func (i Ini) GetUIntValue(key string) (uint, error) {
	section, err := i.getDefaultSection()
	if err != nil {
		return 0, err
	}
	return section.Key(key).Uint()
}

func (i Ini) GetUInt64Value(key string) (uint64, error) {
	section, err := i.getDefaultSection()
	if err != nil {
		return 0, err
	}
	return section.Key(key).Uint64()
}

func (i Ini) GetStringValue(key string) (string, error) {
	section, err := i.getDefaultSection()
	if err != nil {
		return "", err
	}
	return section.Key(key).String(), nil
}

func (i Ini) GetBinaryBlobValue(key string) ([]byte, error) {
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

func (i Ini) HasValue(key string) bool {
	section, err := i.getDefaultSection()
	if err != nil {
		return false
	}
	return section.HasValue(key)
}

func (i Ini) getDefaultSection() (*gini.Section, error) {
	return i.mConfigStore.GetSection("DEFAULT")
}
