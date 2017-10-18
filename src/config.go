package src

func loadConfigFromFile(config string) []Host {
	var host Host
	host["name"] = "hahah"
	host["url"] = "sdfasdfasdf"
	host["healthurl"] = "localhsot"
	return []Host{}
}

func loadConfigFromURL(url string) []Host {
	return []Host{}
}
