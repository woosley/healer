package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

//looksLikeUrl return true when a string can be parse by url.ParseRequestURI
//and starts with 'http'
func looksLikeUrl(s string) bool {
	if _, err := url.ParseRequestURI(s); err != nil {
		return false
	} else {
		if strings.HasPrefix(strings.ToLower(s), "http") {
			return true
		} else {
			return false
		}
	}
}

//looksLikeHostPort return ture when it is a host:port style string
func looksLikeHostport(s string) bool {
	if !strings.Contains(s, ":") {
		return false
	} else {
		hostport := strings.SplitN(s, ":", 2)
		if _, err := regexp.MatchString("[0-9]+", hostport[1]); err == nil {
			return true
		} else {
			return false
		}
	}
}

//return key name
func guessKey(h Health) string {
	keys := []string{"key", "name", "hostname", "ip"}
	for _, key := range keys {
		if v, okay := h[key]; okay {
			return v
		}
	}
	return ""
}

func doKeyChan(t KeyChan, data map[string]Health) {
	// first read key from keyChan
	key := t.Key
	if v, ok := data[key]; ok {
		//send data to KeyChan
		t.KChan <- v
		// sync Data
		<-t.Sync
	} else {
		// send a nil map
		t.KChan <- make(Health)
	}
}

//{state|status: running, message|content: xxxxx, error: true|false, code: "500"}
func checkHealth(d map[string]interface{}, code int) (state string, reason string, _code string) {
	state = getState(d, code)
	reason = getReason(d, code)
	_code = getCode(d, code)
	return
}

func getState(dd map[string]interface{}, code int) (s string) {
	if v, ok := dd["error"]; ok {
		if v.(bool) {
			return "error"
		} else {
			return "running"
		}
	}
	if v, ok := dd["state"]; ok {
		s = v.(string)
	} else if v, ok := dd["status"]; ok {
		s = v.(string)
	} else {
		if code == 200 {
			s = "running"
		} else if code == 500 {
			s = "error"
		}
	}
	return s
}

func getReason(d map[string]interface{}, code int) string {
	if v, ok := d["reason"]; ok {
		return v.(string)
	} else if v, ok := d["message"]; ok {
		return v.(string)
	}
	return fmt.Sprintf("server retured code: %v", code)
}

func getCode(d map[string]interface{}, code int) string {
	if v, ok := d["code"]; ok {
		switch t := v.(type) {
		case int:
			return fmt.Sprintf("%v", t)
		case string:
			return t
		}
	}
	return fmt.Sprintf("%v", code)
}

func getHealthFromHostPort(hostport string) (string, string, string) {
	conn, _ := net.DialTimeout("tcp", hostport, 10*time.Second)
	if conn != nil {
		defer conn.Close()
		return "running", fmt.Sprintf("%s is open", hostport), "200"
	}
	return "error", fmt.Sprintf("can not open %s", hostport), "500"
}

func checkHealthAPI(url string) (string, string, string) {
	if looksLikeUrl(url) {
		return getHealthFromURL(url)
	} else if looksLikeHostport(url) {
		return getHealthFromHostPort(url)
	}
	return "", "", ""
}

func getHealthFromURL(url string) (string, string, string) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "error", fmt.Sprintf("%s", err), "500"
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		m := make(map[string]interface{})
		err := json.Unmarshal(body, &m)
		//
		if err != nil {
			return "error", fmt.Sprintf("%s", err), "500"
		} else {
			return checkHealth(m, resp.StatusCode)
		}
	}
}
