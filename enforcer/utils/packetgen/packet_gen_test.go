//Packet generator test for TCP and IP
//Change values of TCP header fields
//Still in beta version
//Updates are coming soon for more options to IP, hopefully ethernet too
//Test cases are created only for generated packets, not for packets on the wire
package packetgen

import "testing"

//. "github.com/smartystreets/goconvey/convey"

// func SetTCPSyn(pBytes []byte) []byte {
// 	// TCPOffset = pBytes[IPOffset] + TCPHeaderLength
// 	// TCPSynOffset = TCPOffset + 4
// 	// pBytes[TCPSynOffset] |= TCPSynFlag
// 	return pBytes
func init() {
	// fmt.Println()
	//fmt.Println(NewPacket().AddIPLayer("122.1.1.1.", "122.2.3.4"))
	// p := NewPacket()
	// p.AddEthernetLayer("aa:aa:fa:aa:ff:aa", "aa:aa:fa:aa:ff:ff")
	// p.AddIPLayer("164.67.228.152", "10.1.10.76")
	// p.AddTCPLayer(666, 80)
	// // //fmt.Println(p.CreatePacketBuffer())
	// pf := NewTCPPacketFlow("aa:aa:fa:aa:ff:aa", "aa:aa:fa:aa:ff:ff", "192.168.1.1", "192.168.2.2", 666, 80)
	//
	// pf.GenerateTCPFlow(PacketFlowTypeGoodFlow)
	// //pf.GetNthPacket(0).NewTCPPayload("Aporeto's Confidential")
	// fmt.Println(pf.GetNthPacket(0).ToBytes())
	// // pfSyn, ok := pf.GetSynPackets()
	// if ok {
	// 	fmt.Printf("Packet Flow:\n%+v\n", pfSyn)
	// } else {
	// 	fmt.Println("Packet Flow: None")
	// }
	//
	// for i := 0; i < pfSyn.GetNumPackets(); i++ {
	// 	p, ok := pf.GetNthPacket(i)
	// 	if ok {
	// 		fmt.Printf("Packet Bytes: %v\n", p.ToBytes())
	// 		p.SetTCPSyn()
	// 	}
	// }
	//
	// for i := 0; i < len(PacketFlowBytes); i++ {
	// 	pBytes := PacketFlowBytes[i]
	// 	if ok {
	// 		fmt.Printf("Packet Bytes: %v\n", pBytes)
	// 		SetTCPSyn(pBytes)
	// 	}
	// }
	//  fmt.Printf("%X",pf.GetNthPacket(2).ToBytes())
	//fmt.Println(pf.GetSynPackets().ToBytes())
	//p.SetTCPSequenceNumber(2345)
	//p.DisplayTCPPacket()
	// layer.SrcIPstr = "164.67.228.152"
	// layer.DstIPstr = "10.1.10.76"
	// ipLayer := layer.GenerateIPPacket(layer.SrcIPstr, layer.DstIPstr)
	// layer.SrcPort = 666
	// layer.DstPort = 80
	// layer.GenerateTCPPacket(&ipLayer, layer.SrcPort, layer.DstPort)
	// layer.SetSynTrue()
	// layer.SequenceNum = 0
	// //layer.InitTemplate()
	// TCPFlow = layer.GenerateTCPFlow(layer.TemplateFlow)
	// //TCPFlow = layer.GenerateTCPFlowPayload("Aporeto Confidential")

}

//
//func TestSample(t *testing.T) {}
//
// //check th enumber of tcp layers generated
// func TestCount(t *testing.T) {
//
// 	t.Parallel()
//
// 	if len(TCPFlow) != 3 {
// 		t.Error("Cannot generate TCP flow, missing either SYN, SYNACK or ACK packets")
// 	}
//
// }
//
// //check for payload in packets
// func TestForPayloadAvailability(t *testing.T) {
//
// 	t.Parallel()
//
// }
//
// //check if Syn is set for the first packet
// func TestForSYNPacket(t *testing.T) {
//
// 	t.Parallel()
//
// 	if layer.Layers[0].SYN != true && layer.Layers[0].ACK == true {
// 		t.Error("No SYN packet in starting flow")
// 	}
//
// }
//
// //check if ethernet is removed from the layers to support datapath_test
// func TestEthernetPresence(t *testing.T) {
//
// 	t.Parallel()
//
// 	for i, _ := range TCPFlow {
// 		if len(TCPFlow[i]) != 46 {
// 			t.Errorf("Ethernet not supported. Check this layer %d", TCPFlow[i])
// 		}
// 	}
//
// }
//
// //check if the TCP flow is good
// func TestGoodPacketFlow(t *testing.T) {
//
// 	t.Parallel()
//
// 	if layer.Layers[0].SrcPort != layer.SrcPort {
// 		t.Error("unexpected source port")
// 	}
//
// 	if layer.Layers[0].DstPort != layer.DstPort {
// 		t.Error("unexpected destination port")
// 	}
//
// 	if layer.Layers[1].SrcPort != layer.DstPort {
// 		t.Error("wrong SynAck port set")
// 	}
//
// 	if layer.Layers[1].DstPort != layer.SrcPort {
// 		t.Error("wrong SynAck port set")
// 	}
//
// }

func TestTypeInterface(t *testing.T) {
	t.Parallel()

	var PktInterface PacketManipulator = (*Packet)(nil)

	if PktInterface != (*Packet)(nil) {

		t.Error("Packet struct does not implement Pkt Interface")

	}

	var PktFlowInterface PacketFlowManipulator = (*PacketFlow)(nil)
	if PktFlowInterface != (*PacketFlow)(nil) {

		t.Error("PacketFlow struct does not implement PktFlow Interface")

	}

}
