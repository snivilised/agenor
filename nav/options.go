package nav

import (
	"errors"
)

// nav defines the core navigators

// mediator to coordinate initialisation?

/*
// WithStandardStyle sets a TermRenderer's styles with a standard (builtin)
// style.
func WithStandardStyle(style string) TermRendererOption {
	return func(tr *TermRenderer) error {
		styles, err := getDefaultStyle(style)
		if err != nil {
			return err
		}
		tr.ansiOptions.Styles = *styles
		return nil
	}
}
*/

// We need to define and formalise the life-cycle of a navigation session from what
// has been learnt from extendio. Do we also need to also allow the client to pass in
// custom callbacks for different stages of the sequence. (These concerns should possibly
// be delegated to the cycle package).
//
// The navigator should have 2 methods:
// - Walk implements sequential navigation; cancellation occurs via a ctrl-c,
// but note that the interrupt process still uses a context, so it is valid
// to pass in a context to Walk (actually, this may not be totally correct, so
// re-check this; ie how does the interrupt work and how to handle; chan/context).
//
// - Run  implements concurrent navigation using rx observable

// TraverseOptions can be used to render markdown content, posing a depth of
// customization and styles to fit your needs.
type TraverseOptions struct {
	Persist Persistables
}

// Persistables represents the traverse options that can be persisted via marshalling
type Persistables struct {
}

// A TraverseOptionFunc sets an option on a TraverseOptions.
type TraverseOptionFunc func(*TraverseOptions) error

func Cuddles(n int) error {
	if n < 0 {
		return errors.New("negative")
	}

	a := n
	if a < 0 {
		return errors.New("negative")
	}
	b := n
	if b < 0 {
		return errors.New("negative")
	}
	return errors.New("positive")
}
