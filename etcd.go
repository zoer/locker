package locker

import (
	"context"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"
)

const etcdNamespace = "/disributed-locker/"
const etcdDefaultKey = "_default"
const etcdDefaultTTL = 20 * time.Second // 20 seconds
const etcdDefaultWaitTTL = 5 * time.Second

type etcdLocker struct {
	client *etcd.Client
	opts   options
}

type etcdLock struct {
	cancel  context.CancelFunc
	leaseID etcd.LeaseID
}

// NewEtcd creates Etcd locker
func NewEtcd(cl *etcd.Client, oo ...Option) Locker {
	ncl := *cl
	ncl.KV = namespace.NewKV(cl.KV, etcdNamespace)
	ncl.Lease = namespace.NewLease(cl.Lease, etcdNamespace)
	ncl.Watcher = namespace.NewWatcher(cl.Watcher, etcdNamespace)
	opts := newEtcdOptions()
	opts.apply(oo...)

	return &etcdLocker{client: &ncl, opts: *opts}
}

func newEtcdOptions() *options {
	o := newOptions()
	o.apply(
		WithTTL(etcdDefaultTTL),
		WithWaitTTL(etcdDefaultWaitTTL),
		WithKey(etcdDefaultKey))

	return o
}

func newEtcdLock(ctx context.Context, cl *etcd.Client, opts options) (Lock, error) {
	ctx, cancel := context.WithCancel(ctx)
	wctx, wcancel := context.WithTimeout(ctx, opts.waitTTL)
	defer wcancel()

	l := &etcdLock{
		cancel: cancel,
	}

	gr, err := cl.Grant(wctx, int64(opts.ttl.Seconds()))
	if err != nil {
		return nil, err
	}
	l.leaseID = gr.ID

	keepc, err := cl.KeepAlive(ctx, l.leaseID)
	if err != nil {
		return nil, err
	}

	go func() {
		defer cl.Revoke(context.TODO(), l.leaseID)
		for {
			if _, live := <-keepc; !live {
				break
			}
		}
	}()

	for {
		rt, err := cl.Txn(wctx).
			If(etcd.Compare(etcd.CreateRevision(opts.key), "=", 0)).
			Then(etcd.OpPut(opts.key, "", etcd.WithLease(l.leaseID))).
			Commit()
		if err != nil {
			cancel()
			return nil, err
		}
		if rt.Succeeded {
			break
		}

		select {
		case <-wctx.Done():
		case <-time.After(200 * time.Millisecond):
		}
	}

	return l, nil
}

// Lock gets the lock, ctx is used only for locking operation(not the lock itself)
func (l etcdLocker) Lock(ctx context.Context, oo ...Option) (Lock, error) {
	return newEtcdLock(ctx, l.client, l.opts.withOptions(oo...))
}

// Unlock unlocks the lock and cancels locking context.
func (l *etcdLock) Unlock() {
	if l.cancel != nil {
		l.cancel()
	}
}
