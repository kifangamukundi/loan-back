package parameters

import (
	"bytes"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func ToSlug(input string) string {
	slug := strings.ToLower(input)

	slug = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(slug, "")

	slug = regexp.MustCompile(`[\s_]+`).ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}

func TrimWhitespace(input string) string {
	return strings.TrimSpace(input)
}

func SanitizeHTML(input string) (string, error) {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = cleanHTML(&buf, doc)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func cleanHTML(w *bytes.Buffer, n *html.Node) error {
	allowedTags := map[string]bool{
		"p":      true,
		"a":      true,
		"b":      true,
		"i":      true,
		"strong": true,
		"em":     true,
		"ul":     true,
		"ol":     true,
		"li":     true,
		"br":     true,
		"img":    true,
		"h1":     true,
		"h2":     true,
		"h3":     true,
		"h4":     true,
		"h5":     true,
		"h6":     true,
	}

	if n.Type == html.ElementNode {
		if !allowedTags[n.Data] {
			return nil
		}
	}

	if err := html.Render(w, n); err != nil {
		return err
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := cleanHTML(w, c); err != nil {
			return err
		}
	}

	return nil
}

func SanitizeText(input string, isHTML bool) string {
	if isHTML {
		sanitized, err := SanitizeHTML(input)
		if err != nil {
			return input
		}
		return sanitized
	}
	reSpecialChars := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	input = reSpecialChars.ReplaceAllString(input, "")
	return strings.TrimSpace(input)
}
