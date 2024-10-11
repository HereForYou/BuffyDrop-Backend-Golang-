package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Friend struct {
	Id      string  `bson:"id" json:"id" validate:"required"`
	Revenue float64 `bson:"revenue,omitempty" json:"revenue"`
}

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserName     string             `bson:"userName" json:"userName"`
	TotalPoints  float64            `bson:"totalPoints" json:"totalPoints"`
	TgId         string             `bson:"tgId" json:"tgId" validate:"required"`
	FirstName    string             `bson:"firstName" json:"firstName"`
	LastName     string             `bson:"lastName" json:"lastName"`
	CurPoints    float32            `bson:"curPoints" json:"curPoints"`
	CountDown    int                `bson:"countDown" json:"countDown"`
	LastLogin    time.Time          `bson:"lastLogin" json:"lastLogin"`
	StartFarming time.Time          `bson:"startFarming" json:"startFarming"`
	Cliamed      bool               `bson:"cliamed" json:"cliamed"`
	IsStarted    bool               `bson:"isStarted" json:"isStarted"`
	InviteLink   string             `bson:"inviteLink" json:"inviteLink"`
	IsInvited    bool               `bson:"isInvited" json:"isInvited"`
	Task         []string           `bson:"task" json:"task"`
	IntervalId   int                `bson:"intervalId" json:"intervalId"`
	JoinRank     int                `bson:"joinRank" json:"joinRank"`
	Style        string             `bson:"style" json:"style"`
	Friends      []Friend           `bson:"friends" json:"friends"`
	// CreatedAt  string `bson:"createdAt"`
	// UpdatedAt  string `bson:"updatedAt"`
}
