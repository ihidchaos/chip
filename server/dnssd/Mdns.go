package dnssd

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"net/netip"
	"strconv"
	"sync"
)

type ResponderHandler interface {
	OnQuery(data *QueryData) (*dns.Msg, error)
}
type MdnsServer interface {
	StartServer(port uint16) error
	ShutdownServer()
	SendTo(message *dns.Msg, port netip.AddrPort, id net.Interface) error
	BroadcastSend(message dns.Msg, port uint16, id net.Interface, addr netip.Addr) error
	Init() error
	RemoveServices() error
}

type MdnsServerImpl struct {
	mServer dns.Server
	Handler ResponderHandler
}

var _mdnsServerImplInstance *MdnsServerImpl
var _instanceOnce sync.Once

func MdnsInstance() *MdnsServerImpl {
	_instanceOnce.Do(func() {
		if _mdnsServerImplInstance == nil {
			_mdnsServerImplInstance = &MdnsServerImpl{}
		}
	})
	return _mdnsServerImplInstance
}

func NewMdnsServerImpl() *MdnsServerImpl {
	return &MdnsServerImpl{}
}

func (mdns *MdnsServerImpl) Init() error {
	return nil
}

func (mdns *MdnsServerImpl) SendTo(message *dns.Msg, addr netip.AddrPort, id net.Interface) error {
	client := dns.Client{
		Net: "udp",
	}
	_, _, err := client.Exchange(message, addr.Addr().String())
	return err
}

func (mdns *MdnsServerImpl) BroadcastSend(message dns.Msg, port uint16, id net.Interface, srcAddr netip.Addr) error {
	clint := new(dns.Client)
	if srcAddr.Is4() {
		clint.Net = "udp"
		_, _, err := clint.Exchange(&message, netip.AddrPortFrom(IPv4LinkLocalMulticast, port).String())
		if err != nil {
			return err
		}
	}
	if srcAddr.Is6() {
		clint.Net = "udp6"
		_, _, err := clint.Exchange(&message, netip.AddrPortFrom(IPv6LinkLocalMulticast, port).String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (mdns *MdnsServerImpl) RemoveServices() error {
	return nil
}

func (mdns *MdnsServerImpl) StartServer(port uint16) error {
	mdns.mServer = dns.Server{
		Addr:     ":" + strconv.Itoa(int(port)),
		Net:      "udp",
		Listener: nil,
		Handler:  mdns,
	}
	go func() {
		err := mdns.mServer.ListenAndServe()
		if err != nil {
			log.Print(err.Error())
			return
		}
	}()
	return nil
}

func (mdns *MdnsServerImpl) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {

	src, _ := netip.ParseAddrPort(w.RemoteAddr().String())
	dest, _ := netip.ParseAddrPort(w.LocalAddr().String())
	queryData := NewQuery(r).SetDestAddr(dest).SetSrcAddr(src)

	msg, err := mdns.Handler.OnQuery(queryData)
	if err != nil {
		log.Printf(err.Error())
	}
	err = w.WriteMsg(msg)
	if err != nil {
		log.Printf(err.Error())
	}
}

func (mdns *MdnsServerImpl) ShutdownServer() {
	_ = mdns.mServer.Shutdown()
}
