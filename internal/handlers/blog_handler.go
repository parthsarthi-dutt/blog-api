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
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Status   string   `json:"status" binding:"required,oneof=draft published"`
	Category string   `json:"category" binding:"required"`
	Tags     []string `json:"tags"`
}


func (h *BlogHandler) Create(c *gin.Context) {
	var req CreateBlogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	email := c.GetString("email")

	var publishUnix int64
	if req.Status == "published" {
		loc, _ := time.LoadLocation("Asia/Kolkata")
		publishUnix = time.Now().In(loc).Unix()
	}

	err := h.blogService.CreateBlog(
		email,
		req.Title,
		req.Content,
		req.Status,
		req.Category,
		req.Tags,
		publishUnix,
	)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
		return
	}

	utils.Created(c, nil, "Blog created successfully")
}
func (h *BlogHandler) ListPublic(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")

	blogs, total, err := h.blogService.ListPublicBlogs(
		page,
		limit,
		category,
	)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
		return
	}

	utils.Success(c, gin.H{
		"data":  blogs,
		"page":  page,
		"limit": limit,
		"total": total,
	}, "Public blogs fetched")
}


func (h *BlogHandler) List(c *gin.Context) {
	email := c.GetString("email")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")
	status:=c.Query("status")

	blogs, total, err := h.blogService.ListBlogs(email, page, limit, category,status)
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

	var publishUnix int64
	if req.Status == "published" {
		loc, _ := time.LoadLocation("Asia/Kolkata")
		publishUnix = time.Now().In(loc).Unix()
	}

	err := h.blogService.UpdateBlog(
		id,
		email,
		req.Title,
		req.Content,
		req.Status,
		req.Category,
		req.Tags,
		publishUnix,
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
