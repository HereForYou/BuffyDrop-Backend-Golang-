package models

type TaskList struct {
	Id     string `bson:"id"`
	Title  string `bson:"title"`
	Image  string `bson:"image"`
	Profit int32  `bson:"profit"`
	Link   string `bson:"link"`
}

type Admin struct {
	Id       int    `bson:"id"`
	Username string `bson:"username"`
}

type Setting struct {
	TaskList      []TaskList `bson:"taskList"`
	Admin         []Admin    `bson:"admin"`
	InviteRevenue float64    `bson:"inviteRevenue"`
	DailyRevenue  float64    `bson:"dailyRevenue"`
}
