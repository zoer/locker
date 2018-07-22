package locker

import (
	"context"
	"testing"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/require"
)

const etcdEndpoint = "127.0.0.1:2379"

func TestEtcdLocker_LockUnlock(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	key := "foo1"

	cl := newEtcdClient(t)
	defer cl.Close()

	ctx := context.TODO()
	lkr := NewEtcd(cl)

	l, err := lkr.Lock(ctx, WithKey(key), WithTTL(20*time.Second))
	r.NoError(err)

	testEtcdLock(t, key)

	r.True(isEtcdKeyExist(cl, key), "etcd key should exist: %v", key)

	l.Unlock()
	time.Sleep(100 * time.Millisecond)

	r.False(isEtcdKeyExist(cl, key), "etcd key should not exist: %v", key)

	l, err = lkr.Lock(ctx, WithKey(key), WithWaitTTL(100*time.Millisecond))
	if err != nil {
		t.Fatalf("error while locking: %v", err)
	}
	l.Unlock()
}

func TestEtcdLocker_ContextCancellation(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	key := "foo2"

	cl := newEtcdClient(t)
	defer cl.Close()

	lkr := NewEtcd(cl)

	r.False(isEtcdKeyExist(cl, key), "etcd key should not exist: %v", key)

	ctx, cancel := context.WithCancel(context.TODO())
	_, err := lkr.Lock(ctx, WithKey(key))
	r.NoError(err)

	r.True(isEtcdKeyExist(cl, key), "etcd key should exist: %v", key)

	cancel()
	time.Sleep(200 * time.Millisecond)

	r.False(isEtcdKeyExist(cl, key), "etcd key should not exist: %v", key)
}

func testEtcdLock(t *testing.T, key string) {
	cl := newEtcdClient(t)
	defer cl.Close()

	lkr := NewEtcd(cl)

	_, err := lkr.Lock(context.TODO(), WithKey(key), WithWaitTTL(200*time.Millisecond))

	if err == nil {
		t.Fatalf("the %q key is not locked", key)
	} else if err != context.DeadlineExceeded {
		t.Fatalf("errored while locking: %v", err)
	}
}

func etcdLockKey(key string) string {
	return etcdNamespace + key
}

func isEtcdKeyExist(cl *etcd.Client, key string) bool {
	res, err := cl.Get(context.TODO(), etcdLockKey(key))
	if err != nil {
		panic(err)
	}

	return res.Count > 0
}

func newEtcdClient(t *testing.T) *etcd.Client {
	cl, err := etcd.NewFromURL(etcdEndpoint)
	if err != nil {
		t.Fatalf("unable connect to Etcd instance: %v", err)
	}

	return cl
}
