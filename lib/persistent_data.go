package lib

import "github.com/galenliu/chip/lib/store"

type PersistentDataBase interface {
	Save(base store.PersistentStorageDelegate)
}

type PersistentDataImpl struct{}
