package src

func loadConfigFromFile(config string) []Health {
	var host Health = make(Health)
	host["name"] = "hahah"
	host["key"] = "first"
	host["url"] = "sdfasdfasdf"
	host["healthurl"] = "localhsot"
	return []Health{}
}

func loadConfigFromURL(url string) []Health {
	return []Health{}
}
