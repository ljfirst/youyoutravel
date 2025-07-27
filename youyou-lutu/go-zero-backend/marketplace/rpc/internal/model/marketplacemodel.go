package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SnapshotData struct {
	EstimatedCost float64 `bson:"estimated_cost"`
	Days          int32   `bson:"days"`
}

type ItineraryProduct struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID            primitive.ObjectID `bson:"author_id"`
	OriginalItineraryID primitive.ObjectID `bson:"original_itinerary_id"`
	Title               string             `bson:"title"`
	CoverImage          string             `bson:"cover_image"`
	Description         string             `bson:"description"`
	Tags                []string           `bson:"tags"`
	Price               float64            `bson:"price"`
	Status              string             `bson:"status"`
	SalesCount          int32              `bson:"sales_count"`
	AverageRating       float64            `bson:"average_rating"`
	SnapshotData        SnapshotData       `bson:"snapshot_data"`
}

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ProductID  primitive.ObjectID `bson:"product_id"`
	UserID     primitive.ObjectID `bson:"user_id"`
	PurchaseAt primitive.DateTime `bson:"purchase_at"`
}

type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProductID primitive.ObjectID `bson:"product_id"`
	AuthorID  primitive.ObjectID `bson:"author_id"`
	Rating    int32              `bson:"rating"`
	Comment   string             `bson:"comment"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}
