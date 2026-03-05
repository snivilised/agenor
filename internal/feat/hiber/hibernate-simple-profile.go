package hiber

import (
	"io/fs"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/filtering"
	"github.com/snivilised/agenor/life"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/li18ngo"
)

type simple struct {
	common
	states  hibernateStates
	current state
}

func (p *simple) init(controls *life.Controls) error {
	p.states = p.create()
	p.controls = controls

	if p.ho.WakeAt != nil {
		filter, err := filtering.New(p.ho.WakeAt, p.fo)
		if err != nil {
			return err
		}

		p.triggers.wake = filter

		if p.triggers.sleep == nil {
			p.triggers.sleep = filtering.NewProhibitiveTraverseFilter(
				&core.FilterDef{
					Description: li18ngo.Text(locale.ProhibitiveWordTemplData{}),
				},
			)
		}
	}

	if p.ho.SleepAt != nil {
		filter, err := filtering.New(p.ho.SleepAt, p.fo)
		if err != nil {
			return err
		}

		p.triggers.sleep = filter

		if p.triggers.wake == nil {
			p.triggers.wake = filtering.NewPermissiveTraverseFilter(
				&core.FilterDef{
					Description: li18ngo.Text(locale.PermissiveWordTemplData{}),
				},
			)
		}
	}

	p.transition(launch(p.ho))

	return nil
}

func (p *simple) transition(en enums.Hibernation) {
	p.current = p.states[en]
}

func (p *simple) next(servant core.Servant, node *core.Node, inspection enclave.Inspection) (bool, error) {
	return p.current.next(servant, node, inspection)
}

func (p *simple) create() hibernateStates {
	return hibernateStates{
		enums.HibernationPending: state{
			next: func(_ core.Servant, node *core.Node, _ enclave.Inspection) (bool, error) {
				if p.triggers.wake.IsMatch(node) {
					p.controls.Wake.Dispatch()(p.triggers.wake.Description())
					p.transition(enums.HibernationActive)

					if p.ho.Behaviour.InclusiveWake {
						return true, nil
					}
				}

				return false, nil
			},
		},

		enums.HibernationActive: state{
			next: func(_ core.Servant, node *core.Node, _ enclave.Inspection) (bool, error) {
				if p.triggers.sleep.IsMatch(node) {
					p.controls.Sleep.Dispatch()(p.triggers.sleep.Description())
					p.transition(enums.HibernationRetired)

					if p.ho.Behaviour.InclusiveSleep {
						return true, nil
					}

					return false, nil
				}

				return true, nil
			},
		},

		enums.HibernationRetired: state{
			next: func(_ core.Servant, _ *core.Node, _ enclave.Inspection) (bool, error) {
				return false, fs.SkipAll
			},
		},
	}
}
