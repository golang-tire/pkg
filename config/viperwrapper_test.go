package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

var sampleConfigW = []byte(`
--- 
apiVersion: v1
server: 
  delay: 3.4
  host: 127.0.0.1
  port: 9090
services: 
  db: 
    connection: "postgres://"
    engine: postgres
`)

func TestMain(t *testing.M) {

	fileName := writeConfig("wrapper_test", sampleConfigW)
	code := t.Run()
	err := os.Remove(fileName)
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func Test_initViper(t *testing.T) {

	var settingChanged bool = false
	onChange := func() error {
		settingChanged = true
		return nil
	}

	err := initViper("wrapper_test", "yaml", "test", onChange)
	if err != nil {
		t.Errorf("init viper failed %v", err)
	}

	serverHost := viper.GetString("server.host")
	if serverHost != "127.0.0.1" {
		t.Errorf("viper settings is not loaded correct")
	}

	var sampleConfigChanged = []byte(`
--- 
apiVersion: v1
server: 
  delay: 3.4
  host1: "0.0.0.0"
  port: 9080
services: 
  db: 
    connection: "postgres://"
    engine: postgres
urls:
  - "url-one"
  - "url-two"
  - "url-three"
	`)

	writeConfig("wrapper_test", sampleConfigChanged)
	if settingChanged != true {
		t.Error("setting file changed but watch not worked")
	}
}
