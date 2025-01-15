package pocketsmith

import "net/url"

// setupQueries largely exists to add default queries to the given queries
// *map[string]string. In this case, setupQueries returns url.Values that
// contain a default `page_size=100`.
func setupQueries(queries *map[string]string) url.Values {
	out := make(url.Values)

	// add any existing queries to output.
	if queries != nil {
		for key, value := range *queries {
			out.Add(key, value)
		}
	}

	// add default "page_size", if it's not already set.
	if _, ok := out["page_size"]; !ok {
		out["page_size"] = []string{"100"}
	}

	return out
}
