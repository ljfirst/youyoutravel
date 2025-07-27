package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Node struct {
	ID                string `bson:"id"`
	Type              string `bson:"type"`
	ThirdPartyID      string `bson:"third_party_id"`
	UserBudgetedPrice string `bson:"user_budgeted_price"`
	Notes             string `bson:"notes"`
	StartTime         string `bson:"start_time"`
	EndTime           string `bson:"end_time"`
}

type Day struct {
	Date  string `bson:"date"`
	Nodes []Node `bson:"nodes"`
}

type Itinerary struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  primitive.ObjectID `bson:"user_id"`
	Title   string             `bson:"title"`
	StartDate string           `bson:"start_date"`
	EndDate string             `bson:"end_date"`
	Days    []Day              `bson:"days"`
}
