package bitflags

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlags_Has(t *testing.T) {
	var value uint64 = 0b01010101010101
	var ff = Some(value)
	assert.Equal(t, value, ff.Value())
	assert.True(t, ff.HasAll(0b01))
	assert.True(t, ff.HasAll(0b01, 0b0101))
	assert.True(t, ff.HasAll(0b01, 0b0101, 0b01010101))
	assert.True(t, ff.HasAll(0b01, 0b0101, 0b01010101, 0b01010101010101))

	ff.Sets(0b101010101)
	assert.Equal(t, value, ff.Value())

	ff.Sets(0b10101010101010)
	fmt.Printf("flag[ %0b ]\t\n", ff.Value())
	assert.Equal(t, uint64(0b11111111111111), ff.Value())

}

type OnList struct {
	list [12]*One
}

func (l *OnList) Delete(o *One) {
	for i, oo := range l.list {
		if oo == o {
			l.list[i] = nil
		}
	}
}

func (l *OnList) Add(o *One) {
	for i, oo := range l.list {
		if oo == nil {
			o.list = l
			l.list[i] = o
			fmt.Printf("add---------%v \t\n", o)
			return
		}
	}
}

func (l *OnList) Add2(o *One) *One {
	for i, oo := range l.list {
		if oo == nil {
			o.list = l
			l.list[i] = o
			return l.list[i]
		}
	}
	return nil
}

type One struct {
	value string
	list  *OnList
}

func (o *One) Release() {
	fmt.Printf("value:%s Release \t\n", o.value)
	o.list.Delete(o)
}

func TestOption(t *testing.T) {
	l := &OnList{}
	o1 := &One{value: "1"}
	o2 := &One{value: "2"}
	o3 := &One{value: "3"}
	o4 := &One{value: "4"}
	o5 := &One{value: "5"}
	o6 := &One{value: "6"}
	o7 := &One{value: "7"}
	o8 := &One{value: "8"}
	o9 := &One{value: "9"}
	l.Add(o1)
	l.Add(o2)
	l.Add(o3)
	l.Add(o4)
	l.Add(o5)
	l.Add(o6)
	l.Add(o7)
	l.Add(o8)

	v := l.Add2(o9)
	fmt.Printf("v.value:%v \t\n", v.value)
	for _, o := range l.list {
		if o != nil {
			fmt.Printf("value:%s \t\n", o.value)
		}
	}

	*v = One{value: "11111", list: l}

	fmt.Printf("-------------------------\t\n")
	o1.Release()
	o2.Release()
	o3.Release()
	o4.Release()
	fmt.Printf("-------------------------\t\n")

	for _, o := range l.list {
		if o != nil && o.value != "" {
			fmt.Printf("value:%s \t\n", o.value)
		}
	}
}
