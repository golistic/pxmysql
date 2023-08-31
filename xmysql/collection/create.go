// Copyright (c) 2023, Geert JM Vanderkelen

package collection

type CreateOptions struct {
	ReuseExisting bool
}

type CreateOption func(opts *CreateOptions)

func NewCreateOptions(opts ...CreateOption) *CreateOptions {
	options := &CreateOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func CreateReuseExisting() CreateOption {
	return func(opts *CreateOptions) {
		opts.ReuseExisting = true
	}
}
