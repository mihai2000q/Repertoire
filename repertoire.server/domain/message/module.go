package message

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/fx"
	"repertoire/server/domain/message/handler"
)

var Module = fx.Options(
	handler.Module,
	fx.Provide(NewRouter),
	fx.Invoke(func(router *message.Router) {}),
)
