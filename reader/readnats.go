package main

import (
	"github.com/nats-io/nats"
	"fmt"
	"runtime"
	"encoding/json"
	"flag"
	"log"
)

type DhcpMessage struct {
	MacAddress  string `json:"macAddress"`
	MessageType string `json:"messageType"`
	Network     string `json:"network"`
	Address     string `json:"address"`
}

func main() {

	natsUrl := flag.String("natsurl", "", "nats server url as in nats://host:port")
	flag.Parse()

	if (*natsUrl == "") {
		flag.Usage()
		log.Fatal("-natsurl required")
	}

	con, err := nats.Connect(*natsUrl)

	if err != nil {
		panic(err)
	}

	jcon, err := nats.NewEncodedConn(con, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}

	jcon.Subscribe("dhcp", func(msg *DhcpMessage) {
		str, err := json.MarshalIndent(msg, "", " ")
		if err != nil {
			fmt.Println("error decoding", err)
			return
		}
		fmt.Println(string(str))
	})

	runtime.Goexit()

}
