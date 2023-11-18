package handlers

import (
	db "backend-development/pkg/database"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id" validate:"required"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at" validate:"required"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at" validate:"required"`
	Text      string             `json:"text" bson:"text" validate:"required"`
	URL       string             `json:"url" bson:"url"`
	Hashtags  string             `json:"hashtags" bson:"hashtags"`
	Via       string             `json:"via" bson:"via"`
	Related   string             `json:"related" bson:"related"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidatePostStruct(p Post) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(p)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}

func CreatePost(c *fiber.Ctx) error {
	post := Post{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := c.BodyParser(&post); err != nil {
		return err
	}

	errors := ValidatePostStruct(post)
	if errors != nil {
		return c.JSON(errors)
	}

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(db.Database).Collection(string(db.PostsCollection))

	_, err = collection.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}

	return c.JSON(post)
}

func GetAllPosts(c *fiber.Ctx) error {
	client, err := db.GetMongoClient()

	var posts []*Post
	if err != nil {
		return err
	}

	collection := client.Database(db.Database).Collection(string(db.PostsCollection))

	cur, err := collection.Find(context.TODO(), bson.D{
		primitive.E{},
	})
	if err != nil {
		return err
	}

	for cur.Next(context.TODO()) {
		var p Post
		err := cur.Decode(&p)
		if err != nil {
			return err
		}

		posts = append(posts, &p)
	}

	return c.JSON(posts)
}

func GetSinglePost(c *fiber.Ctx) error {
	id := c.Params("id")

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(db.Database).Collection(string(db.PostsCollection))

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	result := collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return result.Err()
	}

	var post Post
	err = result.Decode(&post)
	if err != nil {
		return err
	}

	return c.JSON(post)
}

func UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(db.Database).Collection(string(db.PostsCollection))

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	var updatedPost bson.M
	if err := c.BodyParser(&updatedPost); err != nil {
		return err
	}

	updatedPost["updated_at"] = time.Now()

	update := bson.D{
		{Key: "$set", Value: updatedPost},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return c.JSON(updatedPost)
}

func DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(db.Database).Collection(string(db.PostsCollection))

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
