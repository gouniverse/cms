package cms

import "net/url"

func (m UiManager) url(path string, params map[string]string) string {
	params["path"] = path

	url := m.endpoint + query(params)
	return url
}

func query(queryData map[string]string) string {
	queryString := ""

	if len(queryData) > 0 {
		v := url.Values{}
		for key, value := range queryData {
			v.Set(key, value)
		}
		queryString += "?" + httpBuildQuery(v)
	}

	return queryString
}

func httpBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}
