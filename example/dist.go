package main

import (
	"fmt"

	"github.com/Natata/georedis"
)

func main() {
	configFilePath := "config_example.json"
	pool, _ := georedis.NewPool(configFilePath)
	geo := georedis.NewGeo(pool)
	key := "japan"
	members := []*georedis.Member{
		georedis.NewMember("tokyo", 35.688825, 139.700804),
		georedis.NewMember("osaka", 34.662707, 135.502293),
	}
	geo.Add(key, members)

	dist, _ := geo.Dist("japan", "tokyo", "osaka", georedis.KM)
	fmt.Println("Distance from toyko to osaka: ", dist)
}
