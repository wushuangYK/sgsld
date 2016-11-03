package battle

import (
	"strconv"
	"errors"
)

const (
	BATTLE_MESSAGE_PREPARE	= 11
	BATTLE_MESSAGE_COMMAND = 12
)

var battleIns map[int]*Battle = make(map[int]*Battle)
var roomLen int = 0
var userRoom map[int]int = make(map[int]int)

func HandleMsg(do int, uid int, info string) error{
	switch do {
	case BATTLE_MESSAGE_PREPARE:
		enemy_id,err := strconv.Atoi(info)
		if(err != nil){
			return errors.New("invalid variables")
		}
		roomId1 := getRoomId(uid)
		roomId2 := getRoomId(enemy_id)
		switch {
		case roomId1 == NIL_INT && roomId2 != NIL_INT:userRoom[uid] = userRoom[enemy_id]
		case roomId1 != NIL_INT && roomId2 == NIL_INT:userRoom[enemy_id] = userRoom[uid]
		case roomId1 == NIL_INT && roomId2 == NIL_INT:
			roomId := newBattle()
			userRoom[uid] = roomId
		}
		roomId1 = getRoomId(uid)
		roomId2 = getRoomId(enemy_id)
		if( roomId1 == roomId2 && roomId1 != NIL_INT ){
			user_info1 := GetTeamByUid(uid)
			user_info2 := GetTeamByUid(enemy_id)
			bt := battleIns[roomId1]
			bt.Init(user_info1, user_info2)
			go bt.Begin()
		}
	case BATTLE_MESSAGE_COMMAND:
		roomId := getRoomId(uid)
		if(roomId != NIL_INT){
			bt := battleIns[roomId]
			//发送指令
			card_order,error := strconv.Atoi(info)
			if error != nil{
				return errors.New("invalid variables")
			}else{
				bt.SetPrepareState(uid,card_order)
			}

		}else{
			return errors.New("invalid variables")
		}
	}
	return nil
}

func newBattle() int{
	roomLen++
	battleIns[roomLen-1] = new(Battle)
	return roomLen-1
}

func getRoomId(uid int) int{
	if roomId, ok := userRoom[uid]; ok {
		return roomId
	} else {
		return NIL_INT
	}
}