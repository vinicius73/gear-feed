package fonts

import (
	"embed"

	"github.com/golang/freetype/truetype"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

//go:embed *.ttf
var ttfFS embed.FS

var (
	ubuntuMonoBold    = newFont("UbuntuMono-Bold.ttf")
	ubuntuMonoRegular = newFont("UbuntuMono-Regular.ttf")
	firaMonoBold      = newFont("FiraMono-Bold.ttf")
	firaMonoRegular   = newFont("FiraMono-Regular.ttf")
)

type font struct {
	file string
}

func newFont(file string) font {
	return font{
		file: file,
	}
}

func readFont(name string) (*truetype.Font, error) {
	file, err := ttfFS.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return truetype.Parse(file)
}

func (f *font) Read() (*truetype.Font, error) {
	return readFont(f.file)
}

func UbuntuMonoRegular() (*truetype.Font, error) {
	return ubuntuMonoRegular.Read()
}

func UbuntuMonoBold() (*truetype.Font, error) {
	return ubuntuMonoBold.Read()
}

func FiraMonoBold() (*truetype.Font, error) {
	return firaMonoBold.Read()
}

func FiraMonoRegular() (*truetype.Font, error) {
	return firaMonoRegular.Read()
}

func RandomBold() (*truetype.Font, error) {
	fonts := []font{
		ubuntuMonoBold,
		firaMonoBold,
	}

	font := support.Random[font](fonts)

	return font.Read()
}

func RandomRegular() (*truetype.Font, error) {
	fonts := []font{
		ubuntuMonoRegular,
		firaMonoRegular,
	}

	font := support.Random[font](fonts)

	return font.Read()
}
