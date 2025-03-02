package storage

import (
	"fmt"

	"github.com/graph-gophers/dataloader"
)

// the IntKey is analog to the StringKey
type IntKey int

func (k IntKey) String() string   { return fmt.Sprintf("%d", k) }
func (k IntKey) Raw() interface{} { return int(k) }

func NewKeysFromInts(ints []int) dataloader.Keys {
	list := make(dataloader.Keys, len(ints))
	for i := range ints {
		list[i] = IntKey(ints[i])
	}
	return list
}

func IntKeysToSlice(keys dataloader.Keys) []int {
	intIds := make([]int, len(keys))
	for i, id := range keys {
		//nolint:errcheck // by design
		intIds[i] = id.Raw().(int)
	}
	return intIds
}
