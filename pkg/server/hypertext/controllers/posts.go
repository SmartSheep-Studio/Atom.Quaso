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
	tx := v.db.Where("published_at <= ?", time.Now())

	tx.Order("created_at desc")

	if c.Query("type", "none") != "none" {
		tx.Where("type = ?", c.Query("type"))
	}

	var postCount int64
	var posts []models.Post
	if err := tx.Model(&models.Post{}).Count(&postCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else if err := tx.Offset(c.QueryInt("skip", 0)).Limit(5).Find(&posts).Error; err != nil {
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
		"posts": lo.Map(posts, func(item models.Post, index int) map[string]any {
			data := hyperutils.CovertStructToMap(item)

			data["account"], _ = lo.Find(authors, func(v models.Account) bool {
				return v.ID == item.AccountID
			})

			var commentCount int64
			if err := v.db.Model(&models.Post{}).Where("belong_id = ?", item.ID).Count(&commentCount).Error; err != nil {
				data["comment_count"] = 0
			} else {
				data["comment_count"] = commentCount
			}

			return data
		}),
	})
}

func (v *PostController) get(c *fiber.Ctx) error {
	tx := v.db.Preload("Comments")

	tx.Where("id = ?", c.Params("post"))
	tx.Where("published_at <= ?", time.Now())

	var post models.Post
	if err := tx.First(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var commentAuthors []models.Account
	if err := v.db.Where("id IN ?", lo.Union(lo.Map(post.Comments, func(item models.Post, index int) uint {
		return item.AccountID
	}))).Find(&commentAuthors).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	return c.JSON(func() map[string]any {
		data := hyperutils.CovertStructToMap(post)

		var author models.Account
		if err := v.db.Where("id = ?", post.AccountID).First(&author).Error; err != nil {
			data["account"] = nil
		} else {
			data["account"] = author
		}

		var commentCount int64
		if err := v.db.Model(&models.Post{}).Where("belong_id = ?", post.ID).Count(&commentCount).Error; err != nil {
			data["comment_count"] = 0
		} else {
			data["comment_count"] = commentCount
		}

		data["comments"] = lo.Map(post.Comments, func(item models.Post, index int) map[string]any {
			data := hyperutils.CovertStructToMap(item)

			data["account"], _ = lo.Find(commentAuthors, func(v models.Account) bool {
				return v.ID == item.AccountID
			})

			var commentCount int64
			if err := v.db.Model(&models.Post{}).Where("belong_id = ?", item.ID).Count(&commentCount).Error; err != nil {
				data["comment_count"] = 0
			} else {
				data["comment_count"] = commentCount
			}

			return data
		})

		return data
	}())
}

func (v *PostController) create(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	var req struct {
		Type        string     `json:"type" validate:"required"`
		Content     string     `json:"content" validate:"required"`
		Tags        []string   `json:"tags"`
		Attachments []string   `json:"attachments"`
		BelongTo    *uint      `json:"belong_to"`
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

	if req.BelongTo != nil {
		var parentPost models.Post
		if err := v.db.Where("id = ?", req.BelongTo).First(&parentPost).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else {
			post.BelongID = lo.ToPtr(parentPost.ID)
		}
	}

	if err := v.db.Save(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(post)
	}
}
