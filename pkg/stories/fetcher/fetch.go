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
	"github.com/vinicius73/gamer-feed/pkg/support"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
)

var (
	ErrMissingImageURL = apperrors.Business("missing image url", "FETCHER:MISSING_IMAGE_URL")
	ErrFailToHash      = apperrors.System(nil, "fail to hash", "STAGES:FAIL_TO_HASH")
)

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
	Hash       string
}

func Fetch(ctx context.Context, opt Options) (Result, error) {
	//nolint:exhaustruct
	httpClient := &http.Client{
		Timeout: requestTimeout,
	}

	//nolint:exhaustruct
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
		title = strings.TrimSuffix(title, "| "+siteName)
	}

	if ogp.URL == "" {
		siteURL = opt.SourceURL
	}

	parsed, err := url.Parse(siteURL)
	if err != nil {
		return Result{}, err
	}

	hash, err := support.HashSHA256(opt.SourceURL)
	if err != nil {
		return Result{}, ErrFailToHash.Wrap(err)
	}

	return Result{
		Title:      strings.TrimSpace(title),
		Text:       strings.TrimSpace(ogp.Description),
		ImageURL:   image.URL,
		Hash:       hash,
		SiteName:   siteName,
		URL:        siteURL,
		DomainName: parsed.Hostname(),
	}, nil
}

func (f Result) FetchImage(ctx context.Context, target io.Writer) error {
	if f.ImageURL == "" {
		return ErrMissingImageURL
	}

	//nolint:exhaustruct
	httpClient := &http.Client{
		Timeout: requestTimeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.ImageURL, nil)
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
		//nolint:exhaustruct
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
