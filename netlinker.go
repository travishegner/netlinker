package netlinker

import (
	"github.com/vishvananda/netlink"
)

//Handle is an interface which should be fulfilled by both netlink and mocklink
//this enables easy unit testing for netlink dependent code without requiring
//root or doing actual kernel state modifications
type Handle interface {
	RouteAdd(*netlink.Route) error
	RouteDel(*netlink.Route) error
	RouteListFiltered(int, *netlink.Route, uint64) ([]netlink.Route, error)
}
