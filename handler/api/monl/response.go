package monl

import "github.com/pinmonl/pinmonl/model"

// PkgResp defines the response body of Pkg.
func PkgResp(m model.Pkg, stats []model.Stat) map[string]interface{} {
	resp := map[string]interface{}{
		"id":          m.ID,
		"vendor":      m.Vendor,
		"vendorUri":   m.VendorURI,
		"title":       m.Title,
		"description": m.Description,
		"readme":      m.Readme,
		"imageId":     m.ImageID,
		"labels":      m.Labels,
	}

	rss := make([]map[string]interface{}, len(stats))
	for i, s := range stats {
		rss[i] = StatResp(s)
	}
	resp["stats"] = rss

	return resp
}

// StatResp defines the response body of Stat.
func StatResp(s model.Stat) map[string]interface{} {
	resp := map[string]interface{}{
		"kind":       s.Kind,
		"value":      s.Value,
		"recordedAt": s.RecordedAt,
		"manifest":   s.Manifest,
	}
	return resp
}
