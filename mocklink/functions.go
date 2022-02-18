package mocklink

import (
	"fmt"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
)

func getRouteFamily(route *netlink.Route) int {
	if route.Dst != nil && route.Dst.IP != nil {
		return nl.GetIPFamily(route.Dst.IP)
	}

	if route.Src != nil {
		return nl.GetIPFamily(route.Src)
	}

	if route.Gw != nil {
		return nl.GetIPFamily(route.Gw)
	}

	return unix.AF_UNSPEC
}

func matchingRoutes(r1, r2 *netlink.Route) bool {
	if (r1.Dst != nil && r2.Dst == nil) || (r1.Dst == nil && r2.Dst != nil) {
		return false
	}

	if r1.Dst != nil && r2.Dst != nil {
		if !r1.Dst.IP.Equal(r2.Dst.IP) {
			return false
		}
	}

	if (r1.Src != nil && r2.Src == nil) || (r1.Src == nil && r2.Src != nil) {
		return false
	}

	if r1.Src != nil && r2.Src != nil {
		if !r1.Src.Equal(r2.Src) {
			return false
		}
	}

	if (r1.Gw != nil && r2.Gw == nil) || (r1.Gw == nil && r2.Gw != nil) {
		return false
	}

	if r1.Gw != nil && r2.Gw != nil {
		if !r1.Gw.Equal(r2.Gw) {
			return false
		}
	}

	if r1.MPLSDst != r2.MPLSDst {
		return false
	}

	if r1.LinkIndex != r2.LinkIndex {
		return false
	}

	return true
}

func flaggedRouteAttrs(route *netlink.Route, flags uint64) *netlink.Route {
	qr := &netlink.Route{}

	if flags&netlink.RT_FILTER_DST != 0 {
		qr.Dst = route.Dst
	}

	if flags&netlink.RT_FILTER_GW != 0 {
		qr.Gw = route.Gw
	}

	if flags&netlink.RT_FILTER_SRC != 0 {
		qr.Src = route.Src
	}

	if flags&netlink.RT_FILTER_TABLE != 0 {
		qr.Table = route.Table
	}

	if flags&netlink.RT_FILTER_OIF != 0 {
		qr.LinkIndex = route.LinkIndex
	}

	if flags&netlink.RT_FILTER_IIF != 0 {
		qr.ILinkIndex = route.ILinkIndex
	}

	return qr
}

func validateRoute(route *netlink.Route) error {
	if (route.Dst == nil || route.Dst.IP == nil) && route.Src == nil && route.Gw == nil && route.MPLSDst == nil {
		return fmt.Errorf("one of Dst.IP, Src, or Gw must not be nil")
	}

	family := -1

	if route.NewDst != nil {
		if family != -1 && family != route.NewDst.Family() {
			return fmt.Errorf("new destination and destination are not the same address family")
		}
	}

	if route.Src != nil {
		srcFamily := nl.GetIPFamily(route.Src)
		if family != -1 && family != srcFamily {
			return fmt.Errorf("source and destination ip are not the same IP family")
		}
		family = srcFamily
	}

	if route.Gw != nil {
		gwFamily := nl.GetIPFamily(route.Gw)
		if family != -1 && family != gwFamily {
			return fmt.Errorf("gateway, source, and destination ip are not the same IP family")
		}
		family = gwFamily
	}

	if len(route.MultiPath) > 0 {
		for _, nh := range route.MultiPath {
			if nh.Gw != nil {
				gwFamily := nl.GetIPFamily(nh.Gw)
				if family != -1 && family != gwFamily {
					return fmt.Errorf("gateway, source, and destination ip are not the same IP family")
				}
			}
			if nh.NewDst != nil {
				if family != -1 && family != nh.NewDst.Family() {
					return fmt.Errorf("new destination and destination are not the same address family")
				}
			}
		}
	}

	return nil
}
