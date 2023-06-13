package main

import (
	"fmt"
	"log"
	protos "nats/protos/execution"
	protos2 "nats/protos/neworder"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"

	stan "github.com/nats-io/stan.go"
)

var Router *gin.Engine

func main() {

	// --- STAN ---
	clusterID := "test-cluster"
	clientID := "3bd44ce8-20be-4129-870d-e2cc767157a"
	url := "nats://localhost:4222"
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url),
		stan.Pings(1, 30),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))

	if err != nil {
		fmt.Printf("Err occured while connecting stan: %v\n\n", err.Error())
	}

	// Simple Synchronous Publisher
	sc.Publish("USER", []byte("Hello World"))

	sc.QueueSubscribe("EXEC", "my-group", func(m *stan.Msg) {
		fmt.Printf("\n\nyaya\n\n")
		fmt.Printf("Data: %v\n\n", m.Data)
	}, stan.StartAtSequence(0))

	// Simple Async Subscriber
	_, err = sc.Subscribe("EXEC", func(m *stan.Msg) {
		fmt.Printf("Data: %v\n\n", m.Data)
		// fmt.Printf("Data: %#v\n\n", m.Data)
		// fmt.Printf("Data: %+v\n\n", m.Data)
		// fmt.Printf("Data: %q\n\n", m.Data)
		// fmt.Printf("Data: %+q\n\n", m.Data)
		fmt.Printf("Data: %#q\n\n", m.Data)
		// fmt.Printf("Data: %x\n\n", m.Data)
		fmt.Printf("Data: %+x\n\n", m.Data)
		fmt.Printf("Data: %+v\n\n", string(m.Data[:]))
		// fmt.Printf("Data: %+v\n\n", hex.EncodeToString(m.Data))

		// val1, size := utf8.DecodeRune(m.Data)
		// fmt.Printf("DecodeRune: %c, %v\n\n", val1, size)

		// val2, size := utf8.DecodeLastRune(m.Data)
		// fmt.Printf("DecodeLastRune: %c, %v\n\n", val2, size)

		// var b bytes.Buffer
		// n, er := b.Read(m.Data)
		// fmt.Printf("n: %v, er: %v\n\n", n, er)

		varr := &protos.Execution{}
		err = proto.Unmarshal(m.Data, varr)
		if err != nil {
			fmt.Printf("Err unmarshalling!: %v\n\n", err.Error())
		} else {
			fmt.Printf("Received a message: %+v\n", varr)
		}

		varr2 := &protos2.NewOrderRequest{}
		err = proto.Unmarshal(m.Data, varr2)
		if err != nil {
			fmt.Printf("Err unmarshalling!: %v\n\n", err.Error())
		} else {
			fmt.Printf("Received a message: %+v\n", varr2)
		}
	})

	if err != nil {
		fmt.Printf("Err while subscribing in stan: %v\n", err.Error())
	}

	// --- NATS ---

	// Connect to a server
	/*
		nc, err := nats.Connect("nats://localhost:4222")
		if err != nil {
			fmt.Printf("Err: %v\n\n", err.Error())
		}

		nc.Publish("USER", []byte("Hello World"))
		nc.Publish("EXEC", []byte("yay!"))
		nc.Publish("MARKET", []byte("non!!"))

		sub2, err := nc.Subscribe("USER", func(m *nats.Msg) {
			if m == nil {
				fmt.Println("Message is nil in USER!")
			} else {
				fmt.Printf("Received a message in USER: %s\n", string(m.Data))
			}
		})
		if err != nil {
			fmt.Printf("Error subscribing to USER: %v\n", err.Error())
		}
		fmt.Printf("Sub user: %+v\n\n\n", sub2)

		nc.Subscribe("EXEC", func(m *nats.Msg) {
			if m == nil {
				fmt.Println("Message is nil in EXEC!")
			} else {
				fmt.Printf("Received a message in EXEC: %s\n", string(m.Data))
			}
		})

		// nc.Subscribe("MARKET", func(m *nats.Msg) {
		// 	if m == nil {
		// 		fmt.Println("Message is nil in MARKET!")
		// 	} else {
		// 		fmt.Printf("Received a message in MARKET: %s\n", string(m.Data))
		// 	}
		// })

		// nc.Subscribe("CHANGES", func(m *nats.Msg) {
		// 	if m == nil {
		// 		fmt.Println("Message is nil in CHANGES!")
		// 	} else {
		// 		fmt.Printf("Received a message in CHANGES: %s\n", string(m.Data))
		// 	}
		// })

		// nc.Subscribe(">", func(m *nats.Msg) {
		// 	if m == nil {
		// 		fmt.Println("Message is nil in wildcard!")
		// 	} else {
		// 		fmt.Printf("Message: %v Subject: %v\n\n", string(m.Data), m.Subject)
		// 	}
		// })
	*/

	// queue group subscribe
	// nc.QueueSubscribe("EXEC", "my_own_queue_xyz", func(m *nats.Msg) {
	// 	fmt.Println("From queue group EXEC: ", string(m.Data))
	// })

	// nc.QueueSubscribe("USER", "my_own_queue_xyz2", func(m *nats.Msg) {
	// 	fmt.Println("From queue group USER: ", string(m.Data))
	// })

	// nc.QueueSubscribe("MARKET", "my_own_queue_xyz3", func(m *nats.Msg) {
	// 	fmt.Println("From queue group MARKET: ", string(m.Data))
	// })

	// channel subcribe
	// ch := make(chan *nats.Msg, 64)
	// sub, err := nc.ChanSubscribe("EXEC", ch)
	// if err != nil {
	// 	fmt.Printf("Error occurred: %v", err.Error())
	// }
	// msg := <-ch

	// fmt.Printf("SUB in channel: %v", sub)
	// fmt.Printf("MSG in channel: %v", string(msg.Data))

	// sync
	// sub, err := nc.SubscribeSync("EXEC")
	// if err != nil {
	// 	fmt.Printf("Err in subsync: %v\n\n", err.Error())
	// }
	// m, err := sub.NextMsg(10 * time.Minute)
	// if err != nil {
	// 	fmt.Printf("Err in nexmsg: %v\n\n", err.Error())
	// }
	// fmt.Printf("Msg: %v\n", m.Data)

	// jetstream
	// js, err := nc.JetStream(nats.PublishAsyncMaxPending(1024))
	// if err != nil {
	// 	fmt.Printf("Err in jet stream: %v\n", err.Error())
	// }

	// js.Subscribe("EXEC", func(m *nats.Msg) {
	// 	fmt.Printf("Received a JetStream message: %s\n", string(m.Data))
	// })

	// js.Subscribe(">", func(m *nats.Msg) {
	// 	fmt.Printf("Jetstream message: %v sub: %v\n ", string(m.Data), m.Subject)
	// })

	gin.SetMode(gin.DebugMode)
	Router = gin.New()
	Router.Run(":1212")
}
