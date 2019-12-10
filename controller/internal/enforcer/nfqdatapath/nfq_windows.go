// +build windows

package nfqdatapath

import (
	"context"
	"fmt"
	"strconv"
	"syscall"
	"unsafe"

	"go.aporeto.io/trireme-lib/controller/internal/enforcer/nfqdatapath/afinetrawsocket"
	"go.aporeto.io/trireme-lib/controller/internal/enforcer/nfqdatapath/nflog"
	"go.aporeto.io/trireme-lib/controller/internal/windows/frontman"
	"go.aporeto.io/trireme-lib/controller/pkg/connection"
	"go.aporeto.io/trireme-lib/controller/pkg/packet"
	"go.aporeto.io/trireme-lib/controller/pkg/pucontext"
	"go.uber.org/zap"
)

func (d *Datapath) startFrontmanPacketFilter(ctx context.Context, nflogger nflog.NFLogger) error {
	driverHandle, err := frontman.GetDriverHandle()
	if err != nil {
		return err
	}

	nflogWin := nflogger.(*nflog.NfLogWindows)

	packetCallback := func(packetInfoPtr, dataPtr uintptr) uintptr {

		packetInfo := *(*frontman.PacketInfo)(unsafe.Pointer(packetInfoPtr))
		packetBytes := (*[1 << 30]byte)(unsafe.Pointer(dataPtr))[:packetInfo.PacketSize:packetInfo.PacketSize]

		var packetType int
		if packetInfo.Outbound != 0 {
			packetType = packet.PacketTypeApplication
		} else {
			packetType = packet.PacketTypeNetwork
		}

		// Parse the packet
		mark := int(packetInfo.Mark)
		parsedPacket, err := packet.New(uint64(packetType), packetBytes, strconv.Itoa(mark), true)

		if parsedPacket.IPProto() == packet.IPProtocolUDP && parsedPacket.SourcePort() == 53 {
			// notify PUs of DNS results
			err := d.dnsProxy.HandleDNSResponsePacket(parsedPacket.GetUDPData(), parsedPacket.SourceAddress(), func(id string) (*pucontext.PUContext, error) {
				puCtx, err1 := d.puFromContextID.Get(id)
				if err1 != nil {
					return nil, err1
				}
				return puCtx.(*pucontext.PUContext), nil
			})
			if err != nil {
				zap.L().Error("Failed to handle DNS response", zap.Error(err))
			}
			// forward packet
			dllRet, _, err := frontman.PacketFilterForwardProc.Call(uintptr(unsafe.Pointer(&packetInfo)), uintptr(unsafe.Pointer(&packetBytes[0])))
			if dllRet == 0 {
				zap.L().Error(fmt.Sprintf("%s failed: %v", frontman.PacketFilterForwardProc.Name, err))
			}
			return 0
		}

		parsedPacket.WindowsMetadata = &afinetrawsocket.WindowsPacketMetadata{PacketInfo: packetInfo, IgnoreFlow: false}
		var processError error
		var tcpConn *connection.TCPConnection
		var udpConn *connection.UDPConnection
		if err != nil {
			parsedPacket.Print(packet.PacketFailureCreate, d.packetLogs)
		} else if parsedPacket.IPProto() == packet.IPProtocolTCP {
			if packetType == packet.PacketTypeNetwork {
				tcpConn, processError = d.processNetworkTCPPackets(parsedPacket)
			} else {
				tcpConn, processError = d.processApplicationTCPPackets(parsedPacket)
			}
		} else if parsedPacket.IPProto() == packet.IPProtocolUDP {
			// process udp packet
			if packetType == packet.PacketTypeNetwork {
				udpConn, processError = d.ProcessNetworkUDPPacket(parsedPacket)
			} else {
				udpConn, processError = d.ProcessApplicationUDPPacket(parsedPacket)
			}
		} else {
			processError = fmt.Errorf("invalid ip protocol: %d", parsedPacket.IPProto())
		}

		if processError != nil {
			if parsedPacket.IPProto() == packet.IPProtocolTCP {
				d.collectTCPPacket(&debugpacketmessage{
					Mark:    mark,
					p:       parsedPacket,
					tcpConn: tcpConn,
					udpConn: nil,
					err:     processError,
					network: packetType == packet.PacketTypeNetwork,
				})
			} else if parsedPacket.IPProto() == packet.IPProtocolUDP {
				d.collectUDPPacket(&debugpacketmessage{
					Mark:    mark,
					p:       parsedPacket,
					tcpConn: nil,
					udpConn: udpConn,
					err:     processError,
					network: packetType == packet.PacketTypeNetwork,
				})
			}
			// drop packet by not forwarding it
			return 0
		}

		// accept the (modified) packet by forwarding it
		var modifiedPacketBytes []byte
		if parsedPacket.IPProto() == packet.IPProtocolTCP {
			modifiedPacketBytes = make([]byte, parsedPacket.IPTotalLen())
			copyIndex := copy(modifiedPacketBytes, parsedPacket.GetBuffer(0))
			copyIndex += copy(modifiedPacketBytes[copyIndex:], parsedPacket.GetTCPOptions())
			copyIndex += copy(modifiedPacketBytes[copyIndex:], parsedPacket.GetTCPData())
			packetInfo.PacketSize = uint32(copyIndex)
		} else {
			modifiedPacketBytes = parsedPacket.GetBuffer(0)
			packetInfo.PacketSize = uint32(len(modifiedPacketBytes))
		}

		if parsedPacket.WindowsMetadata.(*afinetrawsocket.WindowsPacketMetadata).IgnoreFlow {
			packetInfo.IgnoreFlow = 1
		}
		dllRet, _, err := frontman.PacketFilterForwardProc.Call(uintptr(unsafe.Pointer(&packetInfo)), uintptr(unsafe.Pointer(&modifiedPacketBytes[0])))
		if dllRet == 0 {
			zap.L().Error(fmt.Sprintf("%s failed: %v", frontman.PacketFilterForwardProc.Name, err))
		}

		if parsedPacket.IPProto() == packet.IPProtocolTCP {
			d.collectTCPPacket(&debugpacketmessage{
				Mark:    mark,
				p:       parsedPacket,
				tcpConn: tcpConn,
				udpConn: nil,
				err:     nil,
				network: packetType == packet.PacketTypeNetwork,
			})
		} else if parsedPacket.IPProto() == packet.IPProtocolUDP {
			d.collectUDPPacket(&debugpacketmessage{
				Mark:    mark,
				p:       parsedPacket,
				tcpConn: nil,
				udpConn: udpConn,
				err:     nil,
				network: packetType == packet.PacketTypeNetwork,
			})
		}

		return 0
	}

	logCallback := func(logPacketInfoPtr, dataPtr uintptr) uintptr {

		logPacketInfo := *(*frontman.LogPacketInfo)(unsafe.Pointer(logPacketInfoPtr))
		packetHeaderBytes := (*[1 << 30]byte)(unsafe.Pointer(dataPtr))[:logPacketInfo.PacketSize:logPacketInfo.PacketSize]

		err := nflogWin.NfLogHandler(&logPacketInfo, packetHeaderBytes)
		if err != nil {
			zap.L().Error("error in log callback", zap.Error(err))
		}

		return 0
	}

	dllRet, _, err := frontman.PacketFilterStartProc.Call(driverHandle, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Aporeto Enforcer"))),
		syscall.NewCallbackCDecl(packetCallback), syscall.NewCallbackCDecl(logCallback))
	if dllRet == 0 {
		return fmt.Errorf("%s failed: %v", frontman.PacketFilterStartProc.Name, err)
	}

	return nil
}

// cleanupPlatform for windows is needed to stop the frontman threads and permit the enforcerd app to shut down
func (d *Datapath) cleanupPlatform() {

	dllRet, _, err := frontman.PacketFilterCloseProc.Call()
	if dllRet == 0 {
		zap.L().Error("Failed to close packet proxy", zap.Error(err))
	}

}
