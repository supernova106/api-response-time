package config

import (
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ConfigSuite struct{}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) TestLoad(c *C) {
	os.Setenv("MONGODB_URL", "mongodb://127.0.0.1:27017/fluentd")
	os.Setenv("GIN_ENV", "testing")

	defer os.Setenv("MONGODB_URL", "")
	defer os.Setenv("GIN_ENV", "")

	config, err := Load(".env")
	c.Assert(err, IsNil)
	c.Assert(config.MysqlUrl, Equals, "mongodb://127.0.0.1:27017/fluentd")
	c.Assert(config.GinEnv, Equals, "testing")
}
