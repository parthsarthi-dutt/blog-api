package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parthsarthi-dutt/blog-api/internal/service"
	"github.com/parthsarthi-dutt/blog-api/internal/utils"
)

type BlogHandler struct {
	blogService *service.BlogService
}

func NewBlogHandler() *BlogHandler {
	return &BlogHandler{
		blogService: service.NewBlogService(),
	}
}

type CreateBlogRequest struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=draft published"`
	PublishDate string `json:"publishDate" binding:"required"` // YYYY-MM-DD
}

func (h *BlogHandler) Create(c *gin.Context) {
	var req CreateBlogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	email := c.GetString("email") // from JWT middleware

	parsed, err := time.Parse("2006-01-02", req.PublishDate)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid date format, expected YYYY-MM-DD")
		return
	}

	err = h.blogService.CreateBlog(
		email,
		req.Title,
		req.Content,
		req.Status,
		parsed.Unix(),
	)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
		return
	}

	utils.Created(c, nil, "Blog created successfully")
}

func (h *BlogHandler) List(c *gin.Context) {
	email := c.GetString("email")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status") // optional filter

	blogs, total, err := h.blogService.ListBlogs(email, page, limit, status)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
		return
	}

	utils.Success(c, gin.H{
		"data":  blogs,
		"page":  page,
		"limit": limit,
		"total": total,
	}, "Blogs fetched")
}

func (h *BlogHandler) Update(c *gin.Context) {
	id := c.Param("id")
	email := c.GetString("email")

	var req CreateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	parsed, err := time.Parse("2006-01-02", req.PublishDate)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid date format, expected YYYY-MM-DD")
		return
	}

	err = h.blogService.UpdateBlog(
		id,
		email,
		req.Title,
		req.Content,
		req.Status,
		parsed.Unix(),
	)
	if err != nil {
		utils.Error(c, http.StatusForbidden, "FORBIDDEN", "You are not allowed to update this blog")
		return
	}

	utils.Success(c, nil, "Blog updated successfully")
}

func (h *BlogHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	email := c.GetString("email")

	err := h.blogService.DeleteBlog(id, email)
	if err != nil {
		utils.Error(c, http.StatusForbidden, "FORBIDDEN", "You are not allowed to delete this blog")
		return
	}

	utils.Success(c, nil, "Blog deleted successfully")
}
