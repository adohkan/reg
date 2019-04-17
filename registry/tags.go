package registry

import (
	"net/url"

	"github.com/peterhellberg/link"
)

type tagsResponse struct {
	Tags []string `json:"tags"`
}

// Tags returns the tags for a specific repository.
func (r *Registry) Tags(repository string) ([]string, error) {
	uri := r.url("/v2/%s/tags/list", repository)
	r.Logf("registry.tags url=%s repository=%s", uri, repository)

	var response tagsResponse
	h, err := r.getJSON(uri, &response)
	if err != nil {
		return nil, err
	}

	for _, l := range link.ParseHeader(h) {
		if l.Rel == "next" {
			unescaped, _ := url.QueryUnescape(l.URI)
			tags, err := r.Tags(unescaped)
			if err != nil {
				return nil, err
			}
			response.Tags = append(response.Tags, tags...)
		}
	}

	return response.Tags, nil
}
