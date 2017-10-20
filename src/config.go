package src

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"time"
)

func loadConfigFromFile(config string) (map[string]Health, error) {
	data, err := ioutil.ReadFile(config)

	if err != nil {
		return make(map[string]Health), err
	}
	return loadConfig(data)
}

func loadConfig(s []byte) (map[string]Health, error) {
	var hosts []Health
	_hosts := make(map[string]Health)
	err := yaml.Unmarshal(s, &hosts)

	if err != nil {
		return _hosts, err
	}

	for _, v := range hosts {
		key := guessKey(v)
		if key == "" {
			return _hosts, errors.New("can not find key: name/ip/key/hostname must set for a host")
		}
		_, ok := _hosts[key]
		if ok {
			return _hosts, errors.New(fmt.Sprintf("duplicated key: %s", key))
		}
		_, ok = v["healthURL"]
		if !ok {
			return _hosts, errors.New(fmt.Sprintf("no healthURL set for: %s", key))
		}

		_hosts[key] = v
	}
	return _hosts, nil

}

func loadConfigFromURL(url string) (map[string]Health, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return make(map[string]Health), err
	} else {
		if resp.StatusCode != http.StatusOK {
			return make(map[string]Health), err
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return loadConfig(body)
	}
}
