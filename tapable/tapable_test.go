package tapable_test

import (
	"io/fs"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	"github.com/onsi/ginkgo/v2/dsl/decorators"
	. "github.com/onsi/gomega" //nolint:revive // ok

	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/tapable"
)

type (
	ReadDirectoryHookFunc func(dirname string) ([]fs.DirEntry, error)
	ReadDirectoryHook     tapable.ActivityL[ReadDirectoryHookFunc]

	ReaderHost struct {
		Name string                // init with "default"
		Hook ReadDirectoryHookFunc // init with default func
	}

	Component struct {
		Reader ReaderHost
	}

	ExternalReaderClient struct{}
)

func (c *ExternalReaderClient) Init(from *Component) {
	def, err := from.Reader.Tap("external reader client", func(_ string) ([]fs.DirEntry, error) {
		return []fs.DirEntry{}, nil
	})

	_, _ = def, err
}

func (r *ReaderHost) Tap(name string, fn ReadDirectoryHookFunc) (ReadDirectoryHookFunc, error) {
	previous := r.Hook
	r.Hook = fn
	r.Name = name

	return previous, nil
}

func (c *Component) DoWork() {

}

var _ = Describe("Legacy Tapable", func() {
	Context("foo", func() {
		It("🧪 should: ", func() {
			// !!! this should be using Role
			//
			component := &Component{
				Reader: ReaderHost{
					Name: "default",
					Hook: func(_ string) ([]fs.DirEntry, error) {
						return []fs.DirEntry{}, nil
					},
				},
			}

			external := ExternalReaderClient{}
			external.Init(component)

			Expect(1).To(Equal(1))
		})
	})

	Context("scenarios", func() {
		When("client taps component's core action with hook", func() {

		})

		When("client augments component's core action with hook", func() {
			// client need to call default functionality

		})

		When("client subscribes to component's life-cycle event", func() {

		})

		Context("Component exposes multiple hooks", func() {

		})
	})
})

type (
	Roles struct {
		ReadDirectory tapable.WithDefault[ReadDirectoryHookFunc, enums.Role]
		Container     *tapable.Container[enums.Role]
	}

	Options struct {
		Hooks Roles
	}
)

var _ = Describe("Tapable", decorators.Label("use-case"), func() {
	Context("foo", func() {
		When("bar", func() {
			It("🧪 should: ", func() {
				// This could be exposed to the client as a WithXXX option,
				// or the client could perform the Tap manually themselves.
				// Container is an internal affair, it should be created and
				// used internally only.
				container := tapable.NewContainer[enums.Role]()

				// perhaps we make the tapping mechanism internal only and if
				// we do it this way, it doesn't matter if the container is passed
				// into the hook. We rely on WithXXX options to setup the hook,
				// and we do the tap internally on the client's behalf.
				//
				// WithReadDirectoryHook ==> options.Hooks.ReadDirectory.Tap(...)
				// This is with the caveat that there should be a separation of
				// the options which can be set directly be the user via With commands
				// and another abstraction which contains functional settings (ie
				// non persistable items; probably contains the original options
				// as a member).
				//
				// Actually, we can define a 'binder' which represents the options
				// used at runtime. IE, the user selects options. Some of the With commands
				// may set non persistable items, but these will go straight into
				// the binder. The binder will contain options as a member. There
				// will be translation from options to binder. Most of the internal
				// functionality will be dependent on the binder rather than the
				// options. Availability of the binder should become a life cycle event
				// during bootstrapping.
				//
				options := &Options{
					Hooks: Roles{
						ReadDirectory: tapable.WithDefault[ReadDirectoryHookFunc, enums.Role]{
							Default: func(_ string) ([]fs.DirEntry, error) {
								return []fs.DirEntry{}, nil
							},
							Container: container,
						},
						Container: container,
					},
				}

				_, _ = options.Hooks.ReadDirectory.Tap(
					"client",
					enums.RoleDirectoryReader,
					func(_ string) ([]fs.DirEntry, error) {
						return []fs.DirEntry{}, nil
					},
				)
			})
		})
	})
})
