# Distributed Locker
[![Build Status](https://travis-ci.org/zoer/locker.svg)](https://travis-ci.org/zoer/locker)
[![Go Report
Card](https://goreportcard.com/badge/github.com/zoer/locker)](https://goreportcard.com/report/github.com/zoer/locker)
[![GoDoc](https://godoc.org/github.com/zoer/locker?status.svg)](https://godoc.org/github.com/zoer/locker)


## Install

```
$ go get github.com/zoer/locker
```
## Usage

```go
package main

import (
	"context"
	"log"

	etcd "github.com/coreos/etcd/clientv3"
)

func main() {
	// connect to Etcd instance
	cl, err := clientv3.NewFromURL("127.0.0.1:2379")
	if err != nil {
		log.Fatalf("unable connect to Etcd")
	}

	lkr := NewEtcd(cl)

	// cancel the lock via context timeout after 30 seconds
	ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)

	l, err := lkr.Lock(ctx, WithKey("lock-key"))
	if err != nil {
		log.Fatalf("unable get a lock: %v", err)
	}
	defer l.Unlock()

	// ...
}
```

# TODO
- [x] Etcd
- [ ] Consul
- [ ] Redis
