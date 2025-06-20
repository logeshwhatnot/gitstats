package main

import "flag"

func main() {
	var folder string
	var email string

	flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repos")
	flag.StringVar(&email, "email", "your@email.com", "your email to scan")
	flag.Parse()

	if folder != "" {
		scan(folder)
		return
	}

	stats(email)
}
