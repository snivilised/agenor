package types

type ContextExpiry interface {
	Expired() // ??? ctx context.Context, cancel context.CancelFunc
}
