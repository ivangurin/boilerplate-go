package chrome

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Client interface {
	PrintToPDF(ctx context.Context, html *bytes.Buffer, landscape bool) (*bytes.Reader, error)
}

type client struct {
	url     string
	timeout int
}

func NewClient(host, port string, timeout int) Client {
	return &client{
		url:     fmt.Sprintf("ws://%s:%s", host, port),
		timeout: timeout,
	}
}

func (c *client) PrintToPDF(ctx context.Context, html *bytes.Buffer, landscape bool) (*bytes.Reader, error) {
	ctx, cancel := chromedp.NewRemoteAllocator(ctx, c.url)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Duration(c.timeout)*time.Second)
	defer cancel()

	var data []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return fmt.Errorf("get frame tree: %w", err)
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html.String()).Do(ctx)
		}),
		chromedp.WaitReady("body"),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			data, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithDisplayHeaderFooter(false).
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithLandscape(landscape).
				WithPreferCSSPageSize(false).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("pd print to pdf: %w", err)
	}

	return bytes.NewReader(data), nil
}
