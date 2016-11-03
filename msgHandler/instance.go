package msgHandler

import (
	"net"
)

type Msg struct {
	Do int
	Uid int
	Data string
}

var _conn map[int]net.Conn = make(map[int]net.Conn)

func Instance() (map[int]net.Conn) {
	if _conn == nil {
		_conn = make(map[int]net.Conn)
	}
	return _conn
}

func Get(id int) (c net.Conn){
	if(_conn[id] == nil){
		return nil
	}
	return _conn[id]
}

func Set(id int, c net.Conn){
	if _conn == nil {
		_conn = make(map[int]net.Conn)
	}
	_conn[id] = c
}

func Isset(id int) (flag bool){
	return _conn[id] == nil
}

func Send(id int, data string){
	_conn[id].Write([]byte(data))
}