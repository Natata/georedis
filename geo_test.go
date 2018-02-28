package georadis

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TODO: mock redis

const testConfigFile = "config_test.json"

func TestGeoSetGet(t *testing.T) {
	Convey("test geo basic functions", t, func() {
		pool, err := NewPool(testConfigFile)
		So(err, ShouldBeNil)
		gc := NewGeo(pool)

		k := "yoyo"
		data := []*MetaData{
			NewMetaData("a1", 23.1, 100.7),
			NewMetaData("a2", 23.9, 100.8),
		}
		err = gc.Set(k, data)
		So(err, ShouldBeNil)

		actualData, err := gc.Get(k, []string{"a1", "a2", "b1"})
		So(err, ShouldBeNil)
		So(len(actualData), ShouldEqual, 3)
		So(actualData[2], ShouldBeNil)
	})
}

func TestGeoNeighbors(t *testing.T) {
	Convey("test geo neighbors function", t, func() {
		pool, err := NewPool(testConfigFile)
		So(err, ShouldBeNil)
		gc := NewGeo(pool)

		k := "yoyo"
		data := []*MetaData{
			NewMetaData("a1", 23.1, 100.7),
			NewMetaData("a2", 23.9, 100.8),
		}
		err = gc.Set(k, data)
		So(err, ShouldBeNil)

		loc := Coordinate{
			Lat: 23.09,
			Lon: 100.69,
		}
		results, err := gc.Neighbors(k, loc, 10, KM)
		So(err, ShouldBeNil)
		So(results, ShouldNotBeNil)
		So(len(results), ShouldEqual, 1)

		results, err = gc.Neighbors(k, loc, 100, KM)
		So(err, ShouldBeNil)
		So(results, ShouldNotBeNil)
		So(len(results), ShouldEqual, 2)
	})
}

func TestGeoNeighborsWithParameter(t *testing.T) {
	Convey("test geo neighbors function with parameter", t, func() {
		pool, err := NewPool(testConfigFile)
		So(err, ShouldBeNil)
		gc := NewGeo(pool)

		k := "yoyo"
		data := []*MetaData{
			NewMetaData("a1", 23.1, 100.7),
			NewMetaData("a2", 23.9, 100.8),
		}
		err = gc.Set(k, data)
		So(err, ShouldBeNil)

		loc := Coordinate{
			Lat: 23.09,
			Lon: 100.69,
		}

		results, err := gc.Neighbors(k, loc, 10, KM, WithHash, WithCoord)
		So(err, ShouldBeNil)
		So(results, ShouldNotBeNil)
		So(len(results), ShouldEqual, 1)
		So(results[0].Name, ShouldEqual, "a1")
		So(results[0].Coord.Lat, ShouldAlmostEqual, 23.1, .001)
		So(results[0].Coord.Lon, ShouldAlmostEqual, 100.7, .001)
	})
}

func TestGeoDist(t *testing.T) {
	Convey("test dist", t, func() {
		pool, _ := NewPool(testConfigFile)
		geo := NewGeo(pool)
		k := "yoyo"
		data := []*MetaData{
			NewMetaData("a1", 23.1, 100.7),
			NewMetaData("a2", 23.9, 100.8),
		}
		err := geo.Set(k, data)
		So(err, ShouldBeNil)

		r, err := geo.Dist(k, "a1", "a2", KM)
		So(err, ShouldBeNil)
		So(r, ShouldAlmostEqual, 89.5638, .001)
	})
}

func TestGeoPos(t *testing.T) {
	Convey("test geopos", t, func() {
		pool, _ := NewPool(testConfigFile)
		geo := NewGeo(pool)
		k := "yoyo"
		data := []*MetaData{
			NewMetaData("a1", 23.1, 100.7),
			NewMetaData("a2", 23.9, 100.8),
		}
		err := geo.Set(k, data)
		So(err, ShouldBeNil)

		coords, err := geo.GeoPos(k, "a1", "a2")
		So(err, ShouldBeNil)
		So(len(coords), ShouldEqual, 2)
		So(coords[0].Lat, ShouldAlmostEqual, 23.1, .001)
		So(coords[0].Lon, ShouldAlmostEqual, 100.7, .001)
		So(coords[1].Lat, ShouldAlmostEqual, 23.9, .001)
		So(coords[1].Lon, ShouldAlmostEqual, 100.8, .001)
	})
}
