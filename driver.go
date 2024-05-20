package tv

import (
	"context"

	"github.com/snivilised/extendio/bus"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/types"
)

const (
	badge = "badge: navigation-driver"
)

type driver struct {
	s session
}

func (d *driver) init() {
	services.Broker.RegisterHandler(badge, bus.Handler{
		Handle: func(_ context.Context, m bus.Message) {
			m.Data.(types.ContextExpiry).Expired()
		},
		Matcher: services.TopicContextExpired,
	})

	services.Broker.RegisterHandler(badge, bus.Handler{
		Handle: func(_ context.Context, m bus.Message) {
			_ = m.Data
			// now invoke session.finish
		},
		Matcher: services.TopicTraverseResult,
	})
}

func (d *driver) Navigate() (core.TraverseResult, error) {
	d.init()
	d.s.start()
	result, err := d.s.exec()

	d.s.finish(result)

	return result, err
}
