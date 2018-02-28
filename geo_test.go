package georedis

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TODO: mock redis

const testConfigFile = "config_test.json"

var (
	k          = "yoyo"
	testMember = []*Member{
		NewMember("a1", 23.1, 100.7),
		NewMember("a2", 23.9, 100.8),
	}
)

func TestGeoSetGet(t *testing.T) {
	Convey("test geo basic functions", t, func() {
		pool, err := NewPool(testConfigFile)
		So(err, ShouldBeNil)
		geo := NewGeo(pool)
		err = geo.Add(k, testMember)
		So(err, ShouldBeNil)

		actualData, err := geo.Pos(k, "a1", "a2", "b1")
		So(err, ShouldBeNil)
		So(len(actualData), ShouldEqual, 3)
		So(actualData[2], ShouldBeNil)
		So(actualData[0].Coord.Lat, ShouldAlmostEqual, 23.1, .0001)
	})
}

func TestGeoNeighbors(t *testing.T) {
	Convey("test geo neighbors function", t, func() {
		pool, err := NewPool(testConfigFile)
		So(err, ShouldBeNil)
		geo := NewGeo(pool)
		err = geo.Add(k, testMember)
		So(err, ShouldBeNil)

		loc := Coordinate{
			Lat: 23.09,
			Lon: 100.69,
		}
		results, err := geo.Radius(k, loc, 10, KM)
		So(err, ShouldBeNil)
		So(results, ShouldNotBeNil)
		So(len(results), ShouldEqual, 1)

		results, err = geo.Radius(k, loc, 100, KM)
		So(err, ShouldBeNil)
		So(results, ShouldNotBeNil)
		So(len(results), ShouldEqual, 2)
	})
}

func TestGeoNeighborsWithParameter(t *testing.T) {
	Convey("test geo neighbors function with parameter", t, func() {
		pool, err := NewPool(testConfigFile)
		So(err, ShouldBeNil)
		geo := NewGeo(pool)
		err = geo.Add(k, testMember)
		So(err, ShouldBeNil)

		loc := Coordinate{
			Lat: 23.09,
			Lon: 100.69,
		}

		results, err := geo.Radius(k, loc, 10, KM, WithHash, WithCoord)
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
		err := geo.Add(k, testMember)
		So(err, ShouldBeNil)

		r, err := geo.Dist(k, "a1", "a2", KM)
		So(err, ShouldBeNil)
		So(r, ShouldAlmostEqual, 89.5638, .001)
	})
}

func TestGeoHash(t *testing.T) {
	Convey("test geohash", t, func() {
		pool, _ := NewPool(testConfigFile)
		geo := NewGeo(pool)
		err := geo.Add(k, testMember)
		So(err, ShouldBeNil)

		hashs, err := geo.Hash(k, "a1", "a2")
		So(err, ShouldBeNil)
		So(len(hashs), ShouldEqual, 2)
		So(hashs[0], ShouldEqual, "whpe7mpx200")
		So(hashs[1], ShouldEqual, "whpxvyb7d50")
	})
}
