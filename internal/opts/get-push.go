package opts

import (
	"github.com/snivilised/agenor/pref"
)

// Get returns a new instance of the options struct and a new instance of the
// binder struct, with the provided settings applied to the options struct. The
// binder struct is used to bind the options struct to the events emitted during
// navigation, allowing for dynamic updates to the options based on user input or
// other events. The Get function also returns an error if any of the provided
// settings are invalid or cannot be applied to the options struct.
func Get(settings ...pref.Option) (o *pref.Options, b *Binder, err error) {
	o = pref.DefaultOptions()

	binder := NewBinder()
	o.Events.Bind(binder.Controls)

	err = apply(o, settings...)

	return o, binder, err
}

func apply(o *pref.Options, settings ...pref.Option) (err error) {
	for _, option := range settings {
		if option != nil {
			err = option(o)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// Push returns a new instance of the binder struct, with the provided options
// struct's events bound to the binder's controls. This allows for dynamic updates
// to the options based on user input or other events emitted during navigation. The
// Push function does not return an error, as it assumes that the provided options
// struct is valid and can be used to bind the events without any issues.
func Push(o *pref.Options) *Binder {
	binder := NewBinder()
	o.Events.Bind(binder.Controls)

	return binder
}
