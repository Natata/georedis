package main

import "georedis"

func main() {
	configFilePath := "../config_test.json"
	pool, _ := georedis.NewPool(configFilePath)
	geo := georedis.NewGeo(pool)
	key := "japan"
	members := []*georedis.Member{
		georedis.NewMember("tokyo", 35.688825, 139.700804),
	}
	geo.Add(key, members)
}
