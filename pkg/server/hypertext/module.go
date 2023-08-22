package hypertext

import (
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/controllers"
	"code.smartsheep.studio/atom/quaso/pkg/server/hypertext/middlewares"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("http",
		middlewares.Module(),
		controllers.Module(),

		fx.Invoke(fx.Annotate(MapControllers, fx.ParamTags(`group:"hypertext_controllers"`))),
		fx.Provide(NewHttpServer),
	)
}
