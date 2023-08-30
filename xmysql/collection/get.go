// Copyright (c) 2023, Geert JM Vanderkelen

package collection

type GetOptions struct {
	ValidateExistence bool
}

type GetOption func(opts *GetOptions)

func NewGetOptions(opts ...GetOption) *GetOptions {
	options := &GetOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func GetValidateExistence() GetOption {
	return func(opts *GetOptions) {
		opts.ValidateExistence = true
	}
}
