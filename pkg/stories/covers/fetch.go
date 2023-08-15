package covers

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

var ErrMissingImageURL = apperrors.Business("missing image url", "COVER:MISSING_IMAGE_URL")

const (
	requestTimeout = 10 * time.Second
)

type FetchResult struct {
	Title      string
	Text       string
	SiteName   string
	ImageURL   string
	DomainName string
	URL        string
}

func Fetch(ctx context.Context, source string) (FetchResult, error) {
	httpClient := &http.Client{
		Timeout: requestTimeout,
	}

	intent := opengraph.Intent{
		Context:    ctx,
		HTTPClient: httpClient,
	}

	ogp, err := opengraph.Fetch(source, intent)
	if err != nil {
		return FetchResult{}, err
	}

	image := findBestImage(ogp.Image)

	title := ogp.Title
	siteName := ogp.SiteName
	siteURL := ogp.URL

	if ogp.SiteName != "" {
		title = strings.TrimSuffix(title, "- "+siteName)
	}

	if ogp.URL == "" {
		siteURL = source
	}

	parsed, err := url.Parse(siteURL)
	if err != nil {
		return FetchResult{}, err
	}

	return FetchResult{
		Title:      title,
		Text:       ogp.Description,
		ImageURL:   image.URL,
		SiteName:   siteName,
		URL:        siteURL,
		DomainName: parsed.Hostname(),
	}, nil
}

func (f FetchResult) FetchImage(target io.Writer) error {
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

func (f FetchResult) ImageName() string {
	return filepath.Base(f.ImageURL)
}

func findBestImage(images []opengraph.Image) opengraph.Image {
	if len(images) == 0 {
		return opengraph.Image{}
	}

	var bestImage opengraph.Image
	found := false

	for _, image := range images {
		if image.Width > defaultWidth && image.Height > defaultHeight {
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
