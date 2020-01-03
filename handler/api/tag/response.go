package tag

import "github.com/pinmonl/pinmonl/model"

// Resp returns basic data of tag.
func Resp(m model.Tag) map[string]interface{} {
	resp := map[string]interface{}{
		"id":       m.ID,
		"name":     m.Name,
		"parentId": m.ParentID,
		"sort":     m.Sort,
	}
	return resp
}
