package db

import (
	"github.com/simar7/gokv/types"
)

type mockConfig struct {
}

func (mc mockConfig) Get(input types.GetItemInput) (bool, error) {
	panic("implement me")
}

func (mc mockConfig) Set(input types.SetItemInput) error {
	panic("implement me")
}

func (mc mockConfig) BatchSet(input types.BatchSetItemInput) error {
	panic("implement me")
}
