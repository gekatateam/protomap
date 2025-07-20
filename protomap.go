package protomap

import (
	"context"
	"errors"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
)

var (
	ErrNoSuchFile    = errors.New("no such file")
	ErrNoSuchMessage = errors.New("no such message in descriptor")
)

type Mapper struct {
	r linker.Resolver
}

func NewMapper(compiler *protocompile.Compiler, files ...string) (*Mapper, error) {
	if compiler == nil {
		compiler = &protocompile.Compiler{
			Resolver: &protocompile.SourceResolver{},
		}
	}

	f, err := compiler.Compile(context.Background(), files...)
	if err != nil {
		return nil, err
	}

	return &Mapper{r: f.AsResolver()}, nil
}
