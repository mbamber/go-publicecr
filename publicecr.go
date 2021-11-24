package publicecr

import (
	"bytes"
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type ECRReader struct {
	opts *ECRReaderOpts
}

func New(optFns ...func(*ECRReaderOpts)) *ECRReader {
	opts := DefaultECRReaderOpts
	for _, f := range optFns {
		f(opts)
	}
	return &ECRReader{opts}
}

// ListTags returns a slice of tags for a public ECR repo given by repoOwner and repoName
func (r *ECRReader) ListTags(ctx context.Context, repoOwner, repoName string) ([]*Image, error) {
	url := r.urlFromRepoOwnerAndName(repoOwner, repoName)
	var images []*Image

	cctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// get the contents of the ECR gallery for further processing
	// need to use a headless browser here as the page is rendered with JS
	var content string
	err := chromedp.Run(cctx,
		chromedp.Navigate(url), // Load the page
		// chromedp.OuterHTML("#root", &fullContent, chromedp.ByQueryAll),          // Dump the entire page for debugging
		chromedp.Click("[aria-controls$='tags-panel']", chromedp.ByQueryAll),    // Click the "tags" button
		chromedp.OuterHTML("[id$='tags-panel']", &content, chromedp.ByQueryAll)) // Dump the tag contents contents for use later
	if err != nil {
		return nil, err
	}
	cancel()
	// r.printf("full page content: %s", fullContent)
	r.printf("tags html content: %s", content)

	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(content))
	if err != nil {
		return nil, err
	}

	doc.Find(".awsui-table-row").Each(func(i int, s *goquery.Selection) {
		image := &Image{}
		for i, child := range s.Children().Nodes {
			switch i {
			case 0: // Image name
				image.Name = child.FirstChild.FirstChild.FirstChild.Data
			case 1: // Image type
				image.Type = child.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.Data
			case 2: // Date pushed
			case 3: // Image URI
			case 4: // Image size
			default:
				r.printf("ignoring unknown child at index %d: %+v", i, child)
			}
		}
		images = append(images, image)
	})

	return images, nil
}

// print prints to the logger if defined
func (r *ECRReader) print(v ...interface{}) {
	if r.opts.logger != nil {
		r.opts.logger.Print(v)
	}
}

// printf prints with a format string to the logger if defined
func (r *ECRReader) printf(format string, v ...interface{}) {
	if r.opts.logger != nil {
		r.opts.logger.Printf(format, v)
	}
}

// urlFromRepoOwnerAndName returns the full URL to the public ECR repo
func (r *ECRReader) urlFromRepoOwnerAndName(owner, name string) string {
	return fmt.Sprintf("%s%s/%s", r.opts.baseURL, owner, name)
}

type Image struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}
