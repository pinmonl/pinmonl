package request

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
)

type ContextKey int

const (
	PaginatorCtxKey ContextKey = iota
	AuthedCtxKey
	UserCtxKey
	ShareCtxKey
	PinlCtxKey
	PkgCtxKey
	TagCtxKey
	SharetagCtxKey
	ImageCtxKey
)

func WithPaginator(ctx context.Context, p *Paginator) context.Context {
	return context.WithValue(ctx, PaginatorCtxKey, p)
}

func PaginatorFrom(ctx context.Context) *Paginator {
	p, ok := ctx.Value(PaginatorCtxKey).(*Paginator)
	if ok {
		return p
	}
	return nil
}

func WithAuthed(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, AuthedCtxKey, user)
}

func AuthedFrom(ctx context.Context) *model.User {
	user, ok := ctx.Value(AuthedCtxKey).(*model.User)
	if ok {
		return user
	}
	return nil
}

func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, UserCtxKey, user)
}

func UserFrom(ctx context.Context) *model.User {
	user, ok := ctx.Value(UserCtxKey).(*model.User)
	if ok {
		return user
	}
	return nil
}

func WithShare(ctx context.Context, share *model.Share) context.Context {
	return context.WithValue(ctx, ShareCtxKey, share)
}

func ShareFrom(ctx context.Context) *model.Share {
	share, ok := ctx.Value(ShareCtxKey).(*model.Share)
	if ok {
		return share
	}
	return nil
}

func WithPinl(ctx context.Context, pinl *model.Pinl) context.Context {
	return context.WithValue(ctx, PinlCtxKey, pinl)
}

func PinlFrom(ctx context.Context) *model.Pinl {
	pinl, ok := ctx.Value(PinlCtxKey).(*model.Pinl)
	if ok {
		return pinl
	}
	return nil
}

func WithPkg(ctx context.Context, pkg *model.Pkg) context.Context {
	return context.WithValue(ctx, PkgCtxKey, pkg)
}

func PkgFrom(ctx context.Context) *model.Pkg {
	pkg, ok := ctx.Value(PkgCtxKey).(*model.Pkg)
	if ok {
		return pkg
	}
	return nil
}

func WithTag(ctx context.Context, tag *model.Tag) context.Context {
	return context.WithValue(ctx, TagCtxKey, tag)
}

func TagFrom(ctx context.Context) *model.Tag {
	tag, ok := ctx.Value(TagCtxKey).(*model.Tag)
	if ok {
		return tag
	}
	return nil
}

func WithSharetag(ctx context.Context, sharetag *model.Sharetag) context.Context {
	return context.WithValue(ctx, SharetagCtxKey, sharetag)
}

func SharetagFrom(ctx context.Context) *model.Sharetag {
	sharetag, ok := ctx.Value(SharetagCtxKey).(*model.Sharetag)
	if ok {
		return sharetag
	}
	return nil
}

func WithImage(ctx context.Context, image *model.Image) context.Context {
	return context.WithValue(ctx, ImageCtxKey, image)
}

func ImageFrom(ctx context.Context) *model.Image {
	image, ok := ctx.Value(ImageCtxKey).(*model.Image)
	if ok {
		return image
	}
	return nil
}
