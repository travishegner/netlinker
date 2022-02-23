package mocklink

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

//Handle provides a mock netlink handle that satisfies the netlinker interface
//Any method which typically reads or manipulates kernel state, should be overridden
//to maintain an in-memory mock state.
type Handle struct {
	routes []*netlink.Route
	//addresses []*netlink.Addr
	//links     []*netlink.Link
}

//NewHandle returns a new mocklink handle
func NewHandle() *Handle {
	return &Handle{
		routes: make([]*netlink.Route, 0),
	}
}

//RouteAdd adds a new route to mocklink
func (h *Handle) RouteAdd(route *netlink.Route) error {
	if err := validateRoute(route); err != nil {
		return err
	}

	for _, rt := range h.routes {
		if matchingRoutes(rt, route) {
			return fmt.Errorf("route already exists")
		}
	}

	h.routes = append(h.routes, route)
	return nil
}

//RouteDel deletes a route from mocklink
func (h *Handle) RouteDel(route *netlink.Route) error {
	if err := validateRoute(route); err != nil {
		return err
	}
	rtIndex := -1
	for i, rt := range h.routes {
		if matchingRoutes(rt, route) {
			rtIndex = i
			break
		}
	}

	if rtIndex == -1 {
		return fmt.Errorf("route not found")
	}

	h.routes = append(h.routes[:rtIndex], h.routes[rtIndex+1:]...)
	return nil
}

//RouteListFiltered returns a filtered list of routes from mocklink
func (h *Handle) RouteListFiltered(family int, route *netlink.Route, flags uint64) ([]netlink.Route, error) {
	results := make([]netlink.Route, 0)

	if route == nil {
		route = &netlink.Route{}
	}

	qr := flaggedRouteAttrs(route, flags)

	for _, rt := range h.routes {
		if family != netlink.FAMILY_ALL {
			f1 := getRouteFamily(rt)
			if f1 != family {
				continue
			}
		}

		r1 := flaggedRouteAttrs(rt, flags)

		if matchingRoutes(r1, qr) {
			results = append(results, *rt)
		}
	}

	return results, nil
}
