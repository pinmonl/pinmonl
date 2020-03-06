package share

import "github.com/pinmonl/pinmonl/model"

// Resp returns basic data of share.
func Resp(m model.Share, mustTags []model.Tag) map[string]interface{} {
	resp := map[string]interface{}{
		"id":          m.ID,
		"name":        m.Name,
		"description": m.Description,
		"readme":      m.Readme,
		"mustTags":    (model.TagList)(mustTags).PluckName(),
	}
	return resp
}

// DetailResp returns tags relations on top of the basic response.
func DetailResp(m model.Share, mustTags []model.Tag, anyTags []model.Tag) map[string]interface{} {
	resp := Resp(m, mustTags)
	resp["anyTags"] = (model.TagList)(anyTags).PluckName()
	return resp
}
