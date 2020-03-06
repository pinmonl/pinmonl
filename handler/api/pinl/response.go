package pinl

import (
	"github.com/pinmonl/pinmonl/handler/api/monl"
	"github.com/pinmonl/pinmonl/model"
)

// Resp returns basic data of pinl.
func Resp(m model.Pinl, tags []model.Tag) map[string]interface{} {
	resp := map[string]interface{}{
		"id":          m.ID,
		"url":         m.URL,
		"title":       m.Title,
		"description": m.Description,
		"imageId":     m.ImageID,
		"createdAt":   m.CreatedAt,
		"updatedAt":   m.UpdatedAt,
		"tags":        (model.TagList)(tags).PluckName(),
	}
	return resp
}

// DetailResp shows more information on top of Resp.
func DetailResp(m model.Pinl, tags []model.Tag, pkgs []model.Pkg, stats []model.Stat) map[string]interface{} {
	resp := Resp(m, tags)
	resp["readme"] = m.Readme

	rps := make([]interface{}, len(pkgs))
	for i, p := range pkgs {
		pss := (model.StatList)(stats).FindPkg(p)
		rps[i] = monl.PkgResp(p, pss)
	}
	resp["pkgs"] = rps

	return resp
}
