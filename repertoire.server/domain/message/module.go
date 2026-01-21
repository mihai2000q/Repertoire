package message

import (
	"repertoire/server/domain/message/handler"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/fx"
)

var Module = fx.Options(
	handler.Module,
	fx.Provide(NewRouter),
	fx.Invoke(func(router *message.Router) {}),
)
