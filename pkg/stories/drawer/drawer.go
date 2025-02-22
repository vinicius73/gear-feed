//nolint:varnamelen,gomnd,exhaustruct,nakedret,nonamedreturns
package drawer

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"

	"github.com/cenkalti/dominantcolor"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/muesli/gamut"
	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/assets/fonts"
	"github.com/vinicius73/gamer-feed/pkg/stories/fetcher"
)

const (
	maxTextLength   = 350
	boxMargin       = 10.0
	innerBoxMargin  = boxMargin * 2
	minFontSize     = 50.0
	footerFontSize  = 35.0
	fontLineSpacing = 1.1
	startFontSize   = 100.0
	footerImageSize = 90.0
)

type CoverColors struct {
	Main   color.Color
	Box    color.Color
	Text   color.Color
	Shadow color.Color
}

type Fonts struct {
	Title       *truetype.Font
	Footer      *truetype.Font
	Description *truetype.Font
}

type Draw struct {
	DrawOptions
	dc     *gg.Context
	fonts  Fonts
	Colors CoverColors
}

type Footer struct {
	Image string `fig:"image" json:"image" yaml:"image"`
	Text  string `fig:"text"  json:"text"  yaml:"text"`
}

type DrawOptions struct {
	Footer Footer
	Width  int
	Height int
}

type DrawPipe func(ctx context.Context, source fetcher.Result) error

func NewCoverColors(im image.Image) CoverColors {
	main := dominantcolor.Find(im)
	R, G, B, _ := color.RGBAModel.Convert(main).RGBA()

	boxColor := color.RGBA{
		R: uint8(R >> 8),
		G: uint8(G >> 8),
		B: uint8(B >> 8),
		A: 204,
	}

	return CoverColors{
		Main:   main,
		Box:    boxColor,
		Text:   gamut.Contrast(main),
		Shadow: gamut.Complementary(main),
	}
}

func NewDraw(opts DrawOptions) (*Draw, error) {
	dc := gg.NewContext(opts.Width, opts.Height)

	ttFontTitle, err := fonts.UbuntuMonoBold()
	if err != nil {
		return nil, err
	}

	ttFontDescription, err := fonts.FiraMonoRegular()
	if err != nil {
		return nil, err
	}

	return &Draw{
		dc:          dc,
		DrawOptions: opts,
		fonts: Fonts{
			Title:       ttFontTitle,
			Description: ttFontDescription,
			Footer:      ttFontTitle,
		},
		Colors: CoverColors{
			Main:   color.White,
			Box:    color.Black,
			Text:   color.White,
			Shadow: color.Opaque,
		},
	}, nil
}

func (d *Draw) Draw(ctx context.Context, source fetcher.Result) error {
	if err := d.DrawBase(ctx, source); err != nil {
		return err
	}

	return d.DrawOver(ctx, source)
}

func (d *Draw) DrawBase(ctx context.Context, source fetcher.Result) error {
	pipes := []DrawPipe{
		d.SetImage,
	}

	for _, pipe := range pipes {
		if err := pipe(ctx, source); err != nil {
			return err
		}
	}

	return nil
}

func (d *Draw) DrawOver(ctx context.Context, source fetcher.Result) error {
	pipes := []DrawPipe{
		d.SetBackground,
		d.SetText,
	}

	for _, pipe := range pipes {
		if err := pipe(ctx, source); err != nil {
			return err
		}
	}

	return nil
}

func (d *Draw) SetBackground(_ context.Context, _ fetcher.Result) error {
	x := boxMargin
	y := boxMargin

	w := d.dc.Width() - (innerBoxMargin)

	h := d.dc.Height() - (innerBoxMargin)

	box := gg.NewContext(w, h)
	box.SetColor(d.Colors.Box)
	box.DrawRectangle(0, 0, float64(w), float64(h))
	box.Fill()

	d.dc.DrawImage(box.Image(), int(x), int(y))

	return nil
}

func (d *Draw) SetImage(ctx context.Context, source fetcher.Result) error {
	tmpFile, err := os.CreateTemp(os.TempDir(), "fetch-*--"+source.ImageName())
	if err != nil {
		return err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	err = source.FetchImage(ctx, tmpFile)
	if err != nil {
		return err
	}

	if _, err = tmpFile.Seek(0, 0); err != nil {
		return fmt.Errorf("error seeking file: %w", err)
	}

	img, _, err := image.Decode(tmpFile)
	if err != nil {
		return fmt.Errorf("error decoding image: %w", err)
	}

	img = imaging.Fill(img, d.Width, d.Height, imaging.Center, imaging.Lanczos)

	d.dc.DrawImage(img, 0, 0)

	d.detectColor()

	return nil
}

func (d *Draw) SetText(ctx context.Context, source fetcher.Result) error {
	if err := d.addHead(source); err != nil {
		return err
	}

	_, titleHeight := d.addTitleText(source)

	if err := d.addFooter(ctx, source); err != nil {
		return err
	}

	return d.addDescription(titleHeight, source.Text)
}

func (d *Draw) Write(target io.Writer) error {
	return d.dc.EncodePNG(target)
}

func (d *Draw) Reset() {
	d.dc.Clear()
	d.dc.SetRGBA(1, 1, 1, 0)
}

func (d *Draw) detectColor() {
	d.Colors = NewCoverColors(d.dc.Image())
}

func (d *Draw) addTitleText(source fetcher.Result) (textWidth, textHeight float64) {
	text := source.Title

	dc := d.dc

	W := dc.Width()
	H := dc.Height()
	P := innerBoxMargin

	yPad := P * (fontLineSpacing + 2)

	//nolint:gomnd
	maxWidth := float64(W) - (P * 2)
	//nolint:gomnd
	maxHeight := (float64(H) - (P * 2)) * 0.9

	fontSize := startFontSize

	updateFont := func() {
		dc.SetFontFace(truetype.NewFace(d.fonts.Title, &truetype.Options{
			Size: fontSize,
		}))
	}

	updateFont()

	for {
		if fontSize < minFontSize {
			break
		}

		updateFont()

		lines := dc.WordWrap(text, maxWidth)
		linesCount := len(lines)
		mls := ""

		for index, line := range lines {
			mls += line
			// last line
			if index != linesCount-1 {
				mls += "\n"
			}
		}

		textWidth, textHeight = dc.MeasureMultilineString(mls, fontLineSpacing)

		if textHeight < (maxHeight - (2 * P)) {
			break
		}

		fontSize -= (fontSize * 0.1)
	}

	dc.SetColor(d.Colors.Shadow)
	dc.DrawStringWrapped(text, P+1, yPad+1, 0, 0, maxWidth, fontLineSpacing, gg.AlignLeft)
	dc.SetColor(d.Colors.Text)
	dc.DrawStringWrapped(text, P, yPad, 0, 0, maxWidth, fontLineSpacing, gg.AlignLeft)

	textHeight += yPad

	return
}

func (d *Draw) addHead(source fetcher.Result) error {
	text := "by " + source.SiteName

	if len(source.SiteName) == 0 {
		text = source.DomainName
	}

	dc := d.dc

	W := dc.Width()
	P := innerBoxMargin * 1.5
	yPad := (innerBoxMargin + (fontLineSpacing + 1)) * 1.5

	maxWidth := float64(W) - (P * 2)

	colors := map[color.Color][2]float64{
		d.Colors.Main: {1.01, -4},
		// d.Colors.Shadow: {1, 0},
		d.Colors.Text: {0.99, 1},
	}

	for color, sizes := range colors {
		dc.SetColor(color)
		dc.SetFontFace(truetype.NewFace(d.fonts.Footer, &truetype.Options{
			Size: footerFontSize * sizes[0],
		}))

		dc.DrawStringWrapped(text, sizes[1]+P, sizes[1]+yPad, 0, 0, maxWidth, fontLineSpacing, gg.AlignLeft)
	}

	return nil
}

func (d *Draw) addFooter(ctx context.Context, _ fetcher.Result) error {
	logger := zerolog.Ctx(ctx)

	text := d.Footer.Text

	dc := d.dc

	W := dc.Width()

	P := innerBoxMargin * 1.5

	maxWidth := float64(W) - (P * 2)

	yPad := float64(dc.Height()) - (P * 4)
	xPad := P

	if d.Footer.hasImage() {
		xPad += footerImageSize + P
	}

	dc.SetFontFace(truetype.NewFace(d.fonts.Footer, &truetype.Options{
		Size: footerFontSize,
	}))

	dc.SetColor(d.Colors.Text)
	dc.DrawStringWrapped(text, xPad, yPad-(footerImageSize/2)+P, 0, 0, maxWidth, fontLineSpacing, gg.AlignLeft)

	if !d.Footer.hasImage() {
		logger.Debug().Any("footer", d.Footer).Msg("no image in footer")

		return nil
	}

	img, err := d.Footer.getImage()
	if err != nil {
		return err
	}

	dc.DrawImage(img, int(P), int(yPad-P))

	return nil
}

func (d *Draw) addDescription(paddingTop float64, text string) error {
	textLen := len(text)

	if textLen == 0 {
		return nil
	} else if textLen > maxTextLength {
		text = text[:maxTextLength] + "[…]"
	}

	dc := d.dc

	W := dc.Width()

	P := innerBoxMargin

	yPad := paddingTop + (P * 4)

	maxWidth := float64(W) - (P * 2)

	dc.SetFontFace(truetype.NewFace(d.fonts.Description, &truetype.Options{
		Size: footerFontSize * 2,
	}))

	dc.SetColor(d.Colors.Main)
	dc.DrawStringWrapped(text, P, yPad+1, 0, 0, maxWidth, fontLineSpacing, gg.AlignLeft)
	dc.SetColor(d.Colors.Text)
	dc.DrawStringWrapped(text, P, yPad, 0, 0, maxWidth, fontLineSpacing, gg.AlignLeft)

	return nil
}

func (c Footer) hasImage() bool {
	return len(c.Image) > 0
}

func (c Footer) getImage() (image.Image, error) {
	file, err := os.Open(c.Image)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return imaging.Fill(img, footerImageSize, footerImageSize, imaging.Center, imaging.Lanczos), nil
}
