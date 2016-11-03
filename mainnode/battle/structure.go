package battle

type Event struct{
	id int
	event_Type int
	info []string
}
type User_info struct {
	uid int
	Generals [MAX_USER_GENERAL_NUM]General
	Cards [CARD_NUM]Command_Card
}
type General struct {
	general_type int
	attack int
	defense int
	hp_max int
}
type Skill struct {
	skill_type int

}
type Command_Card struct {
	name string
	effect []string
}
//武将状态信息
type General_State struct {
	isAlive bool
	hp int
}