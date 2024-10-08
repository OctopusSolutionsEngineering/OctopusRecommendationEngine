package excluder

import (
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"strings"
)

type DefaultExcluder struct {
}

func (e DefaultExcluder) IsResourceExcluded(resourceName string, excludeAll bool, excludeThese []string, excludeAllButThese []string) bool {
	if strings.TrimSpace(resourceName) == "" {
		return true
	}

	if excludeAll {
		return true
	}

	if excludeThese != nil && slices.Index(excludeThese, resourceName) != -1 {
		return true
	}

	if excludeAllButThese != nil && len(excludeAllButThese) != 0 {
		// Ignore any empty strings
		filteredList := lo.Filter(excludeAllButThese, func(item string, index int) bool {
			return strings.TrimSpace(item) != ""
		})

		if len(filteredList) != 0 && slices.Index(filteredList, resourceName) == -1 {
			return true
		}
	}

	return false
}
