package sheriff

import (
	"net"
	"strconv"
	"sync"
	"time"
)

type ConnectionInfo struct {
	name     string
	addr     *net.IPAddr
	lastPing time.Time
}

type Registry struct {
	members map[string]ConnectionInfo
	mtx     sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		members: make(map[string]ConnectionInfo),
	}
}

func (r *Registry) Set(name string, ip net.IP, port uint16) error {
	addr, err := net.ResolveIPAddr("ipv4", net.JoinHostPort(ip.String(), strconv.Itoa(int(port))))
	if err != nil {
		return err
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.members[name] = ConnectionInfo{
		name:     name,
		addr:     addr,
		lastPing: time.Now().UTC(),
	}
	return nil
}

func (r *Registry) Delete(name string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	delete(r.members, name)
}

func (r *Registry) Len() int {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return len(r.members)
}

func (r *Registry) Lookup(dur time.Duration) []ConnectionInfo {
	threshold := time.Now().UTC().Add(-dur)
	result := []ConnectionInfo{}

	r.mtx.RLock()
	defer r.mtx.RUnlock()
	for _, conn := range r.members {
		if threshold.Before(conn.lastPing) {
			result = append(result, conn)
		}
	}

	return result
}
