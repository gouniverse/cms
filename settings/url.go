package cms

// func (m UiManager) url(path string, params map[string]string) string {
// 	params["path"] = path

// 	url := m.endpoint + m.query(params)
// 	return url
// }

// func (m UiManager) query(queryData map[string]string) string {
// 	queryString := ""

// 	if len(queryData) > 0 {
// 		v := url.Values{}
// 		for key, value := range queryData {
// 			v.Set(key, value)
// 		}
// 		queryString += "?" + m.httpBuildQuery(v)
// 	}

// 	return queryString
// }

// func (m UiManager) httpBuildQuery(queryData url.Values) string {
// 	return queryData.Encode()
// }
