package main

import (
	"fmt"
	"log"
	"WebCrawler/Helpers"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	var Url string
	fmt.Print("Please Enter Your Scraping Url: ")
	_, err := fmt.Scanln(&Url)
	if err != nil {
		log.Fatal("failed to read url", err)
	}
	var res = Helpers.
}
