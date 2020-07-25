package pinmonl

import (
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
)

type (
	Pinl struct {
		ID          string       `json:"id"`
		UserID      string       `json:"userId"`
		MonlID      string       `json:"monlId"`
		URL         string       `json:"url"`
		Title       string       `json:"title"`
		Description string       `json:"description"`
		ImageID     string       `json:"imageId"`
		Status      model.Status `json:"status"`
		CreatedAt   field.Time   `json:"createdAt"`
		UpdatedAt   field.Time   `json:"updatedAt"`

		Tags []string `json:"tags"`
	}

	Pkg struct {
		ID            string     `json:"id"`
		URL           string     `json:"url"`
		Provider      string     `json:"provider"`
		ProviderHost  string     `json:"providerHost"`
		ProviderURI   string     `json:"providerUri"`
		ProviderProto string     `json:"providerProto"`
		CreatedAt     field.Time `json:"createdAt"`
		UpdatedAt     field.Time `json:"updatedAt"`
	}

	Stat struct {
		ID          string              `json:"id"`
		PkgID       string              `json:"pkgId"`
		ParentID    string              `json:"parentId"`
		RecordedAt  field.Time          `json:"recordedAt"`
		Kind        model.StatKind      `json:"kind"`
		Value       string              `json:"value"`
		ValueType   model.StatValueType `json:"valueType"`
		Checksum    string              `json:"checksum"`
		Weight      int                 `json:"weight"`
		IsLatest    bool                `json:"isLatest"`
		HasChildren bool                `json:"hasChildren"`

		Substats []*Stat `json:"substats"`
	}

	User struct {
		ID        string           `json:"id"`
		Login     string           `json:"login"`
		Password  string           `json:"password"`
		Name      string           `json:"name"`
		ImageID   string           `json:"imageId"`
		Hash      string           `json:"-"`
		Role      model.UserRole   `json:"role"`
		Status    model.UserStatus `json:"status"`
		LastSeen  field.Time       `json:"lastSeen"`
		CreatedAt field.Time       `json:"createdAt"`
		UpdatedAt field.Time       `json:"updatedAt"`
	}

	Tag struct {
		ID          string     `json:"id"`
		Name        string     `json:"name"`
		UserID      string     `json:"userId"`
		ParentID    string     `json:"parentId"`
		Level       int        `json:"level"`
		Color       string     `json:"color"`
		BgColor     string     `json:"bgColor"`
		HasChildren bool       `json:"hasChildren"`
		CreatedAt   field.Time `json:"createdAt"`
		UpdatedAt   field.Time `json:"updatedAt"`
	}

	Share struct {
		ID          string       `json:"id"`
		UserID      string       `json:"userId"`
		Slug        string       `json:"slug"`
		Name        string       `json:"name"`
		Description string       `json:"description"`
		ImageID     string       `json:"imageId"`
		Status      model.Status `json:"status"`
		CreatedAt   field.Time   `json:"createdAt"`
		UpdatedAt   field.Time   `json:"updatedAt"`

		User *User `json:"user,omitempty"`
	}

	Sharetag struct {
		ID          string             `json:"id"`
		ShareID     string             `json:"shareId"`
		TagID       string             `json:"tagId"`
		Kind        model.SharetagKind `json:"kind"`
		ParentID    string             `json:"parentId"`
		Level       int                `json:"level"`
		Status      model.Status       `json:"status"`
		HasChildren bool               `json:"hasChildren"`

		Tag *Tag `json:"tag"`
	}

	Sharepin struct {
		ID      string       `json:"id"`
		ShareID string       `json:"shareId"`
		PinlID  string       `json:"pinlId"`
		Status  model.Status `json:"status"`

		Pinl *Pinl `json:"pinl"`
	}

	ServerInfo struct {
		Version string `json:"version"`
	}

	Token struct {
		Token    string    `json:"token"`
		ExpireAt time.Time `json:"expireAt"`
	}
)
