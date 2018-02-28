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
		georedis.NewMember("okinawa", 26.212313, 127.679153),
	}
	geo.Add(key, members)

	neighbors, _ := geo.RadiusByName("japan", "tokyo", 400, georedis.KM, georedis.WithDist)
	fmt.Println("Neighbors of from tokyo in 400 km: ")
	for _, n := range neighbors {
		fmt.Println(n.Name, "distance ", n.Dist)
	}
}
