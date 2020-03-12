// package config 支持字符串、整型、以及数组  布尔型
package snailframe

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

/*
type configNormalType map[string]interface{}

type config struct {
	data interface{}
}*/

//初始化Conf
func NewConf(configStrcut interface{},configName string) (redata toml.MetaData) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}



	var path string = dir+"/"+ configName
	if data, err := toml.DecodeFile(path, configStrcut); err != nil {
		panic(err)
	}else{
		redata = data
	}

	return redata
}
/*
//加载字符串数组
func (this config)GetSliceString(key string) (conList []string) {

	if value, ok := this.data[key]; ok {

		if val, rightType := value.([]interface{});rightType {
			for _,v := range val {
				if vstr, isString := v.(string);isString {
					conList = append(conList, vstr)
				}
			}
		}else {
			panic("the config [" + key + "] type is not SliceString")
		}
	} else {
		panic("the config [" + key + "] not exist")
	}

	return
}

//加载整数数组
func (this config)GetSliceInt(key string) (conList []int) {

	if value, ok := this.data[key]; ok {

		if val, rightType := value.([]interface{});rightType {
			for _,v := range val {
				if vstr, isString := v.(int64);isString {
					conList = append(conList, int(vstr))
				}
			}
		}else {
			panic("the config [" + key + "] type is not SliceInt")
		}
	} else {
		panic("the config [" + key + "] not exist")
	}

	return
}

//以字符串形式加载配置项
func (this config)GetString(key string) string {
	if value, ok := this.data[key]; ok {
		if val, rightType := value.(string); rightType {
			return val
		} else {
			panic("the config [" + key + "] type is not string")
		}
	} else {
		panic("the config [" + key + "] not exist")
	}
}

//以整数形式加载配置项
func (this config)GetInt(key string) int {
	if value, ok := this.data[key]; ok {
		if val, rightType := value.(int64); rightType {
			return int(val)
		} else {
			panic("the config [" + key + "] type is not int")
		}
	} else {
		panic("the config [" + key + "] not exist")
	}
}

//以布尔形式加载配置项
func (this config)GetBool(key string) bool {
	if value, ok := this.data[key]; ok {
		if val, rightType := value.(bool); rightType {
			return val
		} else {
			panic("the config [" + key + "] type is not bool")
		}
	} else {
		panic("the config [" + key + "] not exist")
	}
}*/
