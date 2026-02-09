package profiles

import "go.mongodb.org/mongo-driver/v2/bson"

// Profile represents a user profile document in MongoDB.
type Profile struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string        `bson:"userId"        json:"userId"`
	PreferredName string        `bson:"preferredName" json:"preferredName"`
}

// ProfileUpdate contains the fields that can be updated on a profile.
type ProfileUpdate struct {
	PreferredName *string `json:"preferredName"`
}
