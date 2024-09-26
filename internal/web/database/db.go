package database

import "github.com/nedpals/supabase-go"

type Params struct {
	Host   string
	Secret string
}

func New(params Params) *supabase.Client {
	sb := supabase.CreateClient(params.Host, params.Secret)

	return sb
}
