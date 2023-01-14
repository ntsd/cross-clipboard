//go:build !darwin && !linux
// +build !darwin,!linux

package main

import "log"

func main() {
	log.Fatalf("this program is not work for Windows")
}
