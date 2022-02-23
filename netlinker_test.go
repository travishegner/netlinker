package netlinker

import (
	"testing"

	"github.com/vishvananda/netlink"
)

func TestNetlinkFulfillsNetlinker(t *testing.T) {
	var _ Handle = (*netlink.Handle)(nil)
}
