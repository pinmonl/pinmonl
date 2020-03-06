package sharing

import (
	"github.com/pinmonl/pinmonl/handler/api/monl"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkg/markdown"
)

// Resp defines the response body of Share.
func Resp(m model.Share, u model.User, mustTags []model.Tag) map[string]interface{} {
	resp := map[string]interface{}{
		"id":          m.ID,
		"name":        m.Name,
		"description": m.Description,
		"readme":      markdown.SafeHTML(m.Readme),
		"imageId":     m.ImageID,
		"owner":       UserResp(u),
	}

	rts := make([]interface{}, len(mustTags))
	for i, t := range mustTags {
		rts[i] = TagResp(t)
	}
	resp["mustTags"] = rts

	return resp
}

// UserResp defines the response body of User.
func UserResp(u model.User) map[string]interface{} {
	return map[string]interface{}{
		"login": u.Login,
		"name":  u.Name,
	}
}

// TagResp defines the response body of Tag.
func TagResp(m model.Tag) map[string]interface{} {
	return map[string]interface{}{
		"id":       m.ID,
		"name":     m.Name,
		"parentId": m.ParentID,
		"sort":     m.Sort,
	}
}

// PinlResp defines the response body of Pinl.
func PinlResp(m model.Pinl, tags []model.Tag) map[string]interface{} {
	return map[string]interface{}{
		"id":          m.ID,
		"title":       m.Title,
		"description": m.Description,
		"imageId":     m.ImageID,
		"url":         m.URL,
		"tags":        (model.TagList)(tags).PluckName(),
	}
}

// PinlDetailResp defines the response body of Pinl with detail.
func PinlDetailResp(m model.Pinl, tags []model.Tag, pkgs []model.Pkg, stats []model.Stat) map[string]interface{} {
	resp := PinlResp(m, tags)

	rps := make([]interface{}, len(pkgs))
	for i, p := range pkgs {
		pss := (model.StatList)(stats).FindPkg(p)
		rps[i] = monl.PkgResp(p, pss)
	}
	resp["pkgs"] = rps

	return resp
}
