package version

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFindVersion(t *testing.T) {
	Convey("查找指定名称的版本", t, func() {
		v0 := &Version{
			Name: "1.12.5",
		}
		v1 := &Version{
			Name: "1.11.10",
		}
		v2 := &Version{
			Name: "1.9.7",
		}

		items := []*Version{v0, v1, v2}

		v, err := FindVersion(items, "1.11.10")
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.Name, ShouldEqual, "1.11.10")

		v, err = FindVersion(items, "1.11.11")
		So(err, ShouldEqual, ErrVersionNotFound)
		So(v, ShouldBeNil)
	})
}
