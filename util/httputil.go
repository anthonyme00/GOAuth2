package util

import "fmt"

type UrlQuery map[string]string

func (query *UrlQuery) CreateQuery() string {
	querystring := "?"

	i := 0
	for k, v := range *query {
		i++
		if i < len(*query) {
			querystring += fmt.Sprintf("%s=%s&", k, v)
		} else {
			querystring += fmt.Sprintf("%s=%s", k, v)
		}
	}

	return querystring
}
