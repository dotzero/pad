package app

import (
	"strings"
)

// CommonOpts is the options that provided into handlers
type CommonOpts struct {
	BoltPath   string
	SecretKey  string
	StaticPath string
	TmlPath    string
	TplExt     string
}

// SetCommon apply the options
func (c *CommonOpts) SetCommon(commonOpts CommonOpts) {
	c.BoltPath = strings.TrimSuffix(commonOpts.BoltPath, "/")
	c.SecretKey = commonOpts.SecretKey
	c.StaticPath = strings.TrimSuffix(commonOpts.StaticPath, "/")
	c.TmlPath = strings.TrimSuffix(commonOpts.TmlPath, "/")
	c.TplExt = commonOpts.TplExt
}
