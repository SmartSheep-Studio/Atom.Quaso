package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	"code.smartsheep.studio/atom/quaso/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type PostController struct {
	db         *gorm.DB
	conn       *subapps.HeLiCoPtErConnection
	gatekeeper *middleware.AuthMiddleware
}

func NewPostController(db *gorm.DB, conn *subapps.HeLiCoPtErConnection, gatekeeper *middleware.AuthMiddleware) *PostController {
	return &PostController{db, conn, gatekeeper}
}

func (v *PostController) Map(router *fiber.App) {
	router.Get(
		"/api/posts",
		v.gatekeeper.Fn(true, hyperutils.GenScope("read:posts"), hyperutils.GenPerms("posts.read")),
		v.list,
	)
	router.Get(
		"/api/posts/:post",
		v.gatekeeper.Fn(true, hyperutils.GenScope("read:posts"), hyperutils.GenPerms("posts.read")),
		v.get,
	)
	router.Post(
		"/api/posts",
		v.gatekeeper.Fn(true, hyperutils.GenScope("create:posts"), hyperutils.GenPerms("posts.create")),
		v.create,
	)
}

func (v *PostController) list(c *fiber.Ctx) error {
	tx := v.db.Offset(c.QueryInt("skip", 0)).Limit(20)

	tx.Where("published_at <= ?", time.Now())
	tx.Order("created_at desc")

	var postCount int64
	var posts []models.Post
	if err := tx.Model(&models.Post{}).Count(&postCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else if err := tx.Find(&posts).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var authors []models.Account
	if err := v.db.Where("id IN ?", lo.Union(lo.Map(posts, func(item models.Post, index int) uint {
		return item.AccountID
	}))).Find(&authors).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	return c.JSON(fiber.Map{
		"total": postCount,
		"posts": posts,
		"related_authors": func() map[uint]models.Account {
			data := make(map[uint]models.Account)
			for _, item := range authors {
				data[item.ID] = item
			}
			return data
		}(),
	})
}

func (v *PostController) get(c *fiber.Ctx) error {
	tx := v.db.Preload("Comments")

	tx.Where("id = ?", c.Params("post"))

	var post models.Post
	if err := tx.First(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(post)
	}
}

func (v *PostController) create(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	var req struct {
		Type        string     `json:"type" validate:"required"`
		Content     string     `json:"content" validate:"required"`
		Tags        []string   `json:"tags"`
		Attachments []string   `json:"attachments"`
		PublishedAt *time.Time `json:"published_at"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	post := models.Post{
		Type:        req.Type,
		Content:     req.Content,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		Attachments: datatypes.NewJSONSlice(req.Attachments),
		PublishedAt: lo.Ternary(req.PublishedAt == nil, time.Now(), lo.FromPtr(req.PublishedAt)),
		IpAddress:   c.IP(),
		AccountID:   u.ID,
	}

	if err := v.db.Save(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(post)
	}
}
