package mocklink

import (
	"testing"

	"github.com/vishvananda/netlink"
)

func TestRouteAddDelListFiltered(t *testing.T) {
	ml := NewHandle()

	dst1, _ := netlink.ParseIPNet("10.10.10.0/24")
	rt1 := &netlink.Route{
		LinkIndex: 5,
		Dst:       dst1,
	}

	err := ml.RouteAdd(rt1)
	if err != nil {
		t.Errorf("error while adding first route: %v", err)
	}

	if len(ml.routes) != 1 {
		t.Errorf("should be 1 routes in state slice, found %v", len(ml.routes))
	}

	dst2, _ := netlink.ParseIPNet("10.10.20.0/24")
	rt2 := &netlink.Route{
		LinkIndex: 5,
		Dst:       dst2,
	}

	err = ml.RouteAdd(rt2)
	if err != nil {
		t.Errorf("error wile adding second route: %v", err)
	}

	if len(ml.routes) != 2 {
		t.Errorf("should be 2 routes in state slice, found %v", len(ml.routes))
	}

	routes, err := ml.RouteListFiltered(netlink.FAMILY_V6, nil, 0)
	if err != nil {
		t.Errorf("error while getting ipv6 filtered list: %v", err)
	}

	if len(routes) != 0 {
		t.Errorf("should have received 0 routes from our filter but got: %v", len(routes))
	}

	routes, err = ml.RouteListFiltered(netlink.FAMILY_V4, nil, 0)
	if err != nil {
		t.Errorf("error while getting ipv4 filtered list: %v", err)
	}

	if len(routes) != 2 {
		t.Errorf("should have received 2 routes from our filter but got: %v", len(routes))
	}

	routes, err = ml.RouteListFiltered(netlink.FAMILY_ALL, nil, 0)
	if err != nil {
		t.Errorf("error while getting all filtered list: %v", err)
	}

	if len(routes) != 2 {
		t.Errorf("should have received 2 routes from our filter but got: %v", len(routes))
	}

	routes, err = ml.RouteListFiltered(netlink.FAMILY_ALL, &netlink.Route{LinkIndex: 2}, netlink.RT_FILTER_OIF)
	if err != nil {
		t.Errorf("error while getting wrong link index filtered list: %v", err)
	}

	if len(routes) != 0 {
		t.Errorf("should have received 0 routes from our filter but got: %v", len(routes))
	}

	routes, err = ml.RouteListFiltered(netlink.FAMILY_ALL, &netlink.Route{LinkIndex: 5}, netlink.RT_FILTER_OIF)
	if err != nil {
		t.Errorf("error while getting right link index filtered list: %v", err)
	}

	if len(routes) != 2 {
		t.Errorf("should have received 2 routes from our filter but got: %v", len(routes))
	}

	routes, err = ml.RouteListFiltered(netlink.FAMILY_ALL, &netlink.Route{Dst: dst1}, netlink.RT_FILTER_DST)
	if err != nil {
		t.Errorf("error while getting dst filtered list: %v", err)
	}

	if len(routes) != 1 {
		t.Errorf("should have recieved 1 route from our filter but got: %v", len(routes))
	}

	err = ml.RouteDel(rt1)
	if err != nil {
		t.Errorf("error while deleting route by pointer: %v", err)
	}

	if len(ml.routes) != 1 {
		t.Errorf("should be 1 routes in state slice, found %v", len(ml.routes))
	}

	err = ml.RouteDel(&netlink.Route{
		LinkIndex: 5,
		Dst:       dst2,
	})
	if err != nil {
		t.Errorf("error while deleting route by copy: %v", err)
	}

	if len(ml.routes) != 0 {
		t.Errorf("should be 0 routes in state slice, found %v", len(ml.routes))
	}
}

func TestInvalidRoute(t *testing.T) {
	ml := NewHandle()

	err := ml.RouteAdd(&netlink.Route{})
	if err == nil {
		t.Errorf("should be an error on empty route")
	}
}

func TestDuplicateRoute(t *testing.T) {
	ml := NewHandle()

	dst1, _ := netlink.ParseIPNet("10.10.10.0/24")
	rt1 := &netlink.Route{
		LinkIndex: 5,
		Dst:       dst1,
	}

	err := ml.RouteAdd(rt1)
	if err != nil {
		t.Errorf("error while adding first route: %v", err)
	}

	err = ml.RouteAdd(rt1)
	if err == nil {
		t.Errorf("should have failed to add duplicate route")
	}
}
