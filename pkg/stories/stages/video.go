package stages

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// --
// ffmpeg -y \
// -loop 1 -t 1 -i ./001--base.png \
// -loop 1 -t 9 -i ./001--over.png \
// -filter_complex "
// 	[0]zoompan=z='min(max(zoom,pzoom)+0.0015,1.3)':d=5:x='iw/2-(iw/zoom/2)':y='ih/2-(ih/zoom/2)':s=1080x1920:fps=60[f0];
// 	[1]fade=d=0.2:t=in:alpha=1,setpts=PTS-STARTPTS+1/TB[f1];
// 	[0][f0]overlay[bg1];
// 	[bg1][f1]overlay,format=yuv420p[v]
// " -map "[v]" \
// -movflags +faststart res.mp4
// --

type BuildVideoOptions struct {
	Stage
}

func BuildVideo(opt BuildVideoOptions) (string, error) {
	inputs := []string{
		"-loglevel", "error",
		"-y",
		"-loop", "1", "-t", "1", "-i", opt.Background,
		"-loop", "1", "-t", "9", "-i", opt.Foreground,

		"-filter_complex",
		"[0]zoompan=z='min(max(zoom,pzoom)+0.0015,1.3)':d=5:x='iw/2-(iw/zoom/2)':y='ih/2-(ih/zoom/2)':s=1080x1920:fps=60[f0];[1]fade=d=0.2:t=in:alpha=1,setpts=PTS-STARTPTS+1/TB[f1];[0][f0]overlay[bg1];[bg1][f1]overlay,format=yuv420p[v]",
		"-map", "[v]",
		"-movflags", "+faststart",
		fmt.Sprintf("/Users/luiz.moreira/ghq/github.com/vinicius73/gamer-feed/outputs/%v.mp4", time.Now().Unix()),
	}

	cmd := exec.Command("ffmpeg", inputs...)
	cmd.Stdout = os.Stdout

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return "", nil
}
