package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type salePrice struct {
	amount uint
	date   time.Time
	tenure string
}

type house struct {
	address     string
	numBedrooms uint
	houseType   string
	sales       []salePrice
}

func main() {
	res := ReadDownloadedRightMoveHtml("rm-pg1.html")
	res = append(res, ReadDownloadedRightMoveHtml("rm-pg2.html")...)
	res = append(res, ReadDownloadedRightMoveHtml("rm-pg3.html")...)
	res = append(res, ReadDownloadedRightMoveHtml("rm-pg4.html")...)
	res = append(res, ReadDownloadedRightMoveHtml("rm-pg5.html")...)
	res = append(res, ReadDownloadedRightMoveHtml("rm-pg6.html")...)
	res = append(res, ReadDownloadedRightMoveHtml("rm-pg7.html")...)
	res = append(res, ReadDownloadedRightMoveHtml("rm-pg8.html")...)
	res = append(res, ReadDownloadedRightMoveHtml("rm-pg9.html")...)
	fmt.Println("There are " + strconv.FormatInt(int64(len(res)), 10) + " results")
	fmt.Println(res)
}

func ReadDownloadedRightMoveHtml(filename string) []house {
	var results []house

	body, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	htm, _ := doc.Html()
	_ = os.WriteFile("test.html", []byte(htm), 0644)

	// Find the review items
	doc.Find("body").Find("div").
		Each(func(i int, s *goquery.Selection) {

			a, exists := s.Attr("class")
			if !exists || a != "main-content" {
				return
			}

			s.Find(".propertyCard").
				Each(func(i int, s *goquery.Selection) {
					s.Children().
						Each(func(i int, s *goquery.Selection) {
							a, exists := s.Attr("class")
							if !exists || a != "propertyCard-content" {
								return
							}

							var h house

							s.Children().
								Each(func(i int, s *goquery.Selection) {
									a, exists := s.Attr("class")
									if !exists {
										return
									}

									// Address
									if a == "title clickable" {
										h.address = s.Text()
									}
									// Num bedrooms
									if a == "subTitle bedrooms" {
										nb := s.Text()
										nbb := strings.Split(nb, ",")
										nbb[0] = strings.Replace(nbb[0], " bed", "", 1)
										nbn, err := strconv.ParseUint(strings.TrimSpace(nbb[0]), 10, 64)
										if err != nil {
											panic(err)
										}

										h.numBedrooms = uint(nbn)
										h.houseType = strings.TrimSpace(nbb[1])
									}

									if a == "transaction-table-container" {
										var sp salePrice

										s.Find("tr").
											Each(func(i int, s *goquery.Selection) {
												s.Children().
													Each(func(i int, s *goquery.Selection) {
														a, exists := s.Attr("class")
														if !exists {
															return
														}

														// Price
														if a == "price" {
															p := s.Text()
															p = strings.Replace(p, "£", "", 1)
															p = strings.Replace(p, ",", "", -1)
															pa, err := strconv.ParseUint(p, 10, 64)
															if err != nil {
																panic(err)
															}
															sp.amount = uint(pa)
														}
														// date-sold
														if a == "date-sold" {
															ds := s.Text()
															ds = strings.Replace(ds, "(New Build)", "", 1)
															parseTime, err := time.Parse("2 Jan 2006", ds)
															if err != nil {
																panic(err)
															}
															sp.date = parseTime
														}
														// Freehold or leasehold
														if a == "table-extra tenure" {
															sp.tenure = s.Text()
														}
													})
											})

										h.sales = append(h.sales, sp)
										sp = salePrice{}
									}
								})

							results = append(results, h)
							h = house{}
						})
				})
		})

	return results
}
