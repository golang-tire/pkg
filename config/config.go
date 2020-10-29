package config

import (
	"sync"
)

var confWatch = configHolder{}

type configHolder struct {
	lock      sync.RWMutex
	confItems []confItem
}

type confItem struct {
	key      string
	ref      interface{}
	defValue interface{}
}

func (cc *configHolder) addRef(key string, ref interface{}, defValue interface{}) {
	cc.lock.Lock()
	defer cc.lock.Unlock()
	cc.confItems = append(cc.confItems, confItem{key: key, ref: ref, defValue: defValue})
}

// RegisterString register an string variable
func RegisterString(key, defValue string) String { return confWatch.RegisterString(key, defValue) }
func (cc *configHolder) RegisterString(key, defValue string) String {
	var v = defValue
	cc.addRef(key, &v, defValue)
	return stringHolder{value: &v}
}

// RegisterInt register an integer variable
func RegisterInt(key string, defValue int) Int { return confWatch.RegisterInt(key, defValue) }
func (cc *configHolder) RegisterInt(key string, defValue int) Int {
	var v = int64(defValue)
	cc.addRef(key, &v, defValue)
	return intHolder{value: &v}
}

// RegisterInt64 register an int64 variable
func RegisterInt64(key string, defValue int64) Int { return confWatch.RegisterInt64(key, defValue) }
func (cc *configHolder) RegisterInt64(key string, defValue int64) Int {
	var v = defValue
	cc.addRef(key, &v, defValue)
	return intHolder{value: &v}
}

// RegisterFloat32 register a float32 variable
func RegisterFloat32(key string, defValue float32) Float {
	return confWatch.RegisterFloat32(key, defValue)
}
func (cc *configHolder) RegisterFloat32(key string, defValue float32) Float {
	var v = float64(defValue)
	cc.addRef(key, &v, defValue)
	return floatHolder{value: &v}
}

// RegisterFloat64 register a float64 variable
func RegisterFloat64(key string, defValue float64) Float {
	return confWatch.RegisterFloat64(key, defValue)
}
func (cc *configHolder) RegisterFloat64(key string, defValue float64) Float {
	var v = defValue
	cc.addRef(key, &v, defValue)
	return floatHolder{value: &v}
}

// RegisterBool register a bool variable
func RegisterBool(key string, defValue bool) Bool { return confWatch.RegisterBool(key, defValue) }
func (cc *configHolder) RegisterBool(key string, defValue bool) Bool {
	var v = defValue
	cc.addRef(key, &v, defValue)
	return boolHolder{value: &v}
}

func (cc *configHolder) handleChange() error {
	cc.lock.RLock()
	defer cc.lock.RUnlock()

	for _, configItem := range cc.confItems {
		switch configItem.defValue.(type) {
		case string:
			v, err := getViperString(configItem.key, configItem.defValue)
			if err != nil {
				return err
			}
			t := configItem.ref.(*string)
			*t = v
		case int:
			v, err := getViperInt64(configItem.key, configItem.defValue)
			if err != nil {
				return err
			}
			t := configItem.ref.(*int64)
			*t = v
		case int64:
			v, err := getViperInt64(configItem.key, configItem.defValue)
			if err != nil {
				return err
			}
			t := configItem.ref.(*int64)
			*t = v
		case float32:
			v, err := getViperFloat32(configItem.key, configItem.defValue)
			if err != nil {
				return err
			}
			t := configItem.ref.(*float32)
			*t = v
		case float64:
			v, err := getViperFloat64(configItem.key, configItem.defValue)
			if err != nil {
				return err
			}
			t := configItem.ref.(*float64)
			*t = v
		case bool:
			v, err := getViperBool(configItem.key, configItem.defValue)
			if err != nil {
				return err
			}
			t := configItem.ref.(*bool)
			*t = v
		}
	}
	return nil
}

// Load configs and set variables
// can use to reload configs
func Load() error {
	return confWatch.handleChange()
}

// Init initialize config module with and accept confName that is config filename
// ext is config file extension
// appName is software name and will use to make search paths for config file
// config module will search for config files in flowing locations:
// /etc/<appName>
// $HOME/.<appName>
// and beside the executable file
func Init(confName, ext, appName string) error {
	return initViper(confName, ext, appName, confWatch.handleChange)
}
