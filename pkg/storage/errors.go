package storage

import "github.com/vinicius73/gear-feed/pkg/support/apperrors"

var (
	ErrFailToMarshalData = apperrors.System(nil, "fail to marshal data", "FAIL_TO_MARSHAL_DATA")
	ErrFailToMarshalMeta = apperrors.System(nil, "fail to marshal data", "FAIL_TO_MARSHAL_META")
)
