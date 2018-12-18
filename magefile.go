// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Build build go command.
func Build() {
	sh.Run("go", "build")
}
