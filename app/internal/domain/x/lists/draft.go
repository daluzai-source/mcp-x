package lists

import (
	"strings"
	"unicode/utf8"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

const (
	MaxNameLength        = 25
	MaxDescriptionLength = 100
)

type Draft struct {
	name        string
	description string
	private     bool
}

type DraftProps struct {
	Name        string `validate:"required,min=1,max=25"`
	Description string `validate:"omitempty,max=100"`
	Private     bool
}

var draftErrors = map[string]error{
	"Name.max":        x.ErrListNameTooLong,
	"Description.max": x.ErrListDescriptionTooLong,
}

func draftError(err error) error {
	key, ok := validator.Key(err)
	if !ok {
		return x.ErrInvalidRequest
	}

	if mapped, ok := draftErrors[key]; ok {
		return mapped
	}

	return x.ErrInvalidRequest
}

func NewDraft(props DraftProps) (Draft, error) {
	props.Name = strings.TrimSpace(props.Name)
	props.Description = strings.TrimSpace(props.Description)

	if err := validator.Validate(props); err != nil {
		return Draft{}, draftError(err)
	}

	return Draft{name: props.Name, description: props.Description, private: props.Private}, nil
}

func (d Draft) Name() string {
	return d.name
}

func (d Draft) Description() string {
	return d.description
}

func (d Draft) Private() bool {
	return d.private
}

type Patch struct {
	name        *string
	description *string
	private     *bool
}

type PatchProps struct {
	Name        *string
	Description *string
	Private     *bool
}

func NewPatch(props PatchProps) (Patch, error) {
	if props.Name == nil && props.Description == nil && props.Private == nil {
		return Patch{}, x.ErrInvalidRequest
	}

	patch := Patch{private: props.Private}

	if props.Name != nil {
		name := strings.TrimSpace(*props.Name)

		if name == "" {
			return Patch{}, x.ErrInvalidRequest
		}

		if utf8.RuneCountInString(name) > MaxNameLength {
			return Patch{}, x.ErrListNameTooLong
		}

		patch.name = &name
	}

	if props.Description != nil {
		description := strings.TrimSpace(*props.Description)

		if utf8.RuneCountInString(description) > MaxDescriptionLength {
			return Patch{}, x.ErrListDescriptionTooLong
		}

		patch.description = &description
	}

	return patch, nil
}

func (p Patch) Name() *string {
	return p.name
}

func (p Patch) Description() *string {
	return p.description
}

func (p Patch) Private() *bool {
	return p.private
}
