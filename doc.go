// Usage:
//  package main
//
//  import (
//  	"context"
//  	"log"
//
//  	etcd "github.com/coreos/etcd/clientv3"
//  )
//
//  func main() {
//  	cl, err := etcd.NewFromURL("127.0.0.1:2379")
//  	if err != nil {
//  		log.Fatalf("unable connect to Etcd")
//  	}
//
//  	lkr := locker.NewEtcd(cl)
//
//  	l, err := lkr.Lock(context.TODO(), WithKey("lock-key"), WithWaitTTL(2 * time.Second))
//  	if err != nil {
//  		log.Fatalf("unable get a lock: %v", err)
//  	}
//  	defer l.Unlock()
//  }

package locker
