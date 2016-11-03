package battle

func GetTeamByUid(uid int) User_info{
	var user_info User_info = User_info{
		uid:uid,
		Generals:[MAX_USER_GENERAL_NUM]General{
			{
				general_type:1,
				attack:100,
				defense:10,
				hp_max:1500,
			},
		},
		Cards:[CARD_NUM]Command_Card{
			{
				name:"杀",
				effect:nil,
			},
			{
				name:"闪",
				effect:nil,
			},
			{
				name:"酒",
				effect:nil,
			},
			{
				name:"桃",
				effect:nil,
			},
			{
				name:"顺手牵羊",
				effect:nil,
			},
			{
				name:"无中生有",
				effect:nil,
			},
			{
				name:"过河拆桥",
				effect:nil,
			},
			{
				name:"无懈可击",
				effect:nil,
			},
		},
	}
	return user_info
}
