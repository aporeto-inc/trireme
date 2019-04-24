package iptablesctrl

import (
	"net"
	"strings"

	"github.com/aporeto-inc/go-ipset/ipset"
	provider "go.aporeto.io/trireme-lib/controller/pkg/aclprovider"
	"go.uber.org/zap"
)

const (
	ipv6String = "v6-"
)

type ipv6 struct {
	ipt provider.IptablesProvider
}

var IPv6Disabled bool
var ipsetV6Param *ipset.Params

func init() {
	ipsetV6Param = &ipset.Params{HashFamily: "inet6"}
	IPv6Disabled = true
}

func GetIPv6Impl() (*ipv6, error) {
	ipt, err := provider.NewGoIPTablesProviderV6([]string{"mangle"})
	if err != nil {
		zap.L().Error("Unable to initialize ipv6 iptables :%s", zap.Error(err))
	}

	return &ipv6{ipt: ipt}, nil
}

func (i *ipv6) GetIPSetPrefix() string {
	return chainPrefix + ipv6String
}

func (i *ipv6) GetIPSetParam() *ipset.Params {
	return ipsetV6Param
}

func (i *ipv6) IPFilter() func(net.IP) bool {
	ipv6Filter := func(ip net.IP) bool {
		if ip.To4() == nil {
			return true
		}

		return false
	}

	return ipv6Filter
}

func (i *ipv6) GetDefaultIP() string {
	return "::/0"
}

func (i *ipv6) NeedICMP() bool {
	return true
}

func (i *ipv6) ProtocolAllowed(proto string) bool {
	if strings.ToLower(proto) == "icmp" {
		return false
	}

	return true
}

func (i *ipv6) Append(table, chain string, rulespec ...string) error {
	if IPv6Disabled || i.ipt == nil {
		return nil
	}

	return i.ipt.Append(table, chain, rulespec...)
}

func (i *ipv6) Insert(table, chain string, pos int, rulespec ...string) error {
	if IPv6Disabled || i.ipt == nil {
		return nil
	}

	return i.ipt.Insert(table, chain, pos, rulespec...)
}

func (i *ipv6) ListChains(table string) ([]string, error) {
	if IPv6Disabled || i.ipt == nil {
		return nil, nil
	}

	return i.ipt.ListChains(table)
}

func (i *ipv6) ClearChain(table, chain string) error {
	if IPv6Disabled || i.ipt == nil {
		return nil
	}

	return i.ipt.ClearChain(table, chain)
}

func (i *ipv6) DeleteChain(table, chain string) error {
	if IPv6Disabled || i.ipt == nil {
		return nil
	}

	return i.ipt.DeleteChain(table, chain)
}

func (i *ipv6) NewChain(table, chain string) error {
	if IPv6Disabled || i.ipt == nil {
		return nil
	}

	return i.ipt.NewChain(table, chain)
}

func (i *ipv6) Commit() error {
	if IPv6Disabled || i.ipt == nil {
		return nil
	}

	return i.ipt.Commit()
}

func (i *ipv6) Delete(table, chain string, rulespec ...string) error {
	if IPv6Disabled || i.ipt == nil {
		return nil
	}

	return i.ipt.Delete(table, chain, rulespec...)
}

func (i *ipv6) RetrieveTable() map[string]map[string][]string {
	return i.ipt.RetrieveTable()
}
