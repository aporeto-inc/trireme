// +build windows

package windows

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/DavidGamba/go-getoptions"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/packet"
	"go.aporeto.io/enforcerd/trireme-lib/utils/frontman"
	"golang.org/x/net/ipv6"
)

// WindowsRuleRange represents a range of values for a rule
type WindowsRuleRange struct { // nolint:golint // ignore type name stutters
	Start int
	End   int
}

// WindowsRuleIcmpMatch represents parameters for an ICMP match
type WindowsRuleIcmpMatch struct { // nolint:golint // ignore type name stutters
	IcmpType      int
	IcmpCodeRange *WindowsRuleRange
	Nomatch       bool
}

// WindowsRuleMatchSet represents result of parsed --match-set
type WindowsRuleMatchSet struct { // nolint:golint // ignore type name stutters
	MatchSetName    string
	MatchSetNegate  bool
	MatchSetDstIP   bool
	MatchSetDstPort bool
	MatchSetSrcIP   bool
	MatchSetSrcPort bool
}

// WindowsRuleSpec represents result of parsed iptables rule
type WindowsRuleSpec struct { // nolint:golint // ignore type name stutters
	Protocol                   int
	Action                     int // FilterAction (allow, drop, nfq, proxy)
	ProxyPort                  int
	Mark                       int
	Log                        bool
	LogPrefix                  string
	GroupID                    int
	ProcessID                  int
	ProcessIncludeChildren     bool
	ProcessIncludeChildrenOnly bool
	MatchSrcPort               []*WindowsRuleRange
	MatchDstPort               []*WindowsRuleRange
	MatchBytesNoMatch          bool
	MatchBytes                 []byte
	MatchBytesOffset           int
	MatchSet                   []*WindowsRuleMatchSet
	IcmpMatch                  []*WindowsRuleIcmpMatch
	TCPFlags                   uint8
	TCPFlagsMask               uint8
	TCPFlagsSpecified          bool
	TCPOption                  uint8
	TCPOptionSpecified         bool
	GotoFilterName             string
	FlowMarkNoMatch            bool
	FlowMark                   int
}

// MakeRuleSpecText converts a WindowsRuleSpec back into a string for an iptables rule
func MakeRuleSpecText(winRuleSpec *WindowsRuleSpec, validate bool) (string, error) {

	rulespec := ""

	if winRuleSpec.Protocol > 0 && winRuleSpec.Protocol < math.MaxUint8 {
		rulespec += fmt.Sprintf("-p %d ", winRuleSpec.Protocol)
	}

	if len(winRuleSpec.MatchBytes) > 0 {
		if winRuleSpec.MatchBytesNoMatch {
			rulespec += fmt.Sprintf("-m string --string ! %s --offset %d ", string(winRuleSpec.MatchBytes), winRuleSpec.MatchBytesOffset)
		} else {
			rulespec += fmt.Sprintf("-m string --string %s --offset %d ", string(winRuleSpec.MatchBytes), winRuleSpec.MatchBytesOffset)
		}
	}

	if len(winRuleSpec.MatchSrcPort) > 0 {
		rulespec += "--sports "
		for i, pr := range winRuleSpec.MatchSrcPort {
			rulespec += strconv.Itoa(pr.Start)
			if pr.Start != pr.End {
				rulespec += fmt.Sprintf(":%d", pr.End)
			}
			if i+1 < len(winRuleSpec.MatchSrcPort) {
				rulespec += ","
			}
		}
		rulespec += " "
	}

	if len(winRuleSpec.MatchDstPort) > 0 {
		rulespec += "--dports "
		for i, pr := range winRuleSpec.MatchDstPort {
			rulespec += strconv.Itoa(pr.Start)
			if pr.Start != pr.End {
				rulespec += fmt.Sprintf(":%d", pr.End)
			}
			if i+1 < len(winRuleSpec.MatchDstPort) {
				rulespec += ","
			}
		}
		rulespec += " "
	}

	if len(winRuleSpec.MatchSet) > 0 {
		for _, ms := range winRuleSpec.MatchSet {
			rulespec += "-m set "
			if ms.MatchSetNegate {
				rulespec += "! "
			}
			rulespec += fmt.Sprintf("--match-set %s ", ms.MatchSetName)
			if ms.MatchSetSrcIP {
				rulespec += "srcIP"
				if ms.MatchSetSrcPort || ms.MatchSetDstPort {
					rulespec += ","
				}
			} else if ms.MatchSetDstIP {
				rulespec += "dstIP"
				if ms.MatchSetSrcPort || ms.MatchSetDstPort {
					rulespec += ","
				}
			}
			if ms.MatchSetSrcPort {
				rulespec += "srcPort"
			} else if ms.MatchSetDstPort {
				rulespec += "dstPort"
			}
			rulespec += " "
		}
	}

	if len(winRuleSpec.IcmpMatch) > 0 {
		for _, im := range winRuleSpec.IcmpMatch {
			if im.Nomatch {
				rulespec += "--icmp-type nomatch"
			} else {
				rulespec += fmt.Sprintf("--icmp-type %d", im.IcmpType)
				if im.IcmpCodeRange != nil {
					rulespec += fmt.Sprintf("/%d", im.IcmpCodeRange.Start)
					if im.IcmpCodeRange.Start != im.IcmpCodeRange.End {
						rulespec += fmt.Sprintf(":%d", im.IcmpCodeRange.End)
					}
				}
			}
			rulespec += " "
		}
	}

	if winRuleSpec.TCPFlagsSpecified {
		rulespec += fmt.Sprintf("--tcp-flags %d,%d ", winRuleSpec.TCPFlagsMask, winRuleSpec.TCPFlags)
	}

	if winRuleSpec.TCPOptionSpecified {
		rulespec += fmt.Sprintf("--tcp-option %d ", winRuleSpec.TCPOption)
	}

	if winRuleSpec.FlowMark != 0 {
		if winRuleSpec.FlowMarkNoMatch {
			rulespec += fmt.Sprintf("-m connmark --mark ! %d ", winRuleSpec.FlowMark)
		} else {
			rulespec += fmt.Sprintf("-m connmark --mark %d ", winRuleSpec.FlowMark)
		}
	}

	switch winRuleSpec.Action {
	case frontman.FilterActionAllow:
		rulespec += "-j ACCEPT "
	case frontman.FilterActionAllowOnce:
		rulespec += "-j ACCEPT_ONCE "
	case frontman.FilterActionBlock:
		rulespec += "-j DROP "
	case frontman.FilterActionProxy:
		rulespec += fmt.Sprintf("-j REDIRECT --to-ports %d ", winRuleSpec.ProxyPort)
	case frontman.FilterActionNfq:
		rulespec += fmt.Sprintf("-j NFQUEUE -j MARK %d ", winRuleSpec.Mark)
	case frontman.FilterActionForceNfq:
		rulespec += fmt.Sprintf("-j NFQUEUE_FORCE -j MARK %d ", winRuleSpec.Mark)
	case frontman.FilterActionGotoFilter:
		rulespec += fmt.Sprintf("-j %s ", winRuleSpec.GotoFilterName)
	case frontman.FilterActionSetMark:
		rulespec += fmt.Sprintf("-j CONNMARK --set-mark %d", winRuleSpec.Mark)
	}
	if winRuleSpec.Log {
		rulespec += fmt.Sprintf("-j NFLOG --nflog-group %d --nflog-prefix \"%s\" ", winRuleSpec.GroupID, winRuleSpec.LogPrefix)
	}

	if winRuleSpec.ProcessID > 0 {
		rulespec += fmt.Sprintf("-m owner --pid-owner %d ", winRuleSpec.ProcessID)
		if winRuleSpec.ProcessIncludeChildrenOnly {
			rulespec += "--pid-childrenonly "
		} else if winRuleSpec.ProcessIncludeChildren {
			rulespec += "--pid-children "
		}
	}

	rulespec = strings.TrimSpace(rulespec)
	if validate {
		if _, err := ParseRuleSpec(strings.Split(rulespec, " ")...); err != nil {
			return "", err
		}
	}
	return rulespec, nil
}

// ParsePortString parses comma-separated list of port or port ranges
func ParsePortString(portString string) ([]*WindowsRuleRange, error) {
	var result []*WindowsRuleRange
	if portString != "" {
		portList := strings.Split(portString, ",")
		for _, portListItem := range portList {
			portEnd := 0
			portStart, err := strconv.Atoi(portListItem)
			if err != nil {
				portRange := strings.SplitN(portListItem, ":", 2)
				if len(portRange) != 2 {
					return nil, errors.New("invalid port string")
				}
				portStart, err = strconv.Atoi(portRange[0])
				if err != nil {
					return nil, errors.New("invalid port string")
				}
				portEnd, err = strconv.Atoi(portRange[1])
				if err != nil {
					return nil, errors.New("invalid port string")
				}
			}
			if portEnd == 0 {
				portEnd = portStart
			}
			result = append(result, &WindowsRuleRange{portStart, portEnd})
		}
	}
	return result, nil
}

// ReduceIcmpProtoString will look at policyRestrictions and return a rulespec substring for matching.
// represents the logic: "icmpProtoTypeCode and (policyRestrictions[0] or policyRestrictions[1] or...)"
// can return empty list if there is a proto match with no restriction.
// will return error if there is no intersection.
func ReduceIcmpProtoString(icmpProtoTypeCode string, policyRestrictions []string) ([]string, error) {

	if len(policyRestrictions) == 0 {
		return TransformIcmpProtoString(icmpProtoTypeCode), nil
	}

	splitIt := func(p string) (string, []*WindowsRuleIcmpMatch, error) {
		var c []*WindowsRuleIcmpMatch
		var err error
		parts := strings.SplitN(p, "/", 2)
		switch len(parts) {
		case 2:
			c, err = ParseIcmpTypeCode(parts[1])
			if err != nil {
				return "", nil, err
			}
			fallthrough
		case 1:
			return parts[0], c, nil
		default:
			return "", nil, fmt.Errorf("invalid icmpProtoTypeCode: %s", icmpProtoTypeCode)
		}
	}

	normalizeProto := func(p string) string {
		switch strings.ToLower(p) {
		case "1":
			return "icmp"
		case "58", "icmp6":
			return "icmpv6"
		}
		return p
	}

	proto, criteria, err := splitIt(icmpProtoTypeCode)
	if err != nil {
		return nil, err
	}

	var positiveMatch bool
	result := make([]string, 0, len(policyRestrictions))
	for _, restriction := range policyRestrictions {
		protoR, criteriaR, err := splitIt(restriction)
		if err != nil {
			return nil, err
		}
		// proto should match
		if proto != protoR && normalizeProto(proto) != normalizeProto(protoR) {
			continue
		}
		if len(criteriaR) == 0 {
			// no restriction
			result = append(result, TransformIcmpProtoString(icmpProtoTypeCode)...)
			positiveMatch = true
			continue
		}
		if len(criteria) == 0 {
			// restriction takes effect
			result = append(result, TransformIcmpProtoString(restriction)...)
			positiveMatch = true
			continue
		}

		if criteria[0].IcmpType != criteriaR[0].IcmpType {
			// types don't match
			continue
		}

		var ranges, rangesR []WindowsRuleRange
		for _, c := range criteria {
			if c.IcmpCodeRange != nil {
				ranges = append(ranges, *c.IcmpCodeRange)
			}
		}
		for _, c := range criteriaR {
			if c.IcmpCodeRange != nil {
				rangesR = append(rangesR, *c.IcmpCodeRange)
			}
		}

		if len(rangesR) == 0 {
			// no code restriction
			result = append(result, TransformIcmpProtoString(icmpProtoTypeCode)...)
			positiveMatch = true
			continue
		}
		if len(ranges) == 0 {
			// use restriction
			result = append(result, TransformIcmpProtoString(restriction)...)
			positiveMatch = true
			continue
		}

		// intersect the code restrictions
		combined := make([]*WindowsRuleRange, 0, len(ranges)+len(rangesR))
		sort.Slice(ranges, func(i, j int) bool {
			return ranges[i].Start < ranges[j].Start
		})
		sort.Slice(rangesR, func(i, j int) bool {
			return rangesR[i].Start < rangesR[j].Start
		})
		for i, j := 0, 0; i < len(ranges) && j < len(rangesR); {
			a, b := ranges[i], rangesR[j]
			// find max of the mins
			maxOfMins := a.Start
			if b.Start > maxOfMins {
				maxOfMins = b.Start
			}
			// find smaller max, and check if it's less than the other min.
			// if not then the intersection is [max(min1,min2),smallermax]
			if a.End < b.End {
				if a.End >= b.Start {
					combined = append(combined, &WindowsRuleRange{Start: maxOfMins, End: a.End})
				}
			} else {
				if b.End >= a.Start {
					combined = append(combined, &WindowsRuleRange{Start: maxOfMins, End: b.End})
				}
			}
			// advance
			if a.End <= b.End {
				i++
			}
			if b.End <= a.End {
				j++
			}
		}

		if len(combined) == 0 {
			// no intersection
			continue
		}

		codeString := ""
		for i, c := range combined {
			if i > 0 {
				codeString += ","
			}
			codeString += fmt.Sprintf("%d:%d", c.Start, c.End)
		}
		combinedString := fmt.Sprintf("%s/%d/%s", proto, criteria[0].IcmpType, codeString)
		result = append(result, TransformIcmpProtoString(combinedString)...)
		positiveMatch = true
	}

	if !positiveMatch {
		return nil, errors.New("policy restrictions do not match")
	}
	return result, nil
}

// TransformIcmpProtoString parses icmp/type/code string coming from ACL rule
// and returns a rulespec subsection
func TransformIcmpProtoString(icmpProtoTypeCode string) []string {
	parts := strings.SplitN(icmpProtoTypeCode, "/", 2)
	if len(parts) != 2 {
		return nil
	}
	typeCodeString := strings.TrimSpace(parts[1])
	if typeCodeString == "" {
		return nil
	}
	return []string{"--icmp-type", typeCodeString}
}

// GetIcmpNoMatch returns a rulespec subsection to indicate that there should be no match
func GetIcmpNoMatch() []string {
	return []string{"--icmp-type", "nomatch"}
}

// ParseIcmpTypeCode parses --icmp-type option
// string is of the form type/code:code,code,code:code
func ParseIcmpTypeCode(icmpTypeCode string) ([]*WindowsRuleIcmpMatch, error) {

	if icmpTypeCode == "" {
		return nil, nil
	}

	if strings.EqualFold(icmpTypeCode, "nomatch") {
		return []*WindowsRuleIcmpMatch{{Nomatch: true}}, nil
	}

	var result []*WindowsRuleIcmpMatch

	parts := strings.SplitN(icmpTypeCode, "/", 2)
	if len(parts) == 0 {
		return nil, nil
	}
	icmpType, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	if icmpType < 0 || icmpType > math.MaxUint8 {
		return nil, errors.New("ICMP type out of range")
	}
	if len(parts) > 1 {
		// parse codes, comma-separated
		for _, code := range strings.Split(parts[1], ",") {
			// parse code range
			codeLower, codeUpper := -1, -1
			codeRange := strings.SplitN(code, ":", 2)
			if len(codeRange) > 0 {
				codeLower, err = strconv.Atoi(codeRange[0])
				if err != nil {
					return nil, err
				}
				codeUpper = codeLower
			}
			if len(codeRange) > 1 {
				codeUpper, err = strconv.Atoi(codeRange[1])
				if err != nil {
					return nil, err
				}
			}
			if codeLower < 0 || codeLower > math.MaxUint8 {
				return nil, errors.New("ICMP code out of range")
			}
			if codeUpper < 0 || codeUpper > math.MaxUint8 || codeUpper < codeLower {
				return nil, errors.New("ICMP code out of range")
			}
			result = append(result, &WindowsRuleIcmpMatch{
				IcmpType:      icmpType,
				IcmpCodeRange: &WindowsRuleRange{codeLower, codeUpper},
			})
		}
	}
	if len(result) == 0 {
		result = append(result, &WindowsRuleIcmpMatch{IcmpType: icmpType})
	}

	return result, nil
}

// ParseRuleSpec parses a windows iptable rule
func ParseRuleSpec(rulespec ...string) (*WindowsRuleSpec, error) {

	opt := getoptions.New()

	protocolOpt := opt.String("p", "")
	sPortOpt := opt.String("sports", "")
	dPortOpt := opt.String("dports", "")
	actionOpt := opt.StringSlice("j", 1, 10, opt.Required())
	modeOpt := opt.StringSlice("m", 1, 10)
	matchSetOpt := opt.StringSlice("match-set", 2, 10)
	matchStringOpt := opt.StringSlice("string", 1, 2)
	matchStringOffsetOpt := opt.Int("offset", 0)
	matchOwnerPidOpt := opt.Int("pid-owner", 0)
	matchOwnerPidChildrenOpt := opt.Bool("pid-children", false)
	matchOwnerPidChildrenOnlyOpt := opt.Bool("pid-childrenonly", false)
	redirectPortOpt := opt.Int("to-ports", 0)
	stateOpt := opt.String("state", "")
	opt.String("match", "") // "--match multiport" ignored
	groupIDOpt := opt.Int("nflog-group", 0)
	logPrefixOpt := opt.String("nflog-prefix", "")
	icmpTypeOpt := opt.StringSlice("icmp-type", 1, 20)
	tcpFlags := opt.StringSlice("tcp-flags", 1, 1)
	tcpOption := opt.StringSlice("tcp-option", 1, 1)
	markOpt := opt.StringSlice("mark", 1, 2)
	setMarkOpt := opt.Int("set-mark", 0)

	_, err := opt.Parse(rulespec)
	if err != nil {
		return nil, err
	}

	result := &WindowsRuleSpec{}

	// protocol
	isProtoAnyRule := false
	switch strings.ToLower(*protocolOpt) {
	case "tcp":
		result.Protocol = packet.IPProtocolTCP
	case "udp":
		result.Protocol = packet.IPProtocolUDP
	case "icmp":
		result.Protocol = 1
	case "icmpv6":
		result.Protocol = ipv6.ICMPType(0).Protocol()
	case "": // not specified = all
		fallthrough
	case "all":
		result.Protocol = -1
		isProtoAnyRule = true
	default:
		result.Protocol, err = strconv.Atoi(*protocolOpt)
		if err != nil {
			return nil, errors.New("rulespec not valid: invalid protocol")
		}
		if result.Protocol < 0 || result.Protocol > math.MaxUint8 {
			return nil, errors.New("rulespec not valid: invalid protocol")
		}
		// iptables man page says protocol zero is equivalent to 'all' (sorry, IPv6 Hop-by-Hop Option)
		if result.Protocol == 0 {
			result.Protocol = -1
		}
	}

	if len(*tcpFlags) == 1 {
		parts := strings.SplitN((*tcpFlags)[0], ",", 2)
		if len(parts) == 0 {
			return nil, errors.New("rulespec not valid: invalid tcp-flags")
		}
		mask, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		if mask < 0 || mask > math.MaxUint8 {
			return nil, errors.New("TCP mask out of range")
		}
		flags, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		if flags < 0 || flags > math.MaxUint8 {
			return nil, errors.New("TCP flags out of range")
		}

		result.TCPFlags = uint8(flags)
		result.TCPFlagsMask = uint8(mask)
		result.TCPFlagsSpecified = true
	}

	if len(*tcpOption) == 1 {
		option, err := strconv.Atoi((*tcpOption)[0])
		if err != nil {
			return nil, err
		}
		if option < 0 || option > math.MaxUint8 {
			return nil, errors.New("TCP option out of range")
		}
		result.TCPOption = uint8(option)
		result.TCPOptionSpecified = true
	}

	for i := 0; i < len(*icmpTypeOpt); i++ {
		im, err := ParseIcmpTypeCode((*icmpTypeOpt)[i])
		if err != nil {
			return nil, fmt.Errorf("rulespec not valid: %s", err.Error())
		}
		result.IcmpMatch = append(result.IcmpMatch, im...)
	}

	// src/dest port: either port or port range or list of such
	result.MatchSrcPort, err = ParsePortString(*sPortOpt)
	if err != nil {
		return nil, errors.New("rulespec not valid: invalid match port")
	}
	result.MatchDstPort, err = ParsePortString(*dPortOpt)
	if err != nil {
		return nil, errors.New("rulespec not valid: invalid match port")
	}

	// -m options
	for i, modeOptSetNum := 0, 0; i < len(*modeOpt); i++ {
		switch (*modeOpt)[i] {
		case "set":
			matchSet := &WindowsRuleMatchSet{}
			// see if negate of --match-set occurred
			if i+1 < len(*modeOpt) && (*modeOpt)[i+1] == "!" {
				matchSet.MatchSetNegate = true
				i++
			}
			// now check corresponding match-set by index
			matchSetIndex := 2 * modeOptSetNum
			modeOptSetNum++
			if matchSetIndex+1 >= len(*matchSetOpt) {
				return nil, errors.New("rulespec not valid: --match-set not found for -m set")
			}
			// first part is the ipset name
			matchSet.MatchSetName = (*matchSetOpt)[matchSetIndex]
			// second part is the dst/src match specifier
			ipPortSpecLower := strings.ToLower((*matchSetOpt)[matchSetIndex+1])
			if strings.HasPrefix(ipPortSpecLower, "dstip") {
				matchSet.MatchSetDstIP = true
			} else if strings.HasPrefix(ipPortSpecLower, "srcip") {
				matchSet.MatchSetSrcIP = true
			}
			if strings.HasSuffix(ipPortSpecLower, "dstport") {
				matchSet.MatchSetDstPort = true
				if result.Protocol < 1 {
					return nil, errors.New("rulespec not valid: ipset match on port requires protocol be set")
				}
			} else if strings.HasSuffix(ipPortSpecLower, "srcport") {
				matchSet.MatchSetSrcPort = true
				if result.Protocol < 1 {
					return nil, errors.New("rulespec not valid: ipset match on port requires protocol be set")
				}
			}
			if !matchSet.MatchSetDstIP && !matchSet.MatchSetDstPort && !matchSet.MatchSetSrcIP && !matchSet.MatchSetSrcPort {
				// look for acl-created iptables-conforming match on 'dst' or 'src'.
				// a dst or src by itself we take to mean match both. otherwise, we take it as ip-match,port-match.
				if strings.HasPrefix(ipPortSpecLower, "dst") {
					matchSet.MatchSetDstIP = true
				} else if strings.HasPrefix(ipPortSpecLower, "src") {
					matchSet.MatchSetSrcIP = true
				}
				if strings.HasSuffix(ipPortSpecLower, "dst") && !isProtoAnyRule {
					matchSet.MatchSetDstPort = true
					if result.Protocol < 1 {
						return nil, errors.New("rulespec not valid: ipset match on port requires protocol be set")
					}
				} else if strings.HasSuffix(ipPortSpecLower, "src") && !isProtoAnyRule {
					matchSet.MatchSetSrcPort = true
					if result.Protocol < 1 {
						return nil, errors.New("rulespec not valid: ipset match on port requires protocol be set")
					}
				}
			}
			if !matchSet.MatchSetDstIP && !matchSet.MatchSetDstPort && !matchSet.MatchSetSrcIP && !matchSet.MatchSetSrcPort {
				return nil, errors.New("rulespec not valid: ipset match needs ip/port specifier")
			}
			result.MatchSet = append(result.MatchSet, matchSet)

		case "string":
			nomatch := false
			matchString := ""
			switch len(*matchStringOpt) {
			case 1:
				matchString = (*matchStringOpt)[0]
			case 2:
				if (*matchStringOpt)[0] != "!" {
					return nil, errors.New("rulespec not valid: match string ! option invalid")
				}
				nomatch = true
				matchString = (*matchStringOpt)[1]
			}
			if len(matchString) == 0 {
				return nil, errors.New("rulespec not valid: no match string given")
			}
			result.MatchBytesNoMatch = nomatch
			result.MatchBytes = []byte(matchString)
			result.MatchBytesOffset = *matchStringOffsetOpt

		case "owner":
			if *matchOwnerPidOpt <= 0 {
				return nil, errors.New("rulespec not valid: valid pid-owner needed")
			}
			result.ProcessID = *matchOwnerPidOpt
			if *matchOwnerPidChildrenOnlyOpt {
				result.ProcessIncludeChildrenOnly = true
			}
			if *matchOwnerPidChildrenOpt {
				if *matchOwnerPidChildrenOnlyOpt {
					return nil, errors.New("rulespec not valid: cannot have both --pid-childrenonly and --pid-children")
				}
				result.ProcessIncludeChildren = true
			}

		case "state":
			if result.Protocol == packet.IPProtocolTCP && *stateOpt == "NEW" {
				result.TCPFlags = 2
				result.TCPFlagsMask = 18
				result.TCPFlagsSpecified = true
			}

		case "connmark":
			nomatch := false
			flowMarkString := ""
			switch len(*markOpt) {
			case 1:
				flowMarkString = (*markOpt)[0]
			case 2:
				if (*markOpt)[0] != "!" {
					return nil, errors.New("rulespec not valid: flowmark ! option invalid")
				}
				nomatch = true
				flowMarkString = (*markOpt)[1]
			}
			value, err := strconv.Atoi(flowMarkString)
			if err != nil {
				return nil, errors.New("rulespec not valid: flowmmark should be int")
			}
			result.FlowMarkNoMatch = nomatch
			result.FlowMark = value

		default:
			return nil, errors.New("rulespec not valid: unknown -m option")
		}
	}

	// action: either NFQUEUE, REDIRECT, MARK, ACCEPT, DROP, NFLOG
	for i := 0; i < len(*actionOpt); i++ {
		// CRITICAL: If you add a case here, you need to update fixRuleSpec in iptablesprovider_windows.go
		switch (*actionOpt)[i] {
		case "NFQUEUE":
			result.Action = frontman.FilterActionNfq
		case "NFQUEUE_FORCE":
			result.Action = frontman.FilterActionForceNfq
		case "REDIRECT":
			result.Action = frontman.FilterActionProxy
		case "ACCEPT":
			result.Action = frontman.FilterActionAllow
		case "ACCEPT_ONCE":
			result.Action = frontman.FilterActionAllowOnce
		case "DROP":
			result.Action = frontman.FilterActionBlock
		case "CONNMARK":
			result.Action = frontman.FilterActionSetMark
			result.Mark = *setMarkOpt
		case "MARK":
			i++
			if i >= len(*actionOpt) {
				return nil, errors.New("rulespec not valid: no mark given")
			}
			result.Mark, err = strconv.Atoi((*actionOpt)[i])
			if err != nil {
				return nil, errors.New("rulespec not valid: mark should be int32")
			}
		case "NFLOG":
			// if no other action specified, it will default to 'continue' action (zero)
			result.Log = true
			result.LogPrefix = *logPrefixOpt
			result.GroupID = *groupIDOpt
		default:
			if i >= len(*actionOpt) {
				return nil, errors.New("rulespec not valid: no goto filter given")
			}
			result.Action = frontman.FilterActionGotoFilter
			result.GotoFilterName = (*actionOpt)[i]
		}
	}

	if result.Mark == 0 && (result.Action == frontman.FilterActionNfq || result.Action == frontman.FilterActionForceNfq) {
		return nil, errors.New("rulespec not valid: nfq action needs to set mark")
	}

	if result.Mark == 0 && (result.Action == frontman.FilterActionSetMark) {
		return nil, errors.New("rulespec not valid: setmark action needs to set mark")
	}

	// redirect port
	result.ProxyPort = *redirectPortOpt
	if result.Action == frontman.FilterActionProxy && result.ProxyPort == 0 {
		return nil, errors.New("rulespec not valid: no redirect port given")
	}

	return result, nil
}

// Equal compares a WindowsRuleMatchSet to another for equality
func (w *WindowsRuleMatchSet) Equal(other *WindowsRuleMatchSet) bool {
	if other == nil {
		return false
	}
	return w.MatchSetName == other.MatchSetName &&
		w.MatchSetNegate == other.MatchSetNegate &&
		w.MatchSetDstIP == other.MatchSetDstIP &&
		w.MatchSetDstPort == other.MatchSetDstPort &&
		w.MatchSetSrcIP == other.MatchSetSrcIP &&
		w.MatchSetSrcPort == other.MatchSetSrcPort
}

// Equal compares a WindowsRuleRange to another for equality
func (w *WindowsRuleRange) Equal(other *WindowsRuleRange) bool {
	if other == nil {
		return false
	}
	return w.Start == other.Start && w.End == other.End
}

// Equal compares a WindowsRuleIcmpMatch to another for equality
func (w *WindowsRuleIcmpMatch) Equal(other *WindowsRuleIcmpMatch) bool {
	if other == nil {
		return false
	}
	if w.Nomatch != other.Nomatch {
		return false
	}
	if w.IcmpType != other.IcmpType {
		return false
	}
	if w.IcmpCodeRange != nil {
		if !w.IcmpCodeRange.Equal(other.IcmpCodeRange) {
			return false
		}
	} else if other.IcmpCodeRange != nil {
		return false
	}
	return true
}

// Equal compares a WindowsRuleSpec to another for equality
func (w *WindowsRuleSpec) Equal(other *WindowsRuleSpec) bool {
	if other == nil {
		return false
	}
	equalSoFar := w.Protocol == other.Protocol &&
		w.Action == other.Action &&
		w.ProxyPort == other.ProxyPort &&
		w.Mark == other.Mark &&
		w.Log == other.Log &&
		w.LogPrefix == other.LogPrefix &&
		w.GroupID == other.GroupID &&
		w.ProcessID == other.ProcessID &&
		w.ProcessIncludeChildren == other.ProcessIncludeChildren &&
		w.ProcessIncludeChildrenOnly == other.ProcessIncludeChildrenOnly &&
		w.MatchBytesNoMatch == other.MatchBytesNoMatch &&
		w.MatchBytesOffset == other.MatchBytesOffset &&
		bytes.Equal(w.MatchBytes, other.MatchBytes) &&
		len(w.IcmpMatch) == len(other.IcmpMatch) &&
		len(w.MatchSrcPort) == len(other.MatchSrcPort) &&
		len(w.MatchDstPort) == len(other.MatchDstPort) &&
		len(w.MatchSet) == len(other.MatchSet) &&
		w.TCPFlags == other.TCPFlags &&
		w.TCPFlagsMask == other.TCPFlagsMask &&
		w.TCPFlagsSpecified == other.TCPFlagsSpecified &&
		w.TCPOption == other.TCPOption &&
		w.TCPOptionSpecified == other.TCPOptionSpecified &&
		w.FlowMark == other.FlowMark &&
		w.FlowMarkNoMatch == other.FlowMarkNoMatch
	if !equalSoFar {
		return false
	}
	// we checked lengths above, but now continue to compare slices for equality.
	// note: we assume equal slices have elements in the same order.
	for i := 0; i < len(w.MatchSrcPort); i++ {
		if w.MatchSrcPort[i] == nil {
			if other.MatchSrcPort[i] != nil {
				return false
			}
			continue
		}
		if !w.MatchSrcPort[i].Equal(other.MatchSrcPort[i]) {
			return false
		}
	}
	for i := 0; i < len(w.MatchDstPort); i++ {
		if w.MatchDstPort[i] == nil {
			if other.MatchDstPort[i] != nil {
				return false
			}
			continue
		}
		if !w.MatchDstPort[i].Equal(other.MatchDstPort[i]) {
			return false
		}
	}
	for i := 0; i < len(w.MatchSet); i++ {
		if w.MatchSet[i] == nil {
			if other.MatchSet[i] != nil {
				return false
			}
			continue
		}
		if !w.MatchSet[i].Equal(other.MatchSet[i]) {
			return false
		}
	}
	for i := 0; i < len(w.IcmpMatch); i++ {
		if w.IcmpMatch[i] == nil {
			if other.IcmpMatch[i] != nil {
				return false
			}
			continue
		}
		if !w.IcmpMatch[i].Equal(other.IcmpMatch[i]) {
			return false
		}
	}
	return true
}
