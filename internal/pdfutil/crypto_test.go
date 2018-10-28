package pdfutil

import (
	"fmt"
	"testing"
)

func TestPathExpansiion(t *testing.T) {
	path, err := getPath("/", "~/ABC")
	fmt.Printf("Got path %s\n", path)
	if err != nil {
		t.Errorf("Could not get home-path %s", err)
	}

	path, err = getPath("/a/b", "./document.pdf")
	fmt.Printf("Got path %s\n", path)
	if err != nil {
		t.Errorf("Could not get home-path %s", err)
	}

	path, err = getPath("/a/b", "../document.pdf")
	fmt.Printf("Got path %s\n", path)
	if err != nil {
		t.Errorf("Could not get home-path %s", err)
	}
}
