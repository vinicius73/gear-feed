package stages

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	_ "golang.org/x/image/webp"

	"github.com/cenkalti/dominantcolor"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/muesli/gamut"
	"github.com/vinicius73/gamer-feed/assets/fonts"
)

const (
	boxMargin       = 15.0
	minFontSize     = 50.0
	footerSize      = 35.0
	fontLineSpacing = 1.2
	startFontSize   = 100.0
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
	dc     *gg.Context
	fonts  Fonts
	Colors CoverColors
	width  int
	height int
}

type drawPipe func(source FetchResult) error

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

func NewDraw(width, height int) (*Draw, error) {
	dc := gg.NewContext(width, height)

	ttFontTitle, err := fonts.UbuntuMonoBold()
	if err != nil {
		return nil, err
	}

	ttFontDescription, err := fonts.FiraMonoRegular()
	if err != nil {
		return nil, err
	}

	return &Draw{
		dc:     dc,
		width:  width,
		height: height,
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

func (d *Draw) Draw(source FetchResult) error {
	if err := d.DrawBase(source); err != nil {
		return err
	}

	if err := d.DrawOver(source); err != nil {
		return err
	}

	return nil
}

func (d *Draw) DrawBase(source FetchResult) error {
	pipes := []drawPipe{
		d.SetImage,
	}

	for _, pipe := range pipes {
		if err := pipe(source); err != nil {
			return err
		}
	}

	return nil
}

func (d *Draw) Reset() {
	d.dc.Clear()
	d.dc.SetRGBA(1, 1, 1, 0)
}

func (d *Draw) DrawOver(source FetchResult) error {
	pipes := []drawPipe{
		d.SetBackground,
		d.SetText,
	}

	for _, pipe := range pipes {
		if err := pipe(source); err != nil {
			return err
		}
	}

	return nil
}

func (d *Draw) SetBackground(_ FetchResult) error {
	x := boxMargin
	y := boxMargin

	//nolint:gomnd
	w := d.dc.Width() - (2.0 * boxMargin)
	//nolint:gomnd
	h := d.dc.Height() - (2.0 * boxMargin)

	box := gg.NewContext(w, h)
	box.SetColor(d.Colors.Box)
	box.DrawRectangle(0, 0, float64(w), float64(h))
	box.Fill()

	d.dc.DrawImage(box.Image(), int(x), int(y))

	return nil
}

func (d *Draw) SetImage(source FetchResult) error {
	tmpFile, err := os.CreateTemp(os.TempDir(), "fetch-*--"+source.ImageName())
	if err != nil {
		return err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	err = source.FetchImage(tmpFile)

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

	img = imaging.Fill(img, d.width, d.height, imaging.Center, imaging.Lanczos)

	d.dc.DrawImage(img, 0, 0)

	d.detectColor()

	return nil
}

func (d *Draw) SetText(source FetchResult) error {
	_, titleHeight, err := d.addTitleText(source.Title)
	if err != nil {
		return err
	}

	if err = d.addFooter(source); err != nil {
		return err
	}

	if err = d.addDescription(titleHeight, source.Text); err != nil {
		return err
	}

	return nil
}

func (d *Draw) Write(target io.Writer) error {
	return d.dc.EncodePNG(target)
}

func (d *Draw) detectColor() {
	d.Colors = NewCoverColors(d.dc.Image())
}

func (d *Draw) addTitleText(text string) (textWidth, textHeight float64, err error) {
	dc := d.dc

	W := dc.Width()
	H := dc.Height()
	P := boxMargin * 1.2

	yPad := P

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
	//nolint:gomnd
	dc.DrawStringWrapped(text, P+1, yPad+1, 0, 0, maxWidth, fontLineSpacing, gg.AlignLeft)
	dc.SetColor(d.Colors.Text)
	dc.DrawStringWrapped(text, P, yPad, 0, 0, maxWidth, fontLineSpacing, gg.AlignLeft)

	textHeight += yPad

	return
}

func (d *Draw) addFooter(source FetchResult) error {
	text := fmt.Sprintf("by %s", source.SiteName)

	if len(source.SiteName) == 0 {
		text = source.DomainName
	}

	dc := d.dc

	W := dc.Width()

	P := boxMargin * 2

	maxWidth := float64(W) - (P * 2)

	yPad := float64(dc.Height()) - (P * 2)

	dc.SetFontFace(truetype.NewFace(d.fonts.Footer, &truetype.Options{
		Size: footerSize,
	}))

	dc.SetColor(d.Colors.Text)
	//nolint:gomnd
	dc.DrawStringWrapped(text, P, yPad, 0, 0, maxWidth, fontLineSpacing, gg.AlignLeft)

	return nil
}

func (d *Draw) addDescription(paddingTop float64, text string) error {
	if len(text) == 0 {
		return nil
	}

	dc := d.dc

	W := dc.Width()

	P := boxMargin * 2

	yPad := paddingTop + (P * 4)

	maxWidth := float64(W) - (P * 2)

	dc.SetFontFace(truetype.NewFace(d.fonts.Description, &truetype.Options{
		Size: footerSize * 2,
	}))

	dc.SetColor(d.Colors.Main)
	dc.DrawStringWrapped(text, P, yPad+1, 0, 0, maxWidth, fontLineSpacing, gg.AlignLeft)
	dc.SetColor(d.Colors.Text)
	dc.DrawStringWrapped(text, P, yPad, 0, 0, maxWidth, fontLineSpacing, gg.AlignLeft)

	return nil
}
