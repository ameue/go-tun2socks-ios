package tun2socks

import (
	"context"
	"log"
	"strings"
	"time"

	vcore "v2ray.com/core"
	vproxyman "v2ray.com/core/app/proxyman"
	vbytespool "v2ray.com/core/common/bytespool"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/shadowsocks"
	"github.com/eycorsican/go-tun2socks/proxy/v2ray"
)

type PacketFlow interface {
	WritePacket(packet []byte)
}

var lwipStack core.LWIPStack

func InputPacket(data []byte) {
	lwipStack.Write(data)
}

func StartV2Ray(packetFlow PacketFlow, configBytes []byte) {
	if packetFlow != nil {
		lwipStack = core.NewLWIPStack()

		core.SetBufferPool(vbytespool.GetPool(core.BufSize))

		v, err := vcore.StartInstance("json", configBytes)
		if err != nil {
			log.Fatal("start V instance failed: %v", err)
		}

		sniffingConfig := &vproxyman.SniffingConfig{
			Enabled:             true,
			DestinationOverride: strings.Split("tls,http", ","),
		}
		ctx := vproxyman.ContextWithSniffingConfig(context.Background(), sniffingConfig)

		vhandler := v2ray.NewHandler(ctx, v)
		core.RegisterTCPConnectionHandler(vhandler)
		core.RegisterUDPConnectionHandler(vhandler)

		core.RegisterOutputFn(func(data []byte) (int, error) {
			packetFlow.WritePacket(data)
			return len(data), nil
		})
	}
}

func StartShadowsocks(packetFlow PacketFlow, proxyHost string, proxyPort int, proxyCipher, proxyPassword string) {
	if packetFlow != nil {
		lwipStack = core.NewLWIPStack()
		core.RegisterTCPConnectionHandler(shadowsocks.NewTCPHandler(core.ParseTCPAddr(proxyHost, uint16(proxyPort)).String(), proxyCipher, proxyPassword))
		core.RegisterUDPConnectionHandler(shadowsocks.NewUDPHandler(core.ParseUDPAddr(proxyHost, uint16(proxyPort)).String(), proxyCipher, proxyPassword, 30*time.Second))
		core.RegisterOutputFn(func(data []byte) (int, error) {
			packetFlow.WritePacket(data)
			return len(data), nil
		})
	}
}
