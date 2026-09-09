package lists

import (
	"time"

	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

type List struct {
	id            ListID
	name          string
	description   string
	private       bool
	ownerID       OwnerID
	memberCount   int
	followerCount int
	createdAt     time.Time
}

type ListProps struct {
	ID            string
	Name          string `validate:"omitempty,max=25"`
	Description   string `validate:"omitempty,max=100"`
	Private       bool
	OwnerID       string
	MemberCount   int `validate:"gte=0"`
	FollowerCount int `validate:"gte=0"`
	CreatedAt     time.Time
}

func NewList(props ListProps) (List, error) {
	id, err := NewListID(props.ID)
	if err != nil {
		return List{}, err
	}

	if err := validator.Validate(props); err != nil {
		return List{}, draftError(err)
	}

	list := List{
		id:            id,
		name:          props.Name,
		description:   props.Description,
		private:       props.Private,
		memberCount:   props.MemberCount,
		followerCount: props.FollowerCount,
		createdAt:     props.CreatedAt,
	}

	if props.OwnerID != "" {
		ownerID, err := NewOwnerID(props.OwnerID)
		if err != nil {
			return List{}, err
		}
		list.ownerID = ownerID
	}

	return list, nil
}

type ListView struct {
	ID            string
	Name          string
	Description   string
	Private       bool
	OwnerID       string
	MemberCount   int
	FollowerCount int
	CreatedAt     time.Time
	URL           string
}

func (l List) View() ListView {
	return ListView{
		ID:            l.id.String(),
		Name:          l.name,
		Description:   l.description,
		Private:       l.private,
		OwnerID:       l.ownerID.String(),
		MemberCount:   l.memberCount,
		FollowerCount: l.followerCount,
		CreatedAt:     l.createdAt,
		URL:           "https://x.com/i/lists/" + l.id.String(),
	}
}

func (l List) ID() ListID {
	return l.id
}

type Lists []List

func (l Lists) Len() int {
	return len(l)
}

func (l Lists) IsEmpty() bool {
	return len(l) == 0
}

func (l Lists) Views() []ListView {
	views := make([]ListView, 0, len(l))
	for _, list := range l {
		views = append(views, list.View())
	}

	return views
}

func (l Lists) Limit(n int) Lists {
	if n <= 0 || len(l) == 0 {
		return Lists{}
	}

	result := make(Lists, min(len(l), n))
	copy(result, l)

	return result
}

func (l Lists) Dedupe() Lists {
	seen := make(map[string]struct{}, len(l))
	result := make(Lists, 0, len(l))

	for _, list := range l {
		key := list.id.String()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, list)
	}

	return result
}
