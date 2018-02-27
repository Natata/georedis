package georadis

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenPool(t *testing.T) {
	Convey("test generate pool", t, func() {
		pool, err := NewPool("config_test.json")
		So(err, ShouldBeNil)
		So(pool.IdleTimeout, ShouldEqual, 300*time.Second)
	})
}
