package profiles

import (
	"context"
	"net/http"
	"time"

	"questweaver_pro_backend/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const collectionName = "profiles"

func collection() *mongo.Collection {
	return database.GetCollection(collectionName)
}

// CreateProfile handles POST /profiles
func CreateProfile(c *gin.Context) {
	var profile Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if profile.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if a profile with this userId already exists.
	var existing Profile
	err := collection().FindOne(ctx, bson.M{"userId": profile.UserID}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "profile with this userId already exists"})
		return
	}

	result, err := collection().InsertOne(ctx, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
		return
	}

	profile.ID = result.InsertedID.(bson.ObjectID)
	c.JSON(http.StatusCreated, profile)
}

// GetProfile handles GET /profiles/:userId
func GetProfile(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var profile Profile
	err := collection().FindOne(ctx, bson.M{"userId": userId}).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile handles PUT /profiles/:userId
func UpdateProfile(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	var update ProfileUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setFields := bson.M{}
	if update.PreferredName != nil {
		setFields["preferredName"] = *update.PreferredName
	}

	if len(setFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"userId": userId}
	result, err := collection().UpdateOne(ctx, filter, bson.M{"$set": setFields})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	// Return the updated document.
	var profile Profile
	err = collection().FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// DeleteProfile handles DELETE /profiles/:userId
func DeleteProfile(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection().DeleteOne(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete profile"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile deleted"})
}
