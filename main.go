package main

import (
	"github.com/nats-io/nats"
	"flag"
	"log"
	"net"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/krolaw/dhcp4"
	"strings"
)

func main() {

	natsUrl := flag.String("natsurl", "", "nats server url as in nats://host:port")
	flag.Parse()

	if (*natsUrl == "") {
		flag.Usage()
		log.Fatal("missing natsurl paramter")
	}

	con, err := nats.Connect(*natsUrl)
	if err != nil {
		panic(err)
	}
	defer con.Close()

	addr, err := net.ResolveUDPAddr("udp4", "224.0.0.1:67")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1500)

	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Println(err)
			continue
		}

		if n < 240 {
			// Packet too small to be DHCP
			fmt.Println("packet too small")
			continue
		}

		req := dhcp4.Packet(buffer[:n])
		if req.HLen() > 16 {
			// Invalid size
			fmt.Println("invalid size", req.HLen())
			continue
		}

		options := req.ParseOptions()
		var reqType dhcp4.MessageType
		if t := options[dhcp4.OptionDHCPMessageType]; len(t) != 1 {
			continue
		} else {
			reqType = dhcp4.MessageType(t[0])
			if reqType < dhcp4.Discover || reqType > dhcp4.Inform {
				continue
			}

			msg := make(map[string]string)
			msg["network"] = addr.Network()
			msg["address"] = addr.String()
			msg["messageType"] = reqType.String()
			msg["macAddress"] = strings.ToLower(req.CHAddr().String())

			jsonbuf := bytes.NewBuffer(nil)
			err = json.NewEncoder(jsonbuf).Encode(&msg)
			if err != nil {
				log.Println("error encoding", err)
				continue
			}

			natsPayload := jsonbuf.Bytes()

			con.Publish("dhcp", natsPayload)

		}
	}

}
