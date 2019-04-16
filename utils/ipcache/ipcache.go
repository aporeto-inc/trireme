package ipcache

import (
	"encoding/binary"
	"net"
)

type IPcache interface {
	Put(net.IP, int, interface{})
	Get(net.IP, int) (interface{}, bool)
	RunVal(func(val interface{}) interface{})
	RunIP(net.IP, func(val interface{}) bool)
}

var (
	ipv4Masks = 32 + 1
	ipv6Masks = 128 + 1
)

type ipcache struct {
	ipv4 []map[uint32]interface{}
	ipv6 []map[[16]byte]interface{}
}

func NewIPCache() *ipcache {
	return &ipcache{
		ipv4: make([]map[uint32]interface{}, ipv4Masks),
		ipv6: make([]map[[16]byte]interface{}, ipv6Masks),
	}
}

// Put takes an argument an ip address, mask and value. It is used as a cache for quick lookup
func (cache *ipcache) Put(ip net.IP, mask int, val interface{}) {
	if ip.To4() != nil {
		ip = ip.To4()

		if cache.ipv4[mask] == nil {
			cache.ipv4[mask] = map[uint32]interface{}{}
		}

		m := cache.ipv4[mask]
		// the following expression is ANDing the ip with the mask
		m[binary.BigEndian.Uint32(ip)&binary.BigEndian.Uint32(net.CIDRMask(mask, 32))] = val
	} else {
		ip = ip.To16()
		if cache.ipv6[mask] == nil {
			cache.ipv6[mask] = map[[16]byte]interface{}{}
		}

		m := cache.ipv6[mask]
		// the following expression is ANDing the ip with the mask
		var maskip [16]byte
		copy(maskip[:], ip.Mask(net.CIDRMask(mask, 128)))
		m[maskip] = val
	}
}

// Get takes an argument the IP address and mask and returns the value that is stored for that key.
func (cache *ipcache) Get(ip net.IP, mask int) (interface{}, bool) {

	if ip.To4() != nil {
		ip = ip.To4()
		m := cache.ipv4[mask]
		if m != nil {
			val, ok := m[binary.BigEndian.Uint32(ip)&binary.BigEndian.Uint32(net.CIDRMask(mask, 32))]

			if ok {
				return val, true
			}
		}
	} else {
		ip = ip.To16()
		m := cache.ipv6[mask]

		if m != nil {
			var maskip [16]byte
			copy(maskip[:], ip.Mask(net.CIDRMask(mask, 128)))
			val, ok := m[maskip]

			if ok {
				return val, true
			}
		}
	}

	return nil, false
}

// RunIP function takes as an argument an IP address and a function. It finds the subnet to which this IP belongs in reverse sorted order of subnet masks
// It then calls the function supplied by the user on the value stored and if it succeeds then it returns.
func (cache *ipcache) RunIP(ip net.IP, f func(val interface{}) bool) {
	if ip.To4() != nil {
		ip = ip.To4()

		for i := 32; i >= 0; i-- {
			m := cache.ipv4[i]
			if m != nil {
				val, ok := m[binary.BigEndian.Uint32(ip)&binary.BigEndian.Uint32(net.CIDRMask(i, 32))]
				if ok && f(val) {
					return
				}
			}
		}
	} else {
		ip = ip.To16()

		for i := 128; i >= 0; i-- {
			m := cache.ipv6[i]
			if m != nil {
				var maskip [16]byte
				copy(maskip[:], ip.Mask(net.CIDRMask(i, 128)))
				val, ok := m[maskip]
				if ok && f(val) {
					return
				}
			}
		}
	}
}

// RunVal takes an argument a function which is called on all the values stored in the
// cache. It updates the cache value with the new value that is returned by the function
func (cache *ipcache) RunVal(f func(val interface{}) interface{}) {
	for mask, m := range cache.ipv4 {
		if m == nil {
			continue
		}

		for ip, val := range m {
			v := f(val)
			if v == nil {
				delete(m, ip)
				continue
			}

			m[ip] = v
		}

		if len(m) == 0 {
			cache.ipv4[mask] = nil
		}
	}

	for mask, m := range cache.ipv6 {
		if m == nil {
			continue
		}

		for ip, val := range m {
			v := f(val)
			if v == nil {
				delete(m, ip)
				continue
			}

			m[ip] = v
		}

		if len(m) == 0 {
			cache.ipv6[mask] = nil
		}
	}

}
