package confaccessor

import (
	"fmt"
)

type GetAddresser interface {
	GetAddress() string
}

type GetPather interface {
	GetPath() string
}
type GetPatherAddresser interface {
	GetAddresser
	GetPather
}
type ConfAccessor interface {
	ReadData(a GetPatherAddresser) (string, error)
	WriteData(a GetPatherAddresser, content string) error
}

type DummyAccessor struct{}

func (d DummyAccessor) ReadData(a GetPatherAddresser) (string, error) {
	return fmt.Sprintf("dummy:///%s", a.GetPath()), nil
}
func (d DummyAccessor) WriteData(a GetPatherAddresser, content string) error {
	return nil
}

var Registry map[string]ConfAccessor

//type ConfAccessorGenerator func() Confaccessr

func init() {
	Registry = map[string]ConfAccessor{}
	Registry["Locahost"] = SameHostConfaccessor{}
	Registry["Dummy"] = DummyAccessor{}
}
