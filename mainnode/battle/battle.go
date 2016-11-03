package battle

import (
	"time"
	"strconv"
	"sgsld/msgHandler"
	"math/rand"
	"fmt"
)

//游戏状态信息
type Battle struct {
	uids [PLAYER_NUM]int
	state int           		//控制游戏状态-未开始,等待玩家准备,正在战斗,游戏结束
	round int           		//回合数,双方玩家下达指令至下一次准备下达之前,为一回合
	command_state [PLAYER_NUM]int	//玩家准备状态(是否已下达指令)
	game_time int64     		//场景时间
	event_queue Queue       	//事件队列
	act_stack Stack     		//行动栈
	delay int64         		//延时
	waiting_time int64     	 	//等待时间倒计时
	//主要数据结构
	library [PLAYER_NUM][CARD_NUM]Command_Card   	//双方牌库
	generals [GENERAL_NUM]General       		//武将序列
	general_state [GENERAL_NUM]General_State	//武将状态信息
	pos_arr [POS_NUM]int            		//位置信息
	//辅助数据记录
	winner bool                 //胜利方, 为true标识first胜利,为false则last胜利
	//其他数据记录
	my_rand *rand.Rand
}

/**-----战斗准备时间部分-----**/
func (bt *Battle)Init(user1 User_info, user2 User_info){
	//toDO 检查user_info格式 包括武将是否符合条件等
	//设定先手和后手
	bt.uids[FIRST] = user1.uid
	bt.uids[LAST] = user2.uid
	//设定战斗状态
	bt.state = STATE_NOT_BEGIN
	//初始化武将序列
	for i:=0;i<=3;i++ {
		bt.generals[i] = user1.Generals[i]
	}
	for i:=4;i<=7;i++ {
		bt.generals[i] = user2.Generals[i-4]
	}
	//武将位置
	bt.pos_arr = const_pos
	//初始化武将状态
	for k,v := range bt.generals{
		if(v.hp_max > 0){
			bt.general_state[k].isAlive = true
			bt.general_state[k].hp = v.hp_max
		}else{
			bt.general_state[k].isAlive = false
		}
	}
	//初始化指令库
	bt.my_rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	bt.library[FIRST] = bt.shuffleCard(user1.Cards)
	bt.library[LAST] = bt.shuffleCard(user2.Cards)
}

func (bt *Battle)Begin(){
	//开始游戏
	bt.game_time = 0
	bt.round = 0
	bt.prepareStart()      //开始回合准备阶段
	bt.update()
}

func (bt *Battle)SetPrepareState(uid int, card_order int){
	user_order := bt.getUserOrder(uid)
	bt.command_state[user_order] = card_order
	bt.simOutIdx(uid," 玩家已选择指令")
	if(bt.isAllReady()){
		bt.prepareEnd()
	}
}

//每帧调用
func (bt *Battle)update(){
	bt.game_time += FRAME
	//判断当前状态, 是否是在等待玩家指令
	if(bt.state == STATE_WAITING){
		bt.waiting_time -= FRAME
		if(bt.waiting_time <=0){
			//玩家选择指令超时
			bt.simOut("选择指令超时!")
			if(bt.command_state[FIRST] == NIL_INT){
				bt.command_state[FIRST] = bt.my_rand.Intn(4)
				bt.sendToFirst("随机选择第"+strconv.Itoa(bt.command_state[FIRST])+"张指令牌")
			}
			if(bt.command_state[LAST] == NIL_INT){
				bt.command_state[LAST] = bt.my_rand.Intn(4)
				bt.sendToLast("随机选择第"+strconv.Itoa(bt.command_state[LAST])+"张指令牌")
			}
			bt.prepareEnd()
		}
	}else{
		//开始战斗,执行战斗事件
		//delay标识该动作需要的时间
		bt.delay -= FRAME
		if(bt.delay < 0){
			bt.delay = 0
		}
		//如果没有delay, 则执行下一个时间
		for bt.delay==0{
			//如果游戏已经结束
			if(bt.isEnd()){
				bt.gameEnd()
				return
			}
			//否则执行下一个时间
			if(!bt.act_stack.Empty()){
				//非空执行栈内部分
				event, err := bt.act_stack.Pop()
				if(err != nil){
					//异常,这种情况不应该发生(因为之前已经判断stack不为空)
				}else {
					bt.processEvent(event)
				}
			}else{
				event, err := bt.event_queue.Dequeue()
				if(err != nil){
					break   //没有事件执行,等待事件中
				}else {
					bt.processEvent(event)
				}
			}
		}
	}
	t := time.NewTimer(time.Duration(FRAME))
	<-t.C
	bt.update()

}

/**-----回合准备时间部分-----**/
func (bt *Battle)prepareStart(){
	bt.command_state[FIRST] = NIL_INT
	bt.command_state[LAST] = NIL_INT
	bt.state = STATE_WAITING
	bt.waiting_time = int64(MAX_WAIT_TIME)
	//发送可选指令牌数据
	bt.simOutCards(true,bt.library[FIRST])
	bt.simOutCards(false,bt.library[LAST])
	bt.simOut("等待玩家指令中")
}

func (bt *Battle)prepareEnd(){
	bt.state = STATE_FIGHTING
	var info []string
	info = append(info, strconv.Itoa(bt.command_state[FIRST]))
	info = append(info, strconv.Itoa(bt.command_state[LAST]))
	bt.addEvent(0, EVENT_EXECUTE_COMMAND, info)
	bt.command_state[FIRST] = NIL_INT
	bt.command_state[LAST] = NIL_INT
	bt.addEvent(0, EVENT_ROUND_BEGIN, nil)
}

func (bt *Battle)isAllReady() bool{
	return bt.command_state[FIRST] != NIL_INT && bt.command_state[LAST] != NIL_INT
}

/**-----战斗结束,结算部分----**/
func (bt *Battle)isEnd() bool{
	if( !bt.hasAlive(false) ){
		bt.winner = true
		return true
	}
	if( !bt.hasAlive(true) ){
		bt.winner = false
		return true
	}
	return false
}

func (bt *Battle)hasAlive(isFirst bool) bool{
	var idx_arr []int
	if(isFirst){
		idx_arr = first_idx
	}else {
		idx_arr = last_idx
	}
	for _,idx := range idx_arr{
		if(bt.general_state[idx].isAlive){
			return true
		}
	}
	return false
}

func (bt *Battle)gameEnd(){
	if(bt.winner){
		bt.simOut("Game End! first win")
	}else{
		bt.simOut("Game End! last win")
	}
}

/**-------事件处理部分------**/
func (bt *Battle)processEvent(eve Event){
	bt.delay += bt.getDelay(eve)
	id := eve.id
	info := eve.info
	//在这里,处理事件格式,即所有事件的格式都在这里定义
	switch eve.event_Type {
	case EVENT_ROUND_BEGIN:
		bt.processRoundBegin(id)
	case EVENT_ROUND_END:
		bt.processRoundEnd(id)
	case EVENT_ACT_BEGIN:
		bt.processActBegin(id)
	case EVENT_ACT_BEFORE:
		bt.processActBefore(id)
	case EVENT_ACT_ATTACK:
		bt.processActAttack(id)
	case EVENT_ACT_AFTER:
		bt.processActAfter(id)
	case EVENT_ACT_END:
		bt.processActEnd(id)
	case EVENT_ATK_BEGIN:
		//info的第一个参数为目标
		to,err := strconv.Atoi(info[0])
		if err != nil{
			//todo 目标不合法错误处理
		}else{
			bt.processAtkBegin(id,to)
		}
	case EVENT_ATK_DMG:
		//info的第一个参数为目标
		to,err := strconv.Atoi(info[0])
		if err != nil{
			//todo 目标不合法错误处理
		}else{
			bt.processAtkDmg(id,to)
		}
	case EVENT_ATK_END:
		//info的第一个参数为目标
		to,err := strconv.Atoi(info[0])
		if err != nil{
			//todo 目标不合法错误处理
		}else{
			bt.processAtkEnd(id,to)
		}
	case EVENT_GENERAL_DIE:
		bt.processGeneralDie(id)
	case EVENT_GENERAL_EXIT:
		bt.processGeneralExit(id)
	case EVENT_EXECUTE_COMMAND:
		card1_order,err1 := strconv.Atoi(info[0])
		card2_order,err2 := strconv.Atoi(info[1])
		if err1 != nil || err2 != nil{
			fmt.Println("字符串转换成整数失败")
		}
		bt.processExecuteCommand(card1_order,card2_order)
	}
}
func (bt *Battle)processExecuteCommand(card1_order int, card2_order int){
	card1_name := bt.library[FIRST][card1_order].name
	bt.simOut("玩家一使用了卡牌-"+card1_name)
	card2_name := bt.library[LAST][card2_order].name
	bt.simOut("玩家二使用了卡牌-"+card2_name)
	//将其放在最后
	bt.chgOrderAfterUseCard(true,card1_order)
	bt.chgOrderAfterUseCard(false,card2_order)
}

func (bt *Battle)processRoundBegin(id int){
	bt.round++
	bt.simOutIdx(bt.round,"回合开始")
	act_order := []int{0,5,1,6,2,7,3,4}
	for _, pos := range act_order{
		bt.addEvent(pos,EVENT_ACT_BEGIN,nil)
	}
	bt.addEvent(0, EVENT_ROUND_END,nil)
}

func (bt *Battle)processRoundEnd(id int){
	bt.simOutIdx(bt.round,"回合结束")
	if(bt.round >= MAX_ROUND){
		bt.gameEnd()
		return
	}
	//开始等待玩家指令
	bt.prepareStart()
}

func (bt *Battle)processActBegin(pos int){
	idx := bt.pos_arr[pos]
	if(bt.checkActState(idx)){
		bt.simOutIdx(idx,"行动开始")
		bt.pushStack(idx,EVENT_ACT_END,nil)
		bt.pushStack(idx,EVENT_ACT_AFTER,nil)
		bt.pushStack(idx,EVENT_ACT_ATTACK,nil)
		bt.pushStack(idx,EVENT_ACT_BEFORE,nil)
	}
}

func (bt *Battle)processActBefore(id int){
	//行动前触发技能
}

func (bt *Battle)processActAttack(id int){
	//行动时, 普攻
	target := bt.getNormalTaget(id)
	if(target == NIL_INT){
		//取不到合适的目标
		return
	}
	target_info := []string{strconv.Itoa(target)}
	bt.pushStack(id, EVENT_ATK_BEGIN, target_info)
}

func (bt *Battle)processActAfter(idx int){
	//行动后
}

func (bt *Battle)processActEnd(idx int){
	//行动结束
	bt.simOutIdx(idx,"结束")
}

func (bt *Battle)processAtkBegin(id int, to int){
	//攻击前触发天香类技能
	//不触发则触发普通攻击动作
	var info []string
	info = append(info,strconv.Itoa(to))
	bt.pushStack(id, EVENT_ATK_END, info)
	bt.pushStack(id, EVENT_ATK_DMG, info)
}

func (bt *Battle)processAtkDmg(from int, to int){
	damage := bt.calDamage(from,to,DMGTYPE_NORMALATK)
	bt.doDamage(from,to,DMGTYPE_NORMALATK,damage)
}

func (bt *Battle)processAtkEnd(id int, to int){
	//攻击后触发技能
}

func (bt *Battle)processGeneralDie(id int){
	//如果有复活技能...
	bt.pushStack(id,EVENT_GENERAL_EXIT,nil)
}

func (bt *Battle)processGeneralExit(id int){
	//下场操作, 这里没有其他武将商场
	bt.general_state[id].isAlive = false
	bt.simOutIdx(id, "died")
}

/**-------增加事件------**/
func (bt *Battle)pushStack(id int, e_type int, extra_info []string){
	eve := Event{id:id,event_Type:e_type,info:extra_info}
	bt.act_stack.Push(eve)
}

func (bt *Battle)addEvent(id int, e_type int, info []string){
	event := Event{id:id,event_Type:e_type,info:info}
	bt.event_queue.Enqueue(event)
}

/**-------战斗计算部分----**/
func (bt *Battle)checkActState(idx int) bool{
	return idx != NIL_INT && bt.general_state[idx].isAlive
}

func (bt *Battle)checkValidTarget(idx int) bool{
	return idx != NIL_INT && bt.general_state[idx].isAlive
}

func (bt *Battle)getNormalTaget(pos int) int{
	var target int
	for _,val := range target_arr{
		target = val[pos]
		if( bt.checkValidTarget(bt.pos_arr[target]) ){
			return bt.pos_arr[target]
		}
	}
	//取不到合适的目标
	return NIL_INT
}

func (bt *Battle)calDamage(from int, to int, dmg_type int) int{
	atk := bt.generals[from].attack
	def := bt.generals[from].defense
	damage := atk - def
	return damage
}

func (bt *Battle)doDamage(from int, to int, dmg_type int, damage int){
	bt.general_state[to].hp -= damage
	if(bt.general_state[to].hp <= 0){
		bt.pushStack(to,EVENT_GENERAL_DIE,nil)
	}
	str := strconv.Itoa(from)+"对"+strconv.Itoa(to)+"造成了"+strconv.Itoa(damage)+"点伤害"
	bt.simOut(str)
}

/**--------指令牌库部分-----**/
func (bt *Battle) shuffleCard(in [CARD_NUM]Command_Card ) [CARD_NUM]Command_Card  {
	for i := CARD_NUM - 1; i > 0; i-- {
		r := bt.my_rand.Intn(CARD_NUM)
		in[r], in[i] = in[i], in[r]
	}
	return in
}

func (bt *Battle) chgOrderAfterUseCard(isFirst bool, card_order int){
	last_index := CARD_NUM - 1	//最后一张卡的idx
	if(isFirst){
		tmp_card := bt.library[FIRST][card_order]
		for i:=card_order;i<last_index;i++{
			bt.library[FIRST][i] = bt.library[FIRST][i+1]
		}
		bt.library[FIRST][last_index] = tmp_card
	}else{
		tmp_card := bt.library[LAST][card_order]
		for i:=card_order;i<last_index;i++{
			bt.library[LAST][i] = bt.library[LAST][i+1]
		}
		bt.library[LAST][last_index] = tmp_card
	}
}
/**---------辅助函数-------**/
func (bt *Battle)getUserOrder(uid int) int{
	switch {
	case bt.uids[FIRST] == uid:return FIRST
	case bt.uids[LAST] == uid:return LAST
	}
	return NIL_INT
}

/**-------模拟器输出------**/
func (bt *Battle)sendToFirst(str string){
	str = str + "\r\n"
	msgHandler.Send(bt.uids[FIRST],str)
}

func (bt *Battle)sendToLast(str string){
	str = str + "\r\n"
	msgHandler.Send(bt.uids[LAST],str)
}

func (bt *Battle)simOut(str string){
	bt.sendToFirst(str)
	bt.sendToLast(str)
}

func (bt *Battle)simOutIdx(idx int, str string){
	bt.simOut(strconv.Itoa(idx)+" "+str)
}

func (bt *Battle)simOutCards(isFirst bool, cards [CARD_NUM]Command_Card){
	str := "可选的卡牌为: "
	for i := 0;i<AVAILABLE_CARD_NUM;i++{
		str = str + strconv.Itoa(i) + ". "
		str = str + cards[i].name + " "
	}
	if(isFirst){
		bt.sendToFirst(str)
	}else{
		bt.sendToLast(str)
	}
}
/**-----获取事件的delay---**/
func (bt *Battle)getDelay(eve Event) int64{
	default_delay := int64(time.Second)
	switch eve.event_Type {
	case EVENT_ATK_BEGIN,EVENT_ATK_DMG,EVENT_ATK_END,EVENT_GENERAL_EXIT:
		return default_delay
	}
	return 0
}
