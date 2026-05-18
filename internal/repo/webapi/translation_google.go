package webapi

import (
	"context"
	"fmt"

	translator "github.com/Conight/go-googletrans"
	"github.com/bhcoder23/gin-clean-template/internal/entity"
)

// TranslationWebAPI is a demo adapter used to illustrate outbound integrations.
// Derived projects should treat it as a replaceable example rather than a
// production-ready translation provider.
type TranslationWebAPI struct {
	conf translator.Config
}

// New builds the demo translation adapter.
func New() *TranslationWebAPI {
	conf := translator.Config{
		UserAgent:   []string{"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1"},
		ServiceUrls: []string{"translate.google.com"},
	}

	return &TranslationWebAPI{
		conf: conf,
	}
}

// Translate -.
func (t *TranslationWebAPI) Translate(ctx context.Context, translation entity.Translation) (entity.Translation, error) {
	if err := ctx.Err(); err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationWebAPI - Translate - ctx: %w", err)
	}

	trans := translator.New(t.conf)

	result, err := trans.Translate(translation.Original, translation.Source, translation.Destination)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationWebAPI - Translate - trans.Translate: %w", err)
	}

	translation.Translation = result.Text

	return translation, nil
}
