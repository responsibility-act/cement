package utils

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	RegPhone = regexp.MustCompile(`^1[3|4|5|7|8]\d{9}$`)
)

// copy from url.ParseQuery
func ParseMatrixPath(query string) (m url.Values, err error) {
	m = make(url.Values)
	err = parseMatrixPath(m, query)
	return
}

// copy from url.ParseQuery
func parseMatrixPath(m url.Values, query string) (err error) {
	for query != "" {
		key := query
		if i := strings.IndexAny(key, ";"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		//		value, err1 = url.QueryUnescape(value)
		//		if err1 != nil {
		//			if err == nil {
		//				err = err1
		//			}
		//			continue
		//		}
		m[key] = append(m[key], value)
	}
	return err
}
