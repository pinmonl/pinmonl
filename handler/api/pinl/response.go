package pinl

import "github.com/pinmonl/pinmonl/model"

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
	}
	if tags != nil {
		resp["tags"] = (model.TagList)(tags).PluckName()
	}
	return resp
}

// DetailResp shows more information on top of Resp.
func DetailResp(m model.Pinl, tags []model.Tag) map[string]interface{} {
	tags = append(make([]model.Tag, 0), tags...)
	resp := Resp(m, tags)
	resp["readme"] = m.Readme
	return resp
}
