package billscan

import "github.com/okcredit/go-common/errors"

var (
	ErrLinkNotFound          = errors.From(404, "link_not_found")
	ErrLinkExists            = errors.From(409, "link_exists")
	ErrInvalidDestinationUrl = errors.Invalid("destination_url")
)
