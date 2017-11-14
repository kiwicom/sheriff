package cluster

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"

	kitlog "github.com/go-kit/kit/log"
)

type pingDelegate struct{ r *Registry }

func (pd pingDelegate) AckPayload() []byte { return []byte{} }

func (pd pingDelegate) NotifyPingComplete(other *memberlist.Node, rtt time.Duration, payload []byte) {
	fmt.Println("ping:", other.Name, other.Addr, rtt)
	_ = pd.r.Set(other.Name, other.Addr, other.Port)
}

type eventDelegate struct{ r *Registry }

func (ed eventDelegate) NotifyJoin(n *memberlist.Node) {
	ed.r.Set(n.Name, n.Addr, n.Port)
}
func (ed eventDelegate) NotifyLeave(n *memberlist.Node) {
	ed.r.Delete(n.Name)
}

// TODO: implement
func (ed eventDelegate) NotifyUpdate(n *memberlist.Node) {}

type Member struct {
	members  *memberlist.Memberlist
	registry *Registry
	notifier Notifier

	ProbeInterval time.Duration
	logger        log.Logger
}

func NewMember(existing []string) (Member, error) {
	logger = kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stderr))

	registry := NewRegistry()
	cfg := memberlist.DefaultLANConfig()
	{
		cfg.BindPort = 0
		cfg.Name = uuid.New().String()
		cfg.Ping = pingDelegate{registry}
		cfg.Events = eventDelegate{registry}
		cfg.ProbeInterval = time.Second
		cfg.Logger = log
		cfg.LogOutput = ioutil.Discard
	}

	m, err := memberlist.Create(cfg)
	if err != nil {
		return Member{}, err
	}

	if len(existing) > 0 {
		_, err = m.Join(existing)
		if err != nil {
			return Member{}, err
		}
	}

	node := m.LocalNode()
	/*
		log.WithFields(logrus.Fields{
			"addr": node.Addr,
			"port": node.Port,
		}).Info("local node")
	*/
	return Member{
		members:       m,
		registry:      registry,
		notifier:      NewMockNotify(node.Name, node.Addr, node.Port),
		ProbeInterval: time.Second,
	}, nil

}

func (n Member) Run() {
	ticker := time.NewTicker(n.ProbeInterval)
	for _ = range ticker.C {
		for _, member := range n.registry.Lookup(n.ProbeInterval) {
			// nodeName is name of node, which should be tested
			_ = member
		}
	}
}
