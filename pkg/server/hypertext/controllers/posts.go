package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	"code.smartsheep.studio/atom/quaso/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/middlewares"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type PostController struct {
	db         *gorm.DB
	conn       *subapps.HeLiCoPtErConnection
	gatekeeper *middlewares.AuthMiddleware
}

func NewPostController(db *gorm.DB, conn *subapps.HeLiCoPtErConnection, gatekeeper *middlewares.AuthMiddleware) *PostController {
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
	router.Put(
		"/api/posts/:post",
		v.gatekeeper.Fn(true, hyperutils.GenScope("update:posts"), hyperutils.GenPerms("posts.update")),
		v.update,
	)
	router.Delete(
		"/api/posts/:post",
		v.gatekeeper.Fn(true, hyperutils.GenScope("delete:posts"), hyperutils.GenPerms("posts.delete")),
		v.delete,
	)
	router.Post(
		"/api/posts/:post/like",
		v.gatekeeper.Fn(true, hyperutils.GenScope("like:posts"), hyperutils.GenPerms("posts.like")),
		v.like,
	)
	router.Post(
		"/api/posts/:post/dislike",
		v.gatekeeper.Fn(true, hyperutils.GenScope("dislike:posts"), hyperutils.GenPerms("posts.dislike")),
		v.dislike,
	)
}

func (v *PostController) list(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	tx := v.db.Where("published_at <= ?", time.Now())

	tx.Where("is_hidden = ?", false)
	tx.Order("published_at desc")

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

	accountsResp := map[uint]subapps.HeLiCoPtErAccountResp{}
	lo.ForEach(authors, func(item models.Account, index int) {
		resp, _ := v.conn.GetAccountWithID(item.UserID)
		accountsResp[item.ID] = resp
	})

	return c.JSON(fiber.Map{
		"total": postCount,
		"posts": lo.Map(posts, func(item models.Post, index int) map[string]any {
			data := hyperutils.CovertStructToMap(item)

			account, _ := lo.Find(authors, func(v models.Account) bool {
				return v.ID == item.AccountID
			})

			data["account"] = account
			data["author"] = accountsResp[account.ID].User

			var commentCount int64
			if err := v.db.Model(&item).Where("belong_id = ?", item.ID).Count(&commentCount).Error; err != nil {
				data["comment_count"] = 0
			} else {
				data["comment_count"] = commentCount
			}

			var likeCount int64
			tx = v.db.Model(&models.Like{}).Where("post_id = ?", item.ID)
			if err := tx.Where("account_id = ?", u.ID).Count(&likeCount).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					data["is_liked"] = false
				}
			} else {
				data["is_liked"] = lo.Ternary(likeCount > 0, true, false)
			}

			var dislikeCount int64
			tx = v.db.Model(&models.Dislike{}).Where("post_id = ?", item.ID)
			if err := tx.Where("account_id = ?", u.ID).Count(&dislikeCount).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					data["is_disliked"] = false
				}
			} else {
				data["is_disliked"] = lo.Ternary(dislikeCount > 0, true, false)
			}

			tx = v.db.Model(&models.Like{}).Where("post_id = ?", item.ID)
			if err := tx.Count(&likeCount).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					data["like_count"] = 0
				}
			} else {
				data["like_count"] = likeCount
			}

			tx = v.db.Model(&models.Dislike{}).Where("post_id = ?", item.ID)
			if err := tx.Count(&dislikeCount).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					data["dislike_count"] = 0
				}
			} else {
				data["dislike_count"] = dislikeCount
			}

			return data
		}),
	})
}

func (v *PostController) get(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	tx := v.db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("published_at desc")
	})

	tx.Where("id = ?", c.Params("post"))
	tx.Where("is_hidden = ?", false)
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

	accountsResp := map[uint]subapps.HeLiCoPtErAccountResp{}
	lo.ForEach(commentAuthors, func(item models.Account, index int) {
		resp, _ := v.conn.GetAccountWithID(item.UserID)
		accountsResp[item.ID] = resp
	})

	return c.JSON(func() map[string]any {
		data := hyperutils.CovertStructToMap(post)

		var author models.Account
		if err := v.db.Where("id = ?", post.AccountID).First(&author).Error; err != nil {
			data["account"] = nil
		} else {
			resp, _ := v.conn.GetAccountWithID(author.UserID)

			data["account"] = author
			data["author"] = resp.User
		}

		var commentCount int64
		if err := v.db.Model(&models.Post{}).Where("belong_id = ?", post.ID).Count(&commentCount).Error; err != nil {
			data["comment_count"] = 0
		} else {
			data["comment_count"] = commentCount
		}

		var likeCount int64
		tx = v.db.Model(&models.Like{}).Where("post_id = ?", post.ID)
		if err := tx.Where("account_id = ?", u.ID).Count(&likeCount).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				data["is_liked"] = false
			}
		} else {
			data["is_liked"] = lo.Ternary(likeCount > 0, true, false)
		}

		var dislikeCount int64
		tx = v.db.Model(&models.Dislike{}).Where("post_id = ?", post.ID)
		if err := tx.Where("account_id = ?", u.ID).Count(&dislikeCount).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				data["is_disliked"] = false
			}
		} else {
			data["is_disliked"] = lo.Ternary(dislikeCount > 0, true, false)
		}

		tx = v.db.Model(&models.Like{}).Where("post_id = ?", post.ID)
		if err := tx.Count(&likeCount).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				data["like_count"] = 0
			}
		} else {
			data["like_count"] = likeCount
		}

		tx = v.db.Model(&models.Dislike{}).Where("post_id = ?", post.ID)
		if err := tx.Count(&dislikeCount).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				data["dislike_count"] = 0
			}
		} else {
			data["dislike_count"] = dislikeCount
		}

		data["comments"] = lo.Map(post.Comments, func(item models.Post, index int) map[string]any {
			data := hyperutils.CovertStructToMap(item)

			account, _ := lo.Find(commentAuthors, func(v models.Account) bool {
				return v.ID == item.AccountID
			})

			data["account"] = account
			data["author"] = accountsResp[account.ID].User

			var commentCount int64
			if err := v.db.Model(&models.Post{}).Where("belong_id = ?", item.ID).Count(&commentCount).Error; err != nil {
				data["comment_count"] = 0
			} else {
				data["comment_count"] = commentCount
			}

			var likeCount int64
			tx = v.db.Model(&models.Like{}).Where("post_id = ?", item.ID)
			if err := tx.Where("account_id = ?", u.ID).Count(&likeCount).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					data["is_liked"] = false
				}
			} else {
				data["is_liked"] = lo.Ternary(likeCount > 0, true, false)
			}

			var dislikeCount int64
			tx = v.db.Model(&models.Dislike{}).Where("post_id = ?", item.ID)
			if err := tx.Where("account_id = ?", u.ID).Count(&dislikeCount).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					data["is_disliked"] = false
				}
			} else {
				data["is_disliked"] = lo.Ternary(dislikeCount > 0, true, false)
			}

			tx = v.db.Model(&models.Like{}).Where("post_id = ?", item.ID)
			if err := tx.Count(&likeCount).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					data["like_count"] = 0
				}
			} else {
				data["like_count"] = likeCount
			}

			tx = v.db.Model(&models.Dislike{}).Where("post_id = ?", item.ID)
			if err := tx.Count(&dislikeCount).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					data["dislike_count"] = 0
				}
			} else {
				data["dislike_count"] = dislikeCount
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
		Scope       *string    `json:"scope"`
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
		Scope:       lo.Ternary(req.Scope == nil, "plaza", lo.FromPtr(req.Scope)),
		Content:     req.Content,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		Attachments: datatypes.NewJSONSlice(req.Attachments),
		PublishedAt: lo.Ternary(req.PublishedAt == nil, time.Now(), lo.FromPtr(req.PublishedAt)),
		IpAddress:   c.IP(),
		IsHidden:    false,
		IsEdited:    false,
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

func (v *PostController) update(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	var req struct {
		Type        string     `json:"type" validate:"required"`
		Content     string     `json:"content" validate:"required"`
		Scope       *string    `json:"scope"`
		Tags        []string   `json:"tags"`
		Attachments []string   `json:"attachments"`
		PublishedAt *time.Time `json:"published_at"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	tx := v.db.Preload("Comments")

	tx.Where("id = ?", c.Params("post"))
	tx.Where("account_id = ?", u.ID)

	var post models.Post
	if err := tx.First(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	post.Type = req.Type
	post.Scope = lo.Ternary(req.Scope == nil, "plaza", lo.FromPtr(req.Scope))
	post.Tags = datatypes.NewJSONSlice(req.Tags)
	post.Attachments = datatypes.NewJSONSlice(req.Attachments)
	post.Content = req.Content
	post.PublishedAt = lo.Ternary(req.PublishedAt == nil, time.Now(), lo.FromPtr(req.PublishedAt))
	post.IsEdited = true

	if err := v.db.Save(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(post)
	}
}

func (v *PostController) delete(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	tx := v.db.Preload("Comments")

	tx.Where("id = ?", c.Params("post"))
	tx.Where("account_id = ?", u.ID)

	var post models.Post
	if err := tx.First(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	if err := v.db.Delete(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(post)
	}
}

func (v *PostController) like(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	tx := v.db.Where("id = ?", c.Params("post", "0"))

	tx.Where("is_hidden = ?", false)
	tx.Where("published_at <= ?", time.Now())

	var post models.Post
	if err := tx.First(&post).Preload("Likes").Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var likeCount int64
	if err := v.db.Model(&models.Like{}).Where("post_id = ? AND account_id = ?", post.ID, u.ID).Count(&likeCount).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return hyperutils.ErrorParser(err)
		}
	}

	if likeCount > 0 {
		// Cancel like
		if err := v.db.Where("post_id = ? AND account_id = ?", post.ID, u.ID).Delete(&models.Like{}).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else {
			return c.SendStatus(fiber.StatusNoContent)
		}
	} else {
		// Like
		like := models.Like{
			AccountID: u.ID,
			PostID:    post.ID,
		}

		if err := v.db.Save(&like).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else {
			return c.SendStatus(fiber.StatusOK)
		}
	}
}

func (v *PostController) dislike(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	tx := v.db.Where("id = ?", c.Params("post", "0"))

	tx.Where("is_hidden = ?", false)
	tx.Where("published_at <= ?", time.Now())

	var post models.Post
	if err := tx.First(&post).Preload("Dislikes").Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var dislikeCount int64
	if err := v.db.Model(&models.Dislike{}).Where("post_id = ? AND account_id = ?", post.ID, u.ID).Count(&dislikeCount).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return hyperutils.ErrorParser(err)
		}
	}

	if dislikeCount > 0 {
		// Cancel dislike
		if err := v.db.Where("post_id = ? AND account_id = ?", post.ID, u.ID).Delete(&models.Dislike{}).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else {
			return c.SendStatus(fiber.StatusNoContent)
		}
	} else {
		// Dislike
		dislike := models.Dislike{
			AccountID: u.ID,
			PostID:    post.ID,
		}

		if err := v.db.Save(&dislike).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else {
			return c.SendStatus(fiber.StatusOK)
		}
	}
}
