package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var Conf Config

type Config struct {
	DB      DB      `yaml:"db"`
	Monitor Monitor `yaml:"monitor"`
}

type DB struct {
	Path string `yaml:"path"`
}

type Monitor struct {
	PushGateway  string `yaml:"push_gateway"`
	PushInterval int    `yaml:"push_interval"`
}

func InitConfig(path string) error {
	// 读取YAML文件内容
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// 定义一个Config类型的变量
	var config Config

	// 解析YAML文件
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err
	}

	// 打印配置信息
	Conf = config
	return nil
}
