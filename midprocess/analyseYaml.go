package midprocess

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Conf struct {
	PodName      string `yaml:"podName"` //yaml：yaml格式 enabled：属性的为enabled
	ImageTagName string `yaml:"imageTagName"`
}

func (c *Conf) GetConf() *Conf {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}
	return c
}
