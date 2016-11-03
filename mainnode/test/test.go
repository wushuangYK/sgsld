package test

import (
	"sgsld/mainnode/battle"
	"fmt"
	"sync"
	"math/rand"
	"time"
)

var wg sync.WaitGroup

func TestBattle(){
	user_info1 := battle.GetTeamByUid(1)
	user_info2 := battle.GetTeamByUid(2)
	bt := new(battle.Battle)
	bt.Init(user_info1,user_info2)
}

func testGo(user_info1 battle.User_info, user_info2 battle.User_info){
	bt := new(battle.Battle)
	bt.Init(user_info1,user_info2)
	bt.Begin()
	defer wg.Done()
}

func TestSwitch(){
	i := 1;
	switch {
	case i==1:fmt.Println("1")
	case i<2:fmt.Println("2")
	}
}

func TestReturn(){
	i := 1
	if(i==1){
		fmt.Println("1")
		return
	}
	fmt.Println("2")
}

func TestSlice(){
	var s []int
	fmt.Println(s[1])
}

func TestRand(){
	var my_rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=1;i<20;i++{
		fmt.Println(my_rand.Intn(5))
	}
}