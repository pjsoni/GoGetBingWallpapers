package main

import (
	"testing"
)

func TestContains(t *testing.T) {
	images := []Image{
		{URL: "https://www.example.com/image1.jpg", Date: "20230101"},
		{URL: "https://www.example.com/image2.jpg", Date: "20230102"},
	}
	image := Image{URL: "https://www.example.com/image1.jpg", Date: "20230101"}

	if !contains(images, image) {
		t.Errorf("Expected contains to return true, but got false")
	}

	image = Image{URL: "https://www.example.com/image3.jpg", Date: "20230103"}
	if contains(images, image) {
		t.Errorf("Expected contains to return false, but got true")
	}
}

func TestFormatDate(t *testing.T) {
	date := "20230101"
	expected := "2023-01-01"
	formattedDate := formatDate(date)
	if formattedDate != expected {
		t.Errorf("Expected formatDate to return %s, but got %s", expected, formattedDate)
	}

	date = "invalid"
	expected = "invalid"
	formattedDate = formatDate(date)
	if formattedDate != expected {
		t.Errorf("Expected formatDate to return %s, but got %s", expected, formattedDate)
	}
}
