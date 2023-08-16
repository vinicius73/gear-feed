package fetcher

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/otiai10/opengraph/v2"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
)

var ErrMissingImageURL = apperrors.Business("missing image url", "FETCHER:MISSING_IMAGE_URL")

const (
	requestTimeout = 10 * time.Second
)

type Options struct {
	SourceURL     string
	DefaultWidth  int
	DefaultHeight int
}

type Result struct {
	Title      string
	Text       string
	SiteName   string
	ImageURL   string
	DomainName string
	URL        string
}

func Fetch(ctx context.Context, opt Options) (Result, error) {
	httpClient := &http.Client{
		Timeout: requestTimeout,
	}

	intent := opengraph.Intent{
		Context:    ctx,
		HTTPClient: httpClient,
	}

	ogp, err := opengraph.Fetch(opt.SourceURL, intent)
	if err != nil {
		return Result{}, err
	}

	image := findBestImage(opt, ogp.Image)

	title := ogp.Title
	siteName := ogp.SiteName
	siteURL := ogp.URL

	if ogp.SiteName != "" {
		title = strings.TrimSuffix(title, "- "+siteName)
	}

	if ogp.URL == "" {
		siteURL = opt.SourceURL
	}

	parsed, err := url.Parse(siteURL)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Title:      title,
		Text:       ogp.Description,
		ImageURL:   image.URL,
		SiteName:   siteName,
		URL:        siteURL,
		DomainName: parsed.Hostname(),
	}, nil
}

func (f Result) FetchImage(target io.Writer) error {
	if f.ImageURL == "" {
		return ErrMissingImageURL
	}

	httpClient := &http.Client{
		Timeout: requestTimeout,
	}

	req, err := http.NewRequest("GET", f.ImageURL, nil)
	if err != nil {
		return err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(target, res.Body)

	return err
}

func (f Result) ImageName() string {
	return filepath.Base(f.ImageURL)
}

func findBestImage(opt Options, images []opengraph.Image) opengraph.Image {
	if len(images) == 0 {
		return opengraph.Image{}
	}

	var bestImage opengraph.Image
	found := false

	for _, image := range images {
		if image.Width > opt.DefaultWidth && image.Height > opt.DefaultHeight {
			if bestImage.Width == 0 || bestImage.Height == 0 {
				bestImage = image
				found = true
			} else if image.Width < bestImage.Width && image.Height < bestImage.Height {
				bestImage = image
				found = true
			}
		}
	}

	if !found {
		return images[0]
	}

	return bestImage
}
