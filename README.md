# georedis

georedis provide the operations of geolocation with redis.

include
GEOADD, GEODIST, GEOHASH, GEOPOS, GEORADIUS and GEORADIUSBYMEMBER

# Install
```
go get github.com/Natata/georedis
```


# Operation

### Init
```
configFilePath := "config_example.json"
pool, _ := georedis.NewPool(configFilePath)
geo := georedis.NewGeo(pool)
```

### Add
```
key := "japan"
members := []*georedis.Member{
    georedis.NewMember("tokyo", 35.688825, 139.700804),
    georedis.NewMember("osaka", 34.662707, 135.502293),
}   
err := geo.Add(key, members)
```

### Dist
```
dist, err := geo.Dist("japan", "tokyo", "osaka", georedis.KM)
fmt.Println("Distance from toyko to osaka: ", dist)
```

### Hash
```
hs, err := geo.Hash("japan", "tokyo")
```

### Pos
```
pos, err := geo.Pos(key, "tokyo")
```

### RadiusByName
```
neighbors, err := geo.RadiusByName(
	"japan", 
	"tokyo", 
	400, 
	georedis.KM, 
	georedis.WithDist)
```

### Radius
```
neighbors, err := geo.Radius(
	"japan", 
	"tokyo", 
	georedis.Coordinate{
		Lat: 35.688825,
		Lon: 139.700804,
	}, 
	georedis.KM, 
	georedis.WithDist)
```
