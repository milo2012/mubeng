package common

import (
	"os"
	"time"

	"ktbs.dev/mubeng/internal/proxymanager"
)

// Options consists of the configuration required.
type Options struct {
	File         string
	Address      string
	Check        bool
	Timeout      time.Duration
	Rotate       int
	Verbose      bool
	Switch       bool
	Output       string
	Result       *os.File
	ProxyManager *proxymanager.ProxyManager
	Daemon       bool
}
