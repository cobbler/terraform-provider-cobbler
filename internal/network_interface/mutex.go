package network_interface

import "sync"

// systemNetworkInterfaceLock holds a *sync.Mutex per system UID. Cobbler's network
// interface mutation API is not safe for concurrent writes that target the same
// parent system; serialize per system so distinct systems still mutate in parallel.
var systemNetworkInterfaceLock sync.Map

func lockForSystem(systemUid string) *sync.Mutex {
	if m, ok := systemNetworkInterfaceLock.Load(systemUid); ok {
		return m.(*sync.Mutex)
	}
	m, _ := systemNetworkInterfaceLock.LoadOrStore(systemUid, &sync.Mutex{})
	return m.(*sync.Mutex)
}
