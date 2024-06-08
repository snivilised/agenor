package hiber

import (
	"context"

	"github.com/snivilised/extendio/bus"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options) types.Plugin {
	active := o.Core.Hibernate.Wake != nil

	if active {
		return &Plugin{}
	}

	return nil
}

type Plugin struct {
}

func (p *Plugin) Name() string {
	return "hibernation"
}

func (p *Plugin) Init() error {
	h := bus.Handler{
		Handle: func(_ context.Context, m bus.Message) {
			if in, ok := m.Data.(types.UsePlugin); ok {
				_ = in.Interceptor() // TODO: call Intercept
				_ = in.Facilitate()
			}
		},
		Matcher: services.TopicInterceptNavigator,
	}

	services.Broker.RegisterHandler(badge, h)

	return nil
}
