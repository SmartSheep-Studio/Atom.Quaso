package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	"code.smartsheep.studio/atom/quaso/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/middlewares"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"time"
)

type AccountController struct {
	db         *gorm.DB
	conn       *subapps.HeLiCoPtErConnection
	gatekeeper *middlewares.AuthMiddleware
}

func NewAccountController(db *gorm.DB, conn *subapps.HeLiCoPtErConnection, gatekeeper *middlewares.AuthMiddleware) *AccountController {
	return &AccountController{db, conn, gatekeeper}
}

func (v *AccountController) Map(router *fiber.App) {
	router.Get(
		"/api/accounts/self",
		v.gatekeeper.Fn(true, hyperutils.GenScope(), hyperutils.GenPerms()),
		v.self,
	)
	router.Get(
		"/api/accounts/:account",
		v.gatekeeper.Fn(true, hyperutils.GenScope(), hyperutils.GenPerms()),
		v.info,
	)
	router.Get(
		"/api/accounts/:account/posts",
		v.gatekeeper.Fn(true, hyperutils.GenScope(), hyperutils.GenPerms()),
		v.posts,
	)
	router.Post(
		"/api/accounts/:account/subscribe",
		v.gatekeeper.Fn(true, hyperutils.GenScope(), hyperutils.GenPerms()),
		v.subscribe,
	)
}

func (v *AccountController) info(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	var account models.Account
	if err := v.db.Where("id = ?", c.Params("account", "0")).First(&account).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	resp, err := v.conn.GetAccountWithID(account.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var subscriberCount int64
	if err := v.db.Model(&models.Subscription{}).Where("provider_id = ?", account.ID).Count(&subscriberCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}
	var subscriptionCount int64
	if err := v.db.Model(&models.Subscription{}).Where("account_id = ?", account.ID).Count(&subscriptionCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var postCount int64
	if err := v.db.Model(&models.Post{}).Where("account_id = ?", account.ID).Count(&postCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}
	var likeCount int64
	if err := v.db.Model(&models.Like{}).Where("account_id = ?", account.ID).Count(&likeCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var subscribed = false
	if err := v.db.Model(&models.Subscription{}).Where("account_id = ? AND provider_id = ?", u.ID, account.ID).Count(&subscriptionCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else if subscriptionCount > 0 {
		subscribed = true
	}

	return c.JSON(fiber.Map{
		"user":          resp.User,
		"subscribers":   subscriberCount,
		"subscriptions": subscriptionCount,
		"posts":         postCount,
		"likes":         likeCount,

		"is_subscribed": subscribed,
	})
}

func (v *AccountController) self(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	var account models.Account
	if err := v.db.Where("id = ?", u.ID).First(&account).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	resp, err := v.conn.GetAccountWithID(account.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var subscriberCount int64
	if err := v.db.Model(&models.Subscription{}).Where("provider_id = ?", account.ID).Count(&subscriberCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}
	var subscriptionCount int64
	if err := v.db.Model(&models.Subscription{}).Where("account_id = ?", account.ID).Count(&subscriptionCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var postCount int64
	if err := v.db.Model(&models.Post{}).Where("account_id = ?", account.ID).Count(&postCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}
	var likeCount int64
	if err := v.db.Model(&models.Like{}).Where("account_id = ?", account.ID).Count(&likeCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var subscribed = false
	var uSubscriptionCount int64
	if err := v.db.Model(&models.Subscription{}).Where("account_id = ? AND provider_id = ?", u.ID, account.ID).Count(&uSubscriptionCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else if uSubscriptionCount > 0 {
		subscribed = true
	}

	return c.JSON(fiber.Map{
		"user":          resp.User,
		"subscribers":   subscriberCount,
		"subscriptions": subscriptionCount,
		"posts":         postCount,
		"likes":         likeCount,

		"is_subscribed": subscribed,
	})
}

func (v *AccountController) subscribe(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	var account models.Account
	if err := v.db.Where("id = ?", c.Params("account", "0")).First(&account).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var subscriptionCount int64
	if err := v.db.Model(&models.Subscription{}).Where("account_id = ? AND provider_id = ?", u.ID, account.ID).Count(&subscriptionCount).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return hyperutils.ErrorParser(err)
		}
	}

	if subscriptionCount > 0 {
		// Cancel subscription
		if err := v.db.Where("account_id = ? AND provider_id = ?", u.ID, account.ID).Delete(&models.Subscription{}).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else {
			return c.SendStatus(fiber.StatusNoContent)
		}
	} else {
		subscription := models.Subscription{
			AccountID:  u.ID,
			ProviderID: account.ID,
		}

		if err := v.db.Save(&subscription).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else {
			return c.SendStatus(fiber.StatusOK)
		}
	}
}

func (v *AccountController) posts(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	var account models.Account
	if err := v.db.Where("id = ?", c.Params("account", "0")).First(&account).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	tx := v.db.Where("published_at <= ?", time.Now())

	tx.Where("is_hidden = ?", false)
	tx.Order("published_at desc")
	tx.Where("account_id = ?", account.ID)

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
