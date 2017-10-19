package src

func guessKey(h Host) string {
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
