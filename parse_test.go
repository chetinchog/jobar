package main_test

import (
	"encoding/xml"
	"fmt"
	"os"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Elements []Element `xml:"elemento"`
	Items    []Element `xml:"item"`
}

type Element struct {
	Title string `xml:"title"`
}

func main() {
	data, _ := os.ReadFile("Works.xml")
	var rss RSS
	err := xml.Unmarshal(data, &rss)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Elements:", len(rss.Channel.Elements))
	fmt.Println("Items:", len(rss.Channel.Items))
}
