package cfg

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type RpcHub struct {
	clients     []*ethclient.Client
	rpcs        []string
	current     int
	mu          sync.RWMutex
	isSwitching bool
}
type ClientService struct {
	rpcHub   *RpcHub
	param    *ParamHub
	stopChan chan struct{}
}

func NewRpcHubs(urls []string) (*RpcHub, error) {
	hub := RpcHub{
		rpcs:    urls,
		current: 0,
	}
	for _, val := range hub.rpcs {
		log.Printf("rpc -> %s", val)
		client, err := ethclient.Dial(val)
		if err != nil {
			log.Printf("newrpchub -> %s : %s", val, err)
			continue
		}
		hub.clients = append(hub.clients, client)
	}
	if len(hub.clients) == 0 {
		return nil, fmt.Errorf("no rpc hub found")
	}
	log.Printf("total clients  -> %v", len(hub.clients))
	return &hub, nil
}
func NewClientService(param *ParamHub, rpc []string) (*ClientService, error) {
	hubs, err := NewRpcHubs(rpc)
	if err != nil {
		return nil, fmt.Errorf("newclientservice ->NewRpcHubs -> %s", err)
	}
	return &ClientService{
		hubs,
		param,
		make(chan struct{}),
	}, nil
}
func (s *ClientService) SwitchRpc() {
	s.rpcHub.mu.Lock()
	defer s.rpcHub.mu.Unlock()
	if s.rpcHub.isSwitching {
		log.Printf("rps is swtiching -> %s", s.rpcHub.rpcs[0])
		return
	}
	s.rpcHub.isSwitching = true
	oldIndex := s.rpcHub.current
	s.rpcHub.current++
	log.Printf("[from -> %s -> to %s] \n ", s.rpcHub.rpcs[oldIndex], s.rpcHub.rpcs[s.rpcHub.current])
	go func() {
		time.Sleep(1 * time.Millisecond)
		s.rpcHub.mu.Lock()
		s.rpcHub.isSwitching = false
		defer s.rpcHub.mu.Unlock()
		log.Printf("using rpc %s", s.rpcHub.rpcs[s.rpcHub.current])
	}()
}
func (s *ClientService) GetCurrentClient() *ethclient.Client {
	s.rpcHub.mu.RLock()
	defer s.rpcHub.mu.RUnlock()
	return s.rpcHub.clients[s.rpcHub.current]
}
func (s *ClientService) GetCurrentRpc() string {
	s.rpcHub.mu.RLock()
	defer s.rpcHub.mu.RUnlock()
	return s.rpcHub.rpcs[s.rpcHub.current]
}
