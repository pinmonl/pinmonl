package pinmonl

import (
	"time"

	"github.com/pinmonl/pinmonl/model/field"
)

type (
	MonpkgKind    int
	SharetagKind  int
	StatKind      string
	Status        int
	StatValueType int
	UserRole      int
	UserStatus    int

	Pinl struct {
		ID          string     `json:"id"`
		UserID      string     `json:"userId"`
		MonlID      string     `json:"monlId"`
		URL         string     `json:"url"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		ImageID     string     `json:"imageId"`
		Status      Status     `json:"status"`
		CreatedAt   field.Time `json:"createdAt"`
		UpdatedAt   field.Time `json:"updatedAt"`

		TagNames *[]string `json:"tagNames,omitempty"`
	}

	Pkg struct {
		ID            string     `json:"id"`
		URL           string     `json:"url"`
		Provider      string     `json:"provider"`
		ProviderHost  string     `json:"providerHost"`
		ProviderURI   string     `json:"providerUri"`
		ProviderProto string     `json:"providerProto"`
		Title         string     `json:"title"`
		Description   string     `json:"description"`
		ImageID       string     `json:"imageId"`
		CustomUri     string     `json:"customUri"`
		CreatedAt     field.Time `json:"createdAt"`
		UpdatedAt     field.Time `json:"updatedAt"`
	}

	Monpkg struct {
		ID     string     `json:"id"`
		MonlID string     `json:"monlId"`
		PkgID  string     `json:"pkgId"`
		Kind   MonpkgKind `json:"kind"`
		Pkg    *Pkg       `json:"pkg"`
	}

	MonpkgListResponse struct {
		TotalCount int64     `json:"totalCount"`
		Page       int64     `json:"page"`
		PageSize   int64     `json:"pageSize"`
		Data       []*Monpkg `json:"data"`
	}

	Stat struct {
		ID          string        `json:"id"`
		PkgID       string        `json:"pkgId"`
		ParentID    string        `json:"parentId"`
		RecordedAt  field.Time    `json:"recordedAt"`
		Kind        StatKind      `json:"kind"`
		Name        string        `json:"name"`
		Value       string        `json:"value"`
		ValueType   StatValueType `json:"valueType"`
		Checksum    string        `json:"checksum"`
		Weight      int           `json:"weight"`
		IsLatest    bool          `json:"isLatest"`
		HasChildren bool          `json:"hasChildren"`

		Substats *[]*Stat `json:"substats,omitempty"`
	}

	StatListResponse struct {
		TotalCount int64   `json:"totalCount"`
		Page       int64   `json:"page"`
		PageSize   int64   `json:"pageSize"`
		Data       []*Stat `json:"data"`
	}

	User struct {
		ID        string     `json:"id"`
		Login     string     `json:"login"`
		Password  string     `json:"password"`
		Name      string     `json:"name"`
		ImageID   string     `json:"imageId"`
		Hash      string     `json:"-"`
		Role      UserRole   `json:"role"`
		Status    UserStatus `json:"status"`
		LastSeen  field.Time `json:"lastSeen"`
		CreatedAt field.Time `json:"createdAt"`
		UpdatedAt field.Time `json:"updatedAt"`
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
		ID          string     `json:"id"`
		UserID      string     `json:"userId"`
		Slug        string     `json:"slug"`
		Name        string     `json:"name"`
		Description string     `json:"description"`
		ImageID     string     `json:"imageId"`
		Status      Status     `json:"status"`
		CreatedAt   field.Time `json:"createdAt"`
		UpdatedAt   field.Time `json:"updatedAt"`

		User *User `json:"user,omitempty"`
	}

	Sharetag struct {
		ID          string       `json:"id"`
		ShareID     string       `json:"shareId"`
		TagID       string       `json:"tagId"`
		Kind        SharetagKind `json:"kind"`
		ParentID    string       `json:"parentId"`
		Level       int          `json:"level"`
		Status      Status       `json:"status"`
		HasChildren bool         `json:"hasChildren"`

		Tag *Tag `json:"tag"`
	}

	Sharepin struct {
		ID      string `json:"id"`
		ShareID string `json:"shareId"`
		PinlID  string `json:"pinlId"`
		Status  Status `json:"status"`

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
