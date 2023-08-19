package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/middleware"
	"github.com/gofiber/fiber/v2"
)

type AttachmentController struct {
	conn       *subapps.HeLiCoPtErConnection
	gatekeeper *middleware.AuthMiddleware
}

func NewAttachmentController(conn *subapps.HeLiCoPtErConnection, gatekeeper *middleware.AuthMiddleware) *AttachmentController {
	return &AttachmentController{conn, gatekeeper}
}

func (v *AttachmentController) Map(router *fiber.App) {
	router.Post(
		"/api/attachments",
		v.gatekeeper.Fn(true, hyperutils.GenScope("create:attachments"), hyperutils.GenPerms("posts.attachments")),
		v.create,
	)
}

func (v *AttachmentController) create(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	io, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if item, err := v.conn.UploadAssets2User(c.Locals("principal-token").(string), file.Filename, io); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(fiber.Map{
			"file": item,
			"url":  item.GetURL(),
		})
	}
}
