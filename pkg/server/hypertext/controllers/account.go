package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	"code.smartsheep.studio/atom/quaso/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
		"/api/accounts/:account",
		v.info,
	)
	router.Get(
		"/api/accounts/self",
		v.gatekeeper.Fn(true, hyperutils.GenScope(), hyperutils.GenPerms()),
		v.self,
	)
}

func (v *AccountController) info(c *fiber.Ctx) error {
	var account models.Account
	if err := v.db.Where("id = ?", c.Params("account", "0")).First(&account).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var subscriberCount int64
	if err := v.db.Model(&models.Subscription{}).Where("provider_id = ?", account.ID).Count(&subscriberCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}
	var subscriptionCount int64
	if err := v.db.Model(&models.Subscription{}).Where("account_id = ?", account.ID).Count(&subscriptionCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	resp, err := v.conn.GetAccountWithID(account.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"user":          resp.User,
		"subscribers":   subscriberCount,
		"subscriptions": subscriptionCount,
	})
}

func (v *AccountController) self(c *fiber.Ctx) error {
	u := c.Locals("quaso-id").(*models.Account)

	var account models.Account
	if err := v.db.Where("id = ?", u.ID).First(&account).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	resp, err := v.conn.GetAccountWithID(account.ID)
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

	return c.JSON(fiber.Map{
		"user":          resp.User,
		"subscribers":   subscriberCount,
		"subscriptions": subscriptionCount,
	})
}
