package mdns

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/mdns"
)

type MdnsService struct {
	*mdns.Server
	localInstanceName  string
	localInterfaceName string
	serviceName        string
	entriesCh          chan *mdns.ServiceEntry
	discoveredSevices  map[string]*mdns.ServiceEntry
}

// mdns机制的服务发现，同时提供注册及发现服务的功能
func NewMdnsService(instanceName, serviceName string, ifaceName string) (*MdnsService, error) {
	IP, err := getIPFromInterface(ifaceName)
	if nil != err {
		return nil, err
	}
	// Setup our service export
	info := []string{fmt.Sprintf("My service: %v.%v", serviceName, instanceName)}
	service, err := mdns.NewMDNSService(instanceName, serviceName, "", "", 9007, []net.IP{IP}, info)
	if nil != err {
		return nil, err
	}
	// Create the mDNS server
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if nil != err {
		return nil, err
	}

	return &MdnsService{
		Server:             server,
		serviceName:        serviceName,
		localInstanceName:  instanceName,
		localInterfaceName: ifaceName,
		discoveredSevices:  make(map[string]*mdns.ServiceEntry),
	}, nil
}

func (m *MdnsService) StartDiscover() error {
	itf, err := net.InterfaceByName(m.localInterfaceName)
	if nil != err {
		return err
	}

	// Make a channel for results and start listening
	m.entriesCh = make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range m.entriesCh {
			m.discoveredSevices[entry.Name] = entry
			log.Printf("discover: %v: %v", entry.Name, entry.AddrV4)
		}
	}()

	// Start the lookup
	params := mdns.DefaultParams(m.serviceName)
	params.DisableIPv6 = true
	params.Entries = m.entriesCh
	params.Interface = itf
	go func() {
		ticker := time.Tick(time.Second * 30)
		for {
			select {
			case <-ticker:
				if err := mdns.Query(params); nil != err {
					close(m.entriesCh)
					log.Printf("mdns query error: %v", err)
				}
			case <-m.entriesCh:
				log.Printf("mdns query end")
				return
			}
		}
	}()
	return nil
}

func (m *MdnsService) GetDiscoveredServiceAddr() (addr []string) {
	for name, entry := range m.discoveredSevices {
		if strings.Contains(name, m.localInstanceName) {
			continue
		}
		addr = append(addr, entry.Addr.String())
	}
	return addr
}

func (m *MdnsService) Shutdown() error {
	close(m.entriesCh)
	return m.Server.Shutdown()
}

func getIPFromInterface(interfaceName string) (ip net.IP, err error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("get interface %v error: %v", interfaceName, err)
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("get interface %v addr error: %v", interfaceName, err)
	}

	if len(addrs) == 0 {
		return nil, fmt.Errorf("get interface %v addr error: no addr", interfaceName)
	}
	return addrs[0].(*net.IPNet).IP, nil
}
