package sheriff

import (
	"net"
	"time"
)

type Notifier interface {
	Notify(name string, ip net.IP, port uint16, lastPing time.Time) error
}

type MockNotify struct {
	Name string
	IP   net.IP
	Port uint16
}

func NewMockNotify(name string, ip net.IP, port uint16) Notifier {
	return MockNotify{name, ip, port}
}

func (n MockNotify) Notify(name string, ip net.IP, port uint16, lastPing time.Time) error { return nil }
