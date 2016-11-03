package battle

import (
	"time"
)

const PLAYER_NUM = 2				//player_number
const MAX_WAIT_TIME = time.Second*15		//最大等待时间
const FRAME = int64(time.Millisecond*500)	//帧长
const MAX_USER_GENERAL_NUM = 4			//最大武将数目
const GENERAL_NUM = 2*MAX_USER_GENERAL_NUM	//总武将数目
const POS_NUM = 8			//战场上的位置数目,双方各4人
const MAX_ROUND = 20			//最大回合数
const NIL_INT = -1			//标识不需要武将id
const MAX_SPIRIT = 6			//最大武魂点数
const COST_SPIRIT = 4			//武魂技能消耗的武魂点数
const MAX_ANGER	= 100			//怒气最大值,当怒气达到最大值时,消耗怒气
const CARD_NUM = 8			//卡牌数量
const AVAILABLE_CARD_NUM = 4		//可用的卡牌数量
var const_pos = [POS_NUM]int{3,0,1,2,4,5,6,7}	//确定位置-配置每个武将所在的位置
var target_arr0 = []int{5,4,7,4,3,0,3,2}//目标选择序列
var target_arr1 = []int{4,6,4,6,1,3,1,3}//目标选择序列,当第一目标死亡时
var target_arr2 = []int{6,5,6,5,0,1,0,1}//目标选择序列,当第一,第二目标都死亡时
var target_arr3 = []int{7,7,5,7,2,2,2,0}//目标选择序列,当第一,第二目标都死亡时
var target_arr = [][]int{target_arr0,target_arr1,target_arr2,target_arr3}
var first_idx = []int{0,1,2,3}		//先手方idx
var last_idx = []int{4,5,6,7}		//后手方idx
const (
	STATE_NOT_BEGIN	= iota	//战斗未开始
	STATE_WAITING		//等待玩家指令中
	STATE_FIGHTING		//战斗中
	STATE_END		//战斗结束
)
const (
	FIRST	= iota		//先手
	LAST			//后手
)
const (
	GENERAL_LORD	= iota	//主将
	GENERAL_FOLLOWER	//随从
)
const (
	NOT_PREPARE = iota	//未准备
	READY			//准备就绪
)
//这里对应于流程的控制
const (
	EVENT_ROUND_BEGIN	= iota	//R标识round
	EVENT_ACT_BEGIN		//G标识武将
	EVENT_ACT_BEFORE
	EVENT_ACT_ATTACK
	EVENT_ACT_AFTER
	EVENT_ACT_END
	EVENT_ROUND_END
	EVENT_ATK_BEGIN		//攻击出手动作
	EVENT_ATK_DMG		//攻击动作--造成伤害
	EVENT_ATK_END		//攻击结束动作
	EVENT_GENERAL_DIE	//G标识武将,武将死亡
	EVENT_GENERAL_EXIT	//武将下场
	EVENT_EXECUTE_COMMAND	//执行卡牌效果
)
const (
	DMGTYPE_NORMALATK = iota
	DMGTYPE_SKILL
)
const (
	SKILL_TYPE_NORMAL = iota	//普通技能
	SKILL_TYPE_PASSIVE		//被动技能
	SKILL_TYPE_SPIRIT		//武魂技能
	SKILL_TYPE_LORD			//主公技能
)
