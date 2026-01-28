package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/parthsarthi-dutt/blog-api/internal/config"
	"github.com/parthsarthi-dutt/blog-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepository() *BlogRepository {
	return &BlogRepository{
		collection: config.DB.Collection("blogs"),
	}
}

func (r *BlogRepository) Create(blog models.Blog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, blog)
	return err
}

func (r *BlogRepository) FindByUser(email string, skip, limit int64, filter bson.M) ([]models.Blog, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter["user_email"] = email

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, 0, err
	}


	var blogs []models.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, 0, err
	}

	count, _ := r.collection.CountDocuments(ctx, filter)

	return blogs, count, nil
}

func (r *BlogRepository) Update(id primitive.ObjectID, email string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": id, "user_email": email},
		bson.M{"$set": update},
	)
	return err
}

func (r *BlogRepository) Delete(id primitive.ObjectID, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id, "user_email": email})
	return err
}
