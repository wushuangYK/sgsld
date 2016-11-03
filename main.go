package main

import (
	"net"
	"encoding/json"
	"strings"
	"sgsld/msgHandler"
	"sgsld/config"
	"fmt"
	"sgsld/mainnode/battle"
	"log"
)

func main(){
	// Listen on TCP port 2000 on all interfaces.
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	msgHandler.Instance()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleCon(conn)
	}
}

const INIT = 0

//处理连接
func handleCon(c net.Conn){
	b := make([]byte, config.MAX_READ_LENGTH)
	var err error
	for {
		c.Read(b)
		var msg msgHandler.Msg
		dec := json.NewDecoder(strings.NewReader(string(b)))
		if err = dec.Decode(&msg); err != nil {
			//todo 错误处理
			fmt.Println("failed to decode")
		}
		if(msg.Do == INIT){
			msgHandler.Set(msg.Uid,c)
		}
		if err = battle.HandleMsg(msg.Do,msg.Uid,msg.Data); err != nil {
			//todo 错误处理
			break
		}
	}
	c.Close()
}