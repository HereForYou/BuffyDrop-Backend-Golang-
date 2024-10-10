package models

type TaskList struct {
	Id     string `bson:"id" json:"id"`
	Title  string `bson:"title" json:"title"`
	Image  string `bson:"image" json:"image"`
	Profit int32  `bson:"profit" json:"profit"`
	Link   string `bson:"link" json:"link"`
}

type Admin struct {
	Id       int    `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
}

type Setting struct {
	TaskList      []TaskList `bson:"taskList" json:"taskList"`
	Admin         []Admin    `bson:"admin" json:"admin"`
	InviteRevenue float64    `bson:"inviteRevenue" json:"inviteRevenue"`
	DailyRevenue  float64    `bson:"dailyRevenue" json:"dailyRevenue"`
}
