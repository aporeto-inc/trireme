package ipsetctrl

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/aporeto-inc/trireme/policy"
	"github.com/bvandewalle/go-ipset/ipset"
)

const (
	appChainPrefix = "TRIREME-App-"
	netChainPrefix = "TRIREME-Net-"
)

// createACLSets creates the sets for a given PU
func (i *Instance) createACLSets(set string, rules *policy.IPRuleList) error {
	appSet, err := i.ips.NewIpset(set, "hash:net,port", &ipset.Params{})
	if err != nil {
		return fmt.Errorf("Couldn't create IPSet for Trireme: %s", err)
	}

	for _, rule := range rules.Rules {
		if err := appSet.Add(rule.Address+","+rule.Port, 0); err != nil {
			return fmt.Errorf("Couldn't create IPSet for Trireme: %s", err)
		}
	}
	return nil
}

// AddAppSetRule adds an ACL rule to the Set
func (i *Instance) addAppSetRule(set string, ip string) error {

	if err := i.ipt.Insert(
		i.appAckPacketIPTableContext, i.appPacketIPTableSection, 3,
		"-m", "state", "--state", "NEW",
		"-m", "set", "--match-set", set, "dst",
		"-s", ip,
		"-j", "ACCEPT",
	); err != nil {
		log.WithFields(log.Fields{
			"package":                      "iptablesutils",
			"i.appAckPacketIPTableContext": i.appAckPacketIPTableContext,
			"error": err.Error(),
		}).Debug("Error when adding app acl rule")
		return err

	}
	return nil
}

// addNetSetRule
func (i *Instance) addNetSetRule(set string, ip string) error {

	if err := i.ipt.Insert(
		i.netPacketIPTableContext, i.netPacketIPTableSection, 2,
		"-m", "state", "--state", "NEW",
		"-m", "set", "--match-set", set, "src",
		"-d", ip,
		"-j", "ACCEPT",
	); err != nil {
		log.WithFields(log.Fields{
			"package":                   "iptablesutils",
			"i.netPacketIPTableContext": i.netPacketIPTableContext,
			"error":                     err.Error(),
		}).Debug("Error when adding app acl rule")
		return err
	}
	return nil
}

// deleteAppSetRule
func (i *Instance) deleteAppSetRule(set string, ip string) error {

	if err := i.ipt.Delete(
		i.appAckPacketIPTableContext, i.appPacketIPTableSection,
		"-m", "state", "--state", "NEW",
		"-m", "set", "--match-set", set, "dst",
		"-s", ip,
		"-j", "ACCEPT",
	); err != nil {
		log.WithFields(log.Fields{
			"package":                   "iptablesutils",
			"i.netPacketIPTableContext": i.netPacketIPTableContext,
			"chain":                     i.appPacketIPTableSection,
			"error":                     err.Error(),
		}).Debug("Error when removing ingress app acl rule")

	}
	return nil
}

// deleteNetSetRule
func (i *Instance) deleteNetSetRule(set string, ip string) error {

	if err := i.ipt.Delete(
		i.netPacketIPTableContext, i.netPacketIPTableSection,
		"-m", "state", "--state", "NEW",
		"-m", "set", "--match-set", set, "src",
		"-d", ip,
		"-j", "ACCEPT",
	); err != nil {
		log.WithFields(log.Fields{
			"package":                   "iptablesutils",
			"i.netPacketIPTableContext": i.netPacketIPTableContext,
			"chain":                     i.appPacketIPTableSection,
			"error":                     err.Error(),
		}).Debug("Error when removing ingress app acl rule")

	}
	return nil
}

func (i *Instance) deleteSet(set string) error {
	ipSet, err := i.ips.NewIpset(set, "hash:net,port", &ipset.Params{})
	if err != nil {
		return fmt.Errorf("Couldn't create IPSet for Trireme: %s", err)
	}

	ipSet.Destroy()
	return nil
}

// setupIpset sets up an ipset
func (i *Instance) setupIpset(target, container string) error {

	ips, err := i.ips.NewIpset(target, "hash:net", &ipset.Params{})
	if err != nil {
		log.WithFields(log.Fields{
			"package": "supervisor",
			"error":   err.Error(),
		}).Debug("Error creating NewIPSet")
		return fmt.Errorf("Couldn't create IPSet for %s: %s", target, err)
	}

	for _, net := range i.targetNetworks {
		if err = ips.Add(net, 0); err != nil {
			log.WithFields(log.Fields{
				"package": "supervisor",
				"error":   err.Error(),
			}).Debug("Error adding ip to IPSet")
			return fmt.Errorf("Error adding ip %s to %s IPSet: %s", net, target, err)
		}
	}

	i.targetSet = ips

	containerSet, err := i.ips.NewIpset(container, "hash:ip", &ipset.Params{})
	if err != nil {
		log.WithFields(log.Fields{
			"package": "supervisor",
			"error":   err.Error(),
		}).Debug("Error creating NewIPSet")
		return fmt.Errorf("Failed to create container set")
	}

	i.containerSet = containerSet

	return nil
}

func (i *Instance) addContainerToSet(ip string) error {
	if err := i.containerSet.Add(ip, 0); err != nil {
		log.WithFields(log.Fields{
			"package": "supervisor",
			"error":   err.Error(),
		}).Debug("Error adding container to set ")
		return fmt.Errorf("Error adding ip %s to container set : %s", ip, err)
	}
	return nil
}

func (i *Instance) delContainerFromSet(ip string) error {
	if err := i.containerSet.Del(ip); err != nil {
		log.WithFields(log.Fields{
			"package": "supervisor",
			"error":   err.Error(),
		}).Debug("Error adding container to set ")
		return fmt.Errorf("Error adding ip %s to container set : %s", ip, err)
	}
	return nil
}

// addIpsetOption
func (i *Instance) addIpsetOption(ip string) error {

	return i.targetSet.AddOption(ip, "nomatch", 0)
}

// deleteIpsetOption
func (i *Instance) deleteIpsetOption(ip string) error {

	return i.targetSet.Del(ip)
}

// setupTrapRules
func (i *Instance) setupTrapRules(set string) error {

	rules := [][]string{
		// Application Syn and Syn/Ack in RAW
		{
			i.appPacketIPTableContext, i.appPacketIPTableSection,
			"-m", "set", "--match-set", set, "dst",
			"-m", "set", "--match-set", containerSet, "src",
			"-p", "tcp", "--tcp-flags", "FIN,SYN,RST,PSH,URG", "SYN",
			"-j", "NFQUEUE", "--queue-balance", i.applicationQueues,
		},

		// Application Matching Trireme SRC and DST. Established connections.
		{
			i.appAckPacketIPTableContext, i.appPacketIPTableSection,
			"-m", "set", "--match-set", containerSet, "src",
			"-m", "set", "--match-set", set, "dst",
			"-p", "tcp", "--tcp-flags", "FIN,SYN,RST,PSH,URG", "SYN",
			"-j", "ACCEPT",
		},

		// Application Matching Trireme SRC and DST. SYN, SYNACK connections.
		{
			i.appAckPacketIPTableContext, i.appPacketIPTableSection,
			"-m", "set", "--match-set", containerSet, "src",
			"-m", "set", "--match-set", set, "dst",
			"-p", "tcp", "--tcp-flags", "SYN,ACK", "ACK",
			"-m", "connbytes", "--connbytes", ":3", "--connbytes-dir", "original", "--connbytes-mode", "packets",
			"-j", "NFQUEUE", "--queue-balance", i.applicationQueues,
		},

		// Default Drop from Trireme to Network
		{
			i.appAckPacketIPTableContext, i.appPacketIPTableSection,
			"-m", "set", "--match-set", containerSet, "src",
			"-p", "tcp", "-m", "state", "--state", "NEW",
			"-j", "DROP",
		},

		// Network Matching Trireme SRC and DST.
		{
			i.netPacketIPTableContext, i.netPacketIPTableSection,
			"-m", "set", "--match-set", set, "src",
			"-m", "set", "--match-set", containerSet, "dst",
			"-p", "tcp",
			"-m", "connbytes", "--connbytes", ":3", "--connbytes-dir", "original", "--connbytes-mode", "packets",
			"-j", "NFQUEUE", "--queue-balance", i.networkQueues,
		},

		// Default Drop from Network to Trireme.
		{
			i.netPacketIPTableContext, i.netPacketIPTableSection,
			"-m", "set", "--match-set", containerSet, "dst",
			"-p", "tcp", "-m", "state", "--state", "NEW",
			"-j", "DROP",
		},
	}

	for _, tr := range rules {
		if err := i.ipt.Append(tr[0], tr[1], tr[2:]...); err != nil {
			log.WithFields(log.Fields{
				"package": "supervisor",
				"error":   err.Error(),
			}).Debug("Failed to add initial rules for TriremeNet IPSet.")
			return err
		}
	}
	return nil
}

// cleanIPSets cleans all the ipsets
func (i *Instance) cleanIPSets() error {

	i.ipt.ClearChain(i.appPacketIPTableContext, i.appPacketIPTableSection)

	i.ipt.ClearChain(i.appAckPacketIPTableContext, i.appPacketIPTableSection)

	i.ipt.ClearChain(i.netPacketIPTableContext, i.netPacketIPTableSection)

	return i.ips.DestroyAll()
}
