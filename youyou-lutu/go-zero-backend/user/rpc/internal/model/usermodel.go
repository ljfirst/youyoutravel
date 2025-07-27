package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	Nickname            string             `bson:"nickname,omitempty"`
	AvatarURL           string             `bson:"avatar_url,omitempty"`
	Gender              string             `bson:"gender,omitempty"`
	Bio                 string             `bson:"bio,omitempty"`
	PublishedItineraries []primitive.ObjectID `bson:"published_itineraries,omitempty"`
	PurchasedItineraries []primitive.ObjectID `bson:"purchased_itineraries,omitempty"`
}
