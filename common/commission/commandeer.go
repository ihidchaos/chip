package commission

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type commandeerCommissionState struct {
	fsCreate sync.Once
	created  chan struct{}
}

type commandeer struct {
	logger logrus.Logger
}
