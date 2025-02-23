package stages

import (
	"github.com/vinicius73/gear-feed/pkg/stories/drawer"
	"github.com/vinicius73/gear-feed/pkg/stories/fetcher"
	"github.com/vinicius73/gear-feed/pkg/stories/filetemplate"
	"github.com/vinicius73/gear-feed/pkg/support/apperrors"
)

var (
	ErrEmptyTarget         = apperrors.Business("target cannot be empty", "STAGES:EMPTY_TARGET")
	ErrFailtoCreateDir     = apperrors.System(nil, "fail to create dir", "STAGES:FAIL_TO_CREATE_DIR")
	ErrFailToParseTemplate = apperrors.System(nil, "fail to parse template", "STAGES:FAIL_TO_PARSE_TEMPLATE")
	ErrFailtoBuildFilename = apperrors.System(nil, "fail to build filename", "STAGES:FAIL_TO_BUILD_FILENAME")
)

type Footer drawer.Footer

type BuildStageOptions struct {
	Footer
	Source   fetcher.Result
	Template filetemplate.Template
}
