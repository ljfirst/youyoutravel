package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Destination struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name"`
	Address        string             `bson:"address"`
	City           string             `bson:"city"`
	GeoCoordinate  string             `bson:"geo_coordinate"`
	Tags           []string           `bson:"tags"`
	OfficialIntro  string             `bson:"official_intro"`
}
