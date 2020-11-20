package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var configSample = []byte(`
--- 
apiVersion: v1
server: 
  delay: 3.4
  host: 127.0.0.1
  port: 9090
db: 
  enable: true
  connection: "postgres://"
  engine: postgres
urls:
  - "url-one"
  - "url-two"
  - "url-three"
`)

func TestConfig(t *testing.T) {
	fileName := writeConfig("config_test", configSample)
	err := Init("config_test", "yaml", "config")
	assert.Nil(t, err)

	apiVersion := RegisterString("apiVersion", "")
	serverDelay := RegisterFloat64("server.delay", 0.0)
	serverHost := RegisterString("server.host", "")
	serverPort := RegisterInt("server.port", 0)
	dbConnection := RegisterString("db.connection", "")
	dbEnable := RegisterBool("db.enable", false)
	dbEngine := RegisterString("db.engine", "")
	urls := RegisterStringSlice("urls", nil)

	err = Load()
	assert.Nil(t, err)

	assert.Equal(t, "v1", apiVersion.String())
	assert.Equal(t, 3.4, serverDelay.Float64())
	assert.Equal(t, "127.0.0.1", serverHost.String())
	assert.Equal(t, 9090, serverPort.Int())
	assert.Equal(t, "postgres://", dbConnection.String())
	assert.Equal(t, true, dbEnable.Bool())
	assert.Equal(t, "postgres", dbEngine.String())
	assert.Equal(t, []string{"url-one", "url-two", "url-three"}, urls.Slice())

	// default value
	notExistString := RegisterString("foo.bar", "foo.bar")
	notExistInt := RegisterInt64("foo.bar.int", 7)
	notExistFloat := RegisterFloat64("foo.bar.float", 7.7)

	err = Load()
	assert.Nil(t, err)

	assert.Equal(t, "foo.bar", notExistString.String())
	assert.Equal(t, 7, notExistInt.Int())
	assert.Equal(t, 7.7, notExistFloat.Float64())

	err = os.Remove(fileName)
	assert.Nil(t, err)
}
