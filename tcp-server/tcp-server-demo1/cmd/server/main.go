package main

import (
	"fmt"
	"net"

	"github.com/bigwhite/tcp-server-demo1/frame"
	"github.com/bigwhite/tcp-server-demo1/packet"
)

func handlePacket(framePayload []byte) (ackFramePayload []byte, err error) {
	var p packet.Packet
	p, err = packet.Decode(framePayload)
	if err != nil {
		fmt.Println("handleConn: packet decode error:", err)
		return
	}

	switch p.(type) {
	case *packet.Submit:
		submit := p.(*packet.Submit)
		fmt.Printf("recv submit: id = %s, payload=%s\n", submit.ID, string(submit.Payload))
		submitAck := &packet.SubmitAck{
			ID:     submit.ID,
			Result: 0, // 0代表成功
		}
		ackFramePayload, err = packet.Encode(submitAck)
		if err != nil {
			fmt.Println("handleConn: packet encode error:", err)
			return nil, err
		}
		return ackFramePayload, nil
	default:
		return nil, fmt.Errorf("unknown packet type")
	}
}

/*
为什么在协程中先处理一次panic
协程隔离:

	每个 handleConn 函数都是在一个单独的协程中执行的。
	如果在 main 函数中统一处理所有协程的 panic,一旦某个协程 panic,可能会影响到其他协程,导致整个服务器崩溃。
	但是在 handleConn 内部处理 panic,可以确保 panic 只影响当前的连接处理,不会影响到其他连接的处理。

错误定位:

	在 handleConn 内部处理 panic,可以更精确地定位问题所在,因为我们知道 panic 发生在哪个连接的处理中。
	而如果只在 main 函数中处理 panic,很难确定 panic 是由哪个连接处理引起的。

连接状态保持:

	在 handleConn 内部处理 panic 并恢复,可以确保连接能够继续处理下一个请求,保持服务的连续性。
	如果在 main 函数中处理 panic,一旦发生 panic,整个连接就会被关闭,客户端需要重新建立连接。

所以,在 handleConn 函数中处理 panic 是一种更好的做法,可以确保服务器的健壮性和可靠性。这样即使某个连接处理出现问题,其他连接也不会受到影响。
*/
func handleConn(c net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("handleConn: recovered from panic:", r)
		}
		c.Close()
	}()
	frameCodec := frame.NewMyFrameCodec()

	for {
		// read from the connection

		// decode the frame to get the payload
		// the payload is undecoded packet
		framePayload, err := frameCodec.Decode(c)
		if err != nil {
			fmt.Println("handleConn: frame decode error:", err)
			return
		}

		// do something with the packet
		ackFramePayload, err := handlePacket(framePayload)
		if err != nil {
			fmt.Println("handleConn: handle packet error:", err)
			return
		}

		// write ack frame to the connection
		err = frameCodec.Encode(c, ackFramePayload)
		if err != nil {
			fmt.Println("handleConn: frame encode error:", err)
			return
		}
	}
}

func main() {
	// panic处理
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("main: recovered from panic:", r)
		}
	}()

	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	fmt.Println("server start ok(on *.8888)")

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		// start a new goroutine to handle
		// the new connection.
		go handleConn(c)
	}
}
