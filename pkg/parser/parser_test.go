package parser_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Raviraj2000/go-web-crawler/pkg/parser"
)

func TestParse(t *testing.T) {
	// Mock HTML content
	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
			<meta name="description" content="This is a test page">
		</head>
		<body>
			<a href="/relative-link">Relative Link</a>
			<a href="https://example.com/absolute-link">Absolute Link</a>
		</body>
		</html>
	`

	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Make HTTP request to mock server
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	data, links, err := parser.Parse(resp)
	if err != nil {
		t.Fatalf("Error parsing response: %v", err)
	}

	// Verify parsed data
	expectedTitle := "Test Page"
	expectedDescription := "This is a test page"
	if data.Title != expectedTitle {
		t.Errorf("Expected title %q, got %q", expectedTitle, data.Title)
	}
	if data.Description != expectedDescription {
		t.Errorf("Expected description %q, got %q", expectedDescription, data.Description)
	}

	// Verify extracted links
	expectedLinks := []string{
		server.URL + "/relative-link",
		"https://example.com/absolute-link",
	}
	if len(links) != len(expectedLinks) {
		t.Errorf("Expected %d links, got %d", len(expectedLinks), len(links))
	}
	for i, link := range links {
		if link != expectedLinks[i] {
			t.Errorf("Expected link %q, got %q", expectedLinks[i], link)
		}
	}
}

func TestParseNon200Response(t *testing.T) {
	// Mock HTTP server with non-200 response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Make HTTP request to mock server
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	_, _, err = parser.Parse(resp)
	if err == nil {
		t.Fatalf("Expected error for non-200 response, got nil")
	}
}

func TestNormalizeURL(t *testing.T) {
	baseURL, _ := http.NewRequest("GET", "https://example.com/base/", nil)
	tests := []struct {
		href     string
		expected string
	}{
		{"/relative", "https://example.com/relative"},
		{"https://other.com/absolute", "https://other.com/absolute"},
		{"   /trimmed ", "https://example.com/trimmed"},
	}

	for _, test := range tests {
		normalized := parser.NormalizeURL(baseURL.URL, test.href)
		if normalized != test.expected {
			t.Errorf("Expected %q, got %q", test.expected, normalized)
		}
	}
}
