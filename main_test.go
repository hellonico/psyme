package main

import (
	"testing"
)

func TestUser(t *testing.T) {
	img := mapForUser("Nico")
	writeImgToFile(img)
}
