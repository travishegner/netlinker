package netlinker

import (
	"flag"
	"log"
	"net"
	"os"
	"reflect"
	"testing"

	"github.com/travishegner/netlinker/mocklink"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func TestMain(m *testing.M) {
	local := flag.Bool("local", false, "set this flag to test whether mocklink behaves the same as netlink (requires root)")

	flag.Parse()

	if !*local {
		log.Println("skipping local only tests")
		os.Exit(0)
	}

	if os.Getuid() != 0 {
		log.Fatalf("local tests must be run as root")
	}

	ns, err := netns.New()
	if err != nil {
		log.Fatalf("failed to create a new network namespace")
	}

	res := m.Run()

	err = ns.Close()
	if err != nil {
		log.Fatalf("failed to close network namespace")
	}

	os.Exit(res)
}

func TestAddFilterRoute(t *testing.T) {
	ml := mocklink.NewHandle()
	nl, err := netlink.NewHandle()
	if err != nil {
		t.Errorf("failed to get new netlink handle: %v", err)
	}

	//err = netlink.LinkAdd(&netlink.Dummy{})
	la := netlink.NewLinkAttrs()
	la.Name = "dummy0"
	dummy0 := &netlink.Dummy{LinkAttrs: la}
	err = netlink.LinkAdd(dummy0)
	if err != nil {
		t.Errorf("failed to add dummy0 link: %v", err)
	}

	addr, err := netlink.ParseAddr("10.10.10.10/24")
	if err != nil {
		t.Errorf("failed to parse address for dummy0: %v", err)
	}

	err = netlink.AddrAdd(dummy0, addr)
	if err != nil {
		t.Errorf("failed to add address to dummy0: %v", err)
	}

	err = netlink.LinkSetUp(dummy0)
	if err != nil {
		t.Errorf("failed to set dummy0 up")
	}

	dst, err := netlink.ParseIPNet("10.10.20.0/24")
	if err != nil {
		t.Errorf("failed to parse destination IPNet: %v", err)
	}

	rt := &netlink.Route{
		Dst: dst,
		Gw:  net.ParseIP("10.10.10.254"),
	}

	err = ml.RouteAdd(rt)
	if err != nil {
		t.Errorf("failed to add route to mocklink: %v", err)
	}

	err = nl.RouteAdd(rt)
	if err != nil {
		t.Errorf("failed to add route to netlink: %v", err)
	}

	ml_routes, err := ml.RouteListFiltered(netlink.FAMILY_ALL, rt, netlink.RT_FILTER_DST|netlink.RT_FILTER_GW)
	if err != nil {
		t.Errorf("failed to get filtered routes from mocklink")
	}

	nl_routes, err := nl.RouteListFiltered(netlink.FAMILY_ALL, rt, netlink.RT_FILTER_DST|netlink.RT_FILTER_GW)
	if err != nil {
		t.Errorf("failed to get filtered routes from netlink")
	}

	if !reflect.DeepEqual(ml_routes, nl_routes) {
		t.Errorf("routes returned should be equal")
	}
}
