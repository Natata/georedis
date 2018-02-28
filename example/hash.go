package main

import (
	"fmt"
	"georedis"
)

func main() {
	configFilePath := "../config_test.json"
	pool, _ := georedis.NewPool(configFilePath)
	geo := georedis.NewGeo(pool)
	key := "japan"
	members := []*georedis.Member{
		georedis.NewMember("tokyo", 35.688825, 139.700804),
	}
	geo.Add(key, members)

	hs, _ := geo.Hash(key, "tokyo")
	fmt.Println("geohash of tokyo: ", hs[0])
}
