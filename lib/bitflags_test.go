package lib

import (
	"errors"
	log "golang.org/x/exp/slog"
	"testing"
)

func TestBitFlag(t *testing.T) {
	var value uint64 = 0b10101010101010
	var flag1 uint64 = 0b10
	var flag2 uint64 = 0b11
	isTrue := HasFlags(value, flag1, flag2)

	log.Warn("Test flags", "flag1", flag1, "flag2", flag2, "isTrue", isTrue)
	log.Info("Test flags", "flag1", flag1, "flag2", flag2, "isTrue", isTrue)
	log.Error("Test flags", errors.New("new error"))
}
