package traverse

import (
	"context"

	"github.com/snivilised/extendio/bus"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/kernel"
	"github.com/snivilised/traverse/pref"
)

const (
	badge = "navigation-driver"
)

type duffNavigator struct {
	root    string
	client  core.Client
	from    string
	options []pref.Option
}

func (n *duffNavigator) Navigate() (core.TraverseResult, error) {
	return types.NavigateResult{}, nil
}

func init() {
	h := bus.Handler{
		Handle: func(_ context.Context, m bus.Message) {
			m.Data.(types.ContextExpiry).Expired()
		},
		Matcher: services.TopicContextExpired,
	}

	services.Broker.RegisterHandler(badge, h)
}

// replaces the runner in extendio, although it not
// entirely the same as we also need to maintain a clear
// separation between using the reactive model and the
// sequential model.

type Driver interface {
	Primary(root string, client core.Client, options ...pref.Option) kernel.Navigator
	Resume(from string) kernel.Navigator
}

func Walk() Driver {
	return &driver{
		session: &session{
			ctx: context.Background(),
		},
	}
}

func Run() Driver {
	return &driver{
		session: &session{},
	}
}

// the driver needs to tap into the event bus via the broker
// so that it can see when the context has expired.
type driver struct {
	session *session
}

func (d *driver) Primary(root string,
	client core.Client, options ...pref.Option,
) kernel.Navigator {
	return &duffNavigator{
		root:    root,
		client:  client,
		options: options,
	}
}

func (d *driver) Resume(from string) kernel.Navigator {
	return &duffNavigator{
		from: from,
	}
}

func (d *driver) Expired() {
	// not sure if this is right; do we also need to invoke the cancel func
	// if it is defined?
	//
	d.session.ctx = nil
}
