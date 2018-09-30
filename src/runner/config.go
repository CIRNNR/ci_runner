package runner

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfigStruct struct {
	Image string
	Steps []struct {
		Name string
		Command string
	}
}

func loadConfig(directory string) ConfigStruct  {
	config := new(ConfigStruct)
	data,err := ioutil.ReadFile("./" + directory + "/ci.yaml")
	if err != nil{
		panic(err)
	}
	err = yaml.Unmarshal(data, config)
	if err != nil{
		panic(err)
	}
	return *config
}