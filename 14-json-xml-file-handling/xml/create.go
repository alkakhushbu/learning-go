package main

import (
	"encoding/xml"
	"fmt"
	"log"
)

// convert xml to go structs
//https://blog.kowalczyk.info/tools/xmltogo/

// Define structs for nested XML elements
type Chapter struct {
	Title string `xml:"title"`
	Text  string `xml:"text"`
}

type Note struct {
	ID      string `xml:"id,attr"`
	To      string `xml:"to"`
	From    string `xml:"from"`
	Heading string `xml:"heading"`
	Body    string `xml:"body"`
}

type Author struct {
	Gender    string `xml:"gender,attr"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
}

// Define the Book struct to model the XML data
// Refer notes at the end of this program to know about the xml tags

type Book struct {
	XMLName     xml.Name `xml:"book"`
	ID          string   `xml:"id,attr"`
	Title       string   `xml:"title"`
	Author      Author   `xml:"author"`
	Messages    []Note   `xml:"messages>note"`
	Description string   `xml:"description"`
	Content     Chapter  `xml:"content>chapter"`
}

//the content>chapter means that the XML element corresponding to the struct field is nested within another XML element.
//Path in XML:
//content is the parent XML element.
//chapter is the child XML element within the content element.

func main() {
	data := []byte(`
<book id="bk101">
    <title>Go Programming</title>
    <author gender="male">
        <first_name>John</first_name>
        <last_name>Doe</last_name>
    </author>
    <messages>
        <note id="501">
            <to>Tom</to>
            <from>Harry</from>
            <heading>Reminder</heading>
            <body>Party this weekend!</body>
        </note>
        <note id="502">
            <to>Harry</to>
            <from>Tom</from>
            <heading>Re: Reminder</heading>
            <body>I will think about it</body>
        </note>
    </messages>
    <description>A comprehensive guide to Go programming.</description>
    <content>
        <chapter>
            <title>Introduction</title>
            <text>This books is so nice.</text>
        </chapter>
    </content>
</book>
    `)

	var book Book
	err := xml.Unmarshal(data, &book)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Unmarshalled Book: %+v\n", book)

	// Marshal the struct to XML
	xmlData, err := xml.MarshalIndent(book, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling to XML: %v\n", err)
		return
	}

	// Convert the XML data to a string and output it
	fmt.Printf("XML Output:\n%s\n", xmlData)
}
