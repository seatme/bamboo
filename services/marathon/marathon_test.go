package marathon

import (
	"testing"

	. "github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/QubitProducts/bamboo/configuration"
)

func TestHealthCheckDisabledAllowsUnhealthy(t *testing.T) {
	Convey("#healthCheckDisabled", t, func() {
		unhealthyTask := marathonTask{}

		config := new(configuration.Configuration)
		config.Marathon.RequireHealthCheck = false

		Convey("should accept unhealthy task if checks disabled", func() {
			So(taskReady(config, unhealthyTask), ShouldEqual, true)
		})

		results := healthCheckResult{LastSuccess: "asdf"}
		unhealthyTask.HealthCheckResults = append(unhealthyTask.HealthCheckResults, results)

		Convey("should accept healthy task if checks disabled", func() {
			So(taskReady(config, unhealthyTask), ShouldEqual, true)
		})
	})
}

func TestHealthCheckEnabledDisallowsUnhealthy(t *testing.T) {
	Convey("#healthCheckEnabled", t, func() {
		unhealthyTask := marathonTask{}

		config := new(configuration.Configuration)
		config.Marathon.RequireHealthCheck = true

		Convey("should not accept unhealthy task if checks enabled", func() {
			So(taskReady(config, unhealthyTask), ShouldEqual, false)
		})

		results := healthCheckResult{LastSuccess: "asdf"}
		unhealthyTask.HealthCheckResults = append(unhealthyTask.HealthCheckResults, results)

		Convey("should accept healthy task if checks enabled", func() {
			So(taskReady(config, unhealthyTask), ShouldEqual, true)
		})
	})
}

func TestParseHealthCheckPathTCP(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []marathonHealthCheck{
			marathonHealthCheck{"/", "TCP", 0},
			marathonHealthCheck{"/foobar", "TCP", 0},
			marathonHealthCheck{"", "TCP", 0},
		}
		Convey("should return no path if all checks are TCP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "")
		})
	})
}

func TestParseHealthCheckPathHTTP(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []marathonHealthCheck{
			marathonHealthCheck{"/first", "HTTP", 0},
			marathonHealthCheck{"/", "HTTP", 0},
			marathonHealthCheck{"", "HTTP", 0},
		}
		Convey("should return the first path if all checks are HTTP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "/first")
		})
	})
}

func TestParseHealthCheckPathMixed(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []marathonHealthCheck{
			marathonHealthCheck{"", "TCP", 0},
			marathonHealthCheck{"/path", "HTTP", 0},
			marathonHealthCheck{"/", "HTTP", 0},
		}
		Convey("should return the first path if some checks are HTTP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "/path")
		})
	})
}
