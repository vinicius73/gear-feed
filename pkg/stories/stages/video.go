package stages

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gear-feed/pkg/support/apperrors"
)

var (
	ErrFailToCreateVideo    = apperrors.System(nil, "fail to create video", "STAGES:FAIL_TO_CREATE_VIDEO")
	ErrTargetVideoMustBeMp4 = apperrors.Business("target video must be a mp4 file", "STAGES:TARGET_VIDEO_MUST_BE_MP4")
)

// --
// ffmpeg -y \
// -loop 1 -t 1 -i ./001--base.png \
// -loop 1 -t 14 -i ./001--over.png \
// -filter_complex "
// 	[0]zoompan=z='min(max(zoom,pzoom)+0.0015,1.3)':d=5:x='iw/2-(iw/zoom/2)':y='ih/2-(ih/zoom/2)':s=1080x1920:fps=60[f0];
// 	[1]fade=d=0.2:t=in:alpha=1,setpts=PTS-STARTPTS+1/TB[f1];
// 	[0][f0]overlay[bg1];
// 	[bg1][f1]overlay,format=yuv420p[v]
// " -map "[v]" \
// -movflags +faststart res.mp4
// --

var ffmpegFilters = strings.Join([]string{
	"[0]zoompan=z='min(max(zoom,pzoom)+0.0015,1.3)':d=5:x='iw/2-(iw/zoom/2)':y='ih/2-(ih/zoom/2)':s=1080x1920:fps=60[f0]",
	"[1]fade=d=0.2:t=in:alpha=1,setpts=PTS-STARTPTS+1/TB[f1]",
	"[0][f0]overlay[bg1]",
	"[bg1][f1]overlay,format=yuv420p[v]",
}, ";")

type BuildVideoOptions struct {
	Stage
	Target string
}

func BuildVideo(ctx context.Context, opt BuildVideoOptions) (string, error) {
	if filepath.Ext(opt.Target) != ".mp4" {
		return "", ErrTargetVideoMustBeMp4
	}

	inputs := []string{
		"-loglevel", "warning",
		"-y",
		"-loop", "1", "-t", "1", "-i", opt.Background,
		"-loop", "1", "-t", "14", "-i", opt.Foreground,

		"-filter_complex",
		ffmpegFilters,
		"-map", "[v]",
		"-movflags", "+faststart",
		opt.Target,
	}

	cmd := exec.Command("ffmpeg", inputs...)

	logger := zerolog.Ctx(ctx).With().
		Str("stage", "build-video").
		Str("name", filepath.Base(opt.Target)).
		Logger()

	cmd.Stdout = logger.With().Str("out", "stdout").Logger()
	cmd.Stderr = logger.With().Str("out", "stderr").Logger()

	logger.Debug().Str("cmd", cmd.String()).Msg("executing ffmpeg")

	err := cmd.Run()
	if err != nil {
		return "", ErrFailToCreateVideo.Wrap(err)
	}

	logger.Info().Str("target", opt.Target).Msg("video created")

	return opt.Target, nil
}
