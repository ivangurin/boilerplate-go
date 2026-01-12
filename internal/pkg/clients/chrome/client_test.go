package chrome_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/ledongthuc/pdf"
	"github.com/stretchr/testify/require"

	suite_provider "boilerplate/internal/pkg/suite/provider"
)

func TestPrintToPDF(t *testing.T) {
	sp, cleanup := suite_provider.NewProvider()
	defer cleanup()

	text := "Hello, World!"
	template := `<html><body><h1>%s</h1></body></html>`
	htmlContent := fmt.Sprintf(template, text)

	pdfBuffer, err := sp.GetChromeClient().PrintToPDF(sp.Context(), bytes.NewBufferString(htmlContent), false)
	require.NoError(t, err)
	require.NotNil(t, pdfBuffer)

	pdfData, err := io.ReadAll(pdfBuffer)
	require.NoError(t, err)
	require.NotEmpty(t, pdfData)

	pdfReader, err := pdf.NewReader(pdfBuffer, int64(len(pdfData)))
	require.NoError(t, err)

	var textContent strings.Builder
	for i := 1; i <= pdfReader.NumPage(); i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		require.NoError(t, err)

		_, err = textContent.WriteString(text)
		require.NoError(t, err)
	}

	require.Contains(t, textContent.String(), text)
}
