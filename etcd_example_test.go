package locker

import (
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func Example() {
	// connect to Etcd instance
	cl, err := clientv3.NewFromURL("127.0.0.1:2379")
	if err != nil {
		log.Fatalf("unable connect to Etcd")
	}

	lkr := NewEtcd(cl)

	// cancel the lock via context timeout after 30 seconds
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	l, err := lkr.Lock(ctx, WithKey("lock-key"))
	if err != nil {
		log.Fatalf("unable get a lock: %v", err)
	}
	defer l.Unlock()

	// ...
}
