package hiber

import (
	"io/fs"

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/filtering"
	"github.com/snivilised/traverse/life"
	"github.com/snivilised/traverse/locale"
)

type simple struct {
	common
	states  hibernateStates
	current state
}

func (p *simple) init(controls *life.Controls) error {
	p.states = p.create()
	p.controls = controls

	if p.common.ho.WakeAt != nil {
		filter, err := filtering.New(p.common.ho.WakeAt, p.common.fo)
		if err != nil {
			return err
		}

		p.common.triggers.wake = filter

		if p.common.triggers.sleep == nil {
			p.common.triggers.sleep = filtering.NewProhibitiveTraverseFilter(
				&core.FilterDef{
					Description: li18ngo.Text(locale.ProhibitiveWordTemplData{}),
				},
			)
		}
	}

	if p.common.ho.SleepAt != nil {
		filter, err := filtering.New(p.common.ho.SleepAt, p.common.fo)
		if err != nil {
			return err
		}

		p.common.triggers.sleep = filter

		if p.common.triggers.wake == nil {
			p.common.triggers.wake = filtering.NewPermissiveTraverseFilter(
				&core.FilterDef{
					Description: li18ngo.Text(locale.PermissiveWordTemplData{}),
				},
			)
		}
	}

	p.transition(launch(p.common.ho))

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
				if p.common.triggers.wake.IsMatch(node) {
					p.controls.Wake.Dispatch()(p.common.triggers.wake.Description())
					p.transition(enums.HibernationActive)

					if p.common.ho.Behaviour.InclusiveWake {
						return true, nil
					}
				}

				return false, nil
			},
		},

		enums.HibernationActive: state{
			next: func(_ core.Servant, node *core.Node, _ enclave.Inspection) (bool, error) {
				if p.common.triggers.sleep.IsMatch(node) {
					p.controls.Sleep.Dispatch()(p.common.triggers.sleep.Description())
					p.transition(enums.HibernationRetired)

					if p.common.ho.Behaviour.InclusiveSleep {
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
