package confaccessor

import (
	"io/ioutil"
)

type SameHostConfaccessor struct {
	//HostName string
	//Path     string
}

func (s SameHostConfaccessor) ReadData(a GetPatherAddresser) (string, error) {
	content, err := ioutil.ReadFile(a.GetPath())
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s SameHostConfaccessor) WriteData(a GetPatherAddresser, data string) error {
	err := ioutil.WriteFile(a.GetPath(), []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}
