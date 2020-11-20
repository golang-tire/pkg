package config

import (
	"io/ioutil"
	"os"
	"path"
)

func writeConfig(testFilename string, data []byte) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fileName := path.Join(dir, testFilename+".yaml")

	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		panic(err)
	}
	return fileName
}

type stringHolderMock struct {
	v string
}

type stringArrayHolderMock struct {
	v []string
}

type intHolderMock struct {
	v int
}

func (i intHolderMock) Int() int {
	return i.v
}

func (i intHolderMock) Int64() int64 {
	return int64(i.v)
}

func (h stringHolderMock) String() string {
	return h.v
}

func (h stringArrayHolderMock) Strings() []string {
	return h.v
}

// RegisterStringMock mock register string
func RegisterStringMock(key, defValue string) String {
	return stringHolderMock{v: defValue}
}

// RegisterStringArrayMock mock register string array
func RegisterStringArrayMock(key, defValue string) String {
	return stringHolderMock{v: defValue}
}

// RegisterIntMock mock register int
func RegisterIntMock(key string, defValue int) Int {
	return intHolderMock{v: defValue}
}
