package lib

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTestFiles(t *testing.T) {

	Convey("Test", t, func() {
		Convey("NormaliseTestFiles", func() {
			Convey("should only keep valid paths", func() {
				implicit := GetImplicitConfig()
				config := SetupConfig{
					From:           ExpandHome(implicit.Home)("~/FROM"),
					ImplicitConfig: implicit,
				}
				output := NormaliseTestFiles(config, []string{
					"foo", "~/oiwef", "/bar", "~/FROM/.foo/bar", ".bax",
				})

				So(output, ShouldResemble, []string{
					"foo", ".foo/bar", ".bax",
				})
			})
		})
	})
}
