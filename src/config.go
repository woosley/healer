package src

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func loadConfigFromFile(config string) map[string]Health {
	var hosts []Health
	data, err := ioutil.ReadFile(config)
	check(err)
	err = yaml.Unmarshal(data, &hosts)
	check(err)
	_hosts := make(map[string]Health)
	for _, v := range hosts {
		key := guessKey(v)
		if key == "" {
			panic("can not find key: name/ip/key/hostname must set for a host")
		}
		_, ok := _hosts[key]
		if ok {
			panic(fmt.Sprintf("duplicated key: %s", key))
		}

		_hosts[key] = v
	}
	return _hosts
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func loadConfigFromURL(url string) map[string]Health {
	return make(map[string]Health)
}
