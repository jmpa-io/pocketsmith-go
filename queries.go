package pocketsmith

import "net/url"

// setupQueries largely exists to add default queries to the given queries
// *map[string]string. In this case, setupQueries returns url.Values that
// contain a default `per_page=100`. Blank string values are skipped so that
// optional filters don't get sent as empty query parameters (e.g. &start_date=).
func setupQueries(queries *map[string]string) url.Values {
	out := make(url.Values)

	// add any existing queries to output, skipping blank values.
	if queries != nil {
		for key, value := range *queries {
			if value != "" {
				out.Add(key, value)
			}
		}
	}

	// add default "per_page", if it's not already set.
	if _, ok := out["per_page"]; !ok {
		out["per_page"] = []string{"100"}
	}

	return out
}
