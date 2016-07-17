package main

import (
	"github.com/nats-io/nats"
	"fmt"
	"runtime"
	"encoding/json"
	"bytes"
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
		log.Fatal("missing natsurl paramter")
	}

	con, err := nats.Connect(*natsUrl)

	if err != nil {
		panic(err)
	}

	con.Subscribe("dhcp", func(msg *nats.Msg) {

		message := &DhcpMessage{}
		json.NewDecoder(bytes.NewBuffer(msg.Data)).Decode(message)

		str, err := json.MarshalIndent(message, "", " ")
		if err != nil {
			fmt.Println("error decoding", err)
			return
		}
		fmt.Println(string(str))
	})

	runtime.Goexit()

}
