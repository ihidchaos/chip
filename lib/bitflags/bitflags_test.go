package bitflags

import (
	"fmt"
	"github.com/moznion/go-optional"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlags_Has(t *testing.T) {
	var value uint64 = 0b01010101010101
	var ff = Some(value)
	assert.Equal(t, value, ff.Unwrap())
	assert.True(t, ff.Has(0b01))
	assert.True(t, ff.Has(0b01, 0b0101))
	assert.True(t, ff.Has(0b01, 0b0101, 0b01010101))
	assert.True(t, ff.Has(0b01, 0b0101, 0b01010101, 0b01010101010101))

	ff.Set(0b101010101)
	assert.Equal(t, value, ff.Unwrap())

	ff.Set(0b10101010101010)
	fmt.Printf("flag[ %0b ]\t\n", ff.Unwrap())
	assert.Equal(t, uint64(0b11111111111111), ff.Unwrap())

}

func TestOption(t *testing.T) {
	var value optional.Option[uint32]
	value = nil
	if !value.IsSome() {
		t.Log("11111111111")
	}
}
