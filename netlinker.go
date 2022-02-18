package netlinker

import (
	"github.com/vishvananda/netlink"
)

type Handle interface {
	RouteAdd(*netlink.Route) error
	RouteDel(*netlink.Route) error
	RouteListFiltered(int, *netlink.Route, uint64) (*[]netlink.Route, error)
}
