package service

import (
	"errors"
	"time"

	"github.com/parthsarthi-dutt/blog-api/internal/models"
	"github.com/parthsarthi-dutt/blog-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogService struct {
	blogRepo *repository.BlogRepository
}

func NewBlogService() *BlogService {
	return &BlogService{
		blogRepo: repository.NewBlogRepository(),
	}
}

func (s *BlogService) CreateBlog(
	email, title, content, category string,tags []string,
	publishUnix int64,
) error {
	blog := models.Blog{
		UserEmail:   email,
		Title:       title,
		Content:     content,
		Category: category,
		Tags: tags,
		PublishDate: publishUnix,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	return s.blogRepo.Create(blog)
}

func (s *BlogService) ListBlogs(
	email string,
	page, limit int,
	category string,
) ([]models.Blog, int64, error) {

	skip := (page - 1) * limit
	filter := bson.M{"user_email": email}

	if category != "" {
		filter["category"] = category
	}

	return s.blogRepo.FindByUser(
		email,
		int64(skip),
		int64(limit),
		filter,
	)
}

func (s *BlogService) UpdateBlog(
	idStr, email, title, content, category string,
	publishUnix int64,
) error {

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return errors.New("invalid id")
	}

	update := bson.M{
		"title":         title,
		"content":       content,
		"category":        category,
		"publish_date":  publishUnix,
		"updated_at":    time.Now().Unix(),
	}

	return s.blogRepo.Update(id, email, update)
}

func (s *BlogService) DeleteBlog(idStr, email string) error {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return errors.New("invalid id")
	}

	return s.blogRepo.Delete(id, email)
}
