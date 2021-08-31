package store

import (
	"github.com/matrixorigin/talent-challenge/matrixbase/distributed/pkg/cfg"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
	"time"
)

// Store the store interface
type Store interface {
	// Set set key-value to store
	Set(key []byte, value []byte) error
	// Get returns the value from store
	Get(key []byte) ([]byte, error)
	// Delete remove the key from store
	Delete(key []byte) error
}

// NewStore create the raft store
func NewStore(cfg cfg.StoreCfg) (Store, error) {
	if cfg.Memory {
		return newMemoryStore()
	}
    
	// TODO: need to implement
    cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})

	embed.StartEtcd()



	return nil, nil
}
