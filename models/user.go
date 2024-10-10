package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Friend struct {
	Id      string  `bson:"id" json:"id" validate:"required"`
	Revenue float64 `bson:"revenue,omitempty" json:"revenue"`
}

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	UserName     string             `bson:"userName"`
	TotalPoints  float64            `bson:"totalPoints"`
	TgId         string             `bson:"tgId"`
	FirstName    string             `bson:"firstName"`
	LastName     string             `bson:"lastName"`
	CurPoints    float32            `bson:"curPoints"`
	CountDown    int                `bson:"countDown"`
	LastLogin    time.Time          `bson:"lastLogin"`
	StartFarming time.Time          `bson:"startFarming"`
	Cliamed      bool               `bson:"cliamed"`
	IsStarted    bool               `bson:"isStarted"`
	InviteLink   string             `bson:"inviteLink"`
	IsInvited    bool               `bson:"isInvited"`
	Task         []string           `bson:"task"`
	IntervalId   int                `bson:"intervalId"`
	JoinRank     int                `bson:"joinRank"`
	Style        string             `bson:"style"`
	Friends      []Friend           `bson:"friends,omitempty" json:"friends"`
	// CreatedAt  string `bson:"createdAt"`
	// UpdatedAt  string `bson:"updatedAt"`
}
