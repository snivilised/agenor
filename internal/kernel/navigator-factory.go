package kernel

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
)

func PrimeArtefacts(inception *Inception,
	sealer enclave.GuardianSealer,
) *Artefacts {
	mediator, err := NewMediator(inception, sealer)

	return &Artefacts{
		Kontroller: mediator,
		Mediator:   mediator,
		Resources:  inception.Resources,
		Error:      err,
	}
}

func newImpl(o *pref.Options,
	inception *Inception,
) (impl NavigatorImpl, err error) {
	subscription := inception.Subscription

	agent := navigatorAgent{
		ao: &agentOptions{
			hooks:   &o.Hooks,
			defects: &o.Defects,
		},
		ro: &readOptions{
			hooks: readHooks{
				read: o.Hooks.ReadDirectory,
				sort: o.Hooks.Sort,
			},
			behaviour: &o.Behaviours.Sort,
		},
		resources: inception.Resources,
		persister: author{
			o:     o,
			perms: core.Perms,
		},
		magnitude: inception.Facade.Magnitude(),
	}

	switch subscription {
	case enums.SubscribeFiles:
		impl = &navigatorFiles{
			navigatorAgent: agent,
		}

	case enums.SubscribeDirectories, enums.SubscribeDirectoriesWithFiles:
		impl = &navigatorDirectories{
			navigatorAgent: agent,
		}

	case enums.SubscribeUniversal:
		impl = &navigatorUniversal{
			navigatorAgent: agent,
		}

	case enums.SubscribeUndefined:
		err = locale.ErrUsageMissingSubscription
		impl = &navigatorUniversal{
			navigatorAgent: agent,
		}
	}

	return impl, err
}
