package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var postcodeRegex = regexp.MustCompile("([A-Z]{1,2}[0-9][A-Z0-9]? ?[0-9][A-Z]{2}|GIR ?0A{2})")

type salePrice struct {
	Amount uint
	Date   time.Time
	Tenure string
}

type postcode struct {
	FirstPart  string
	SecondPart string
}

type house struct {
	Address     string
	Postcode    postcode
	NumBedrooms uint
	HouseType   string
	Sales       []salePrice
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
	writeJsonFile(res, "house-data.json")
	writeCSVFile(res)
	writeEncodedCSVFile(res)
}

func writeJsonFile(houses any, filename string) {
	bb, err := json.MarshalIndent(houses, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filename, bb, 0644)
	if err != nil {
		panic(err)
	}
}

func writeCSVFile(houses []house) {
	records := [][]string{}

	records = append(records, []string{
		"postcode_first_part",
		"num_bedrooms",
		"house_type",
		"num_sale_records",
		"sale_amount",
		"sale_year",
		"sale_tenure",
	})

	for _, h := range houses {
		for _, s := range h.Sales {
			records = append(records, []string{
				h.Postcode.FirstPart,
				strconv.FormatInt(int64(h.NumBedrooms), 10),
				h.HouseType,
				strconv.FormatInt(int64(len(h.Sales)), 10),
				strconv.FormatInt(int64(s.Amount), 10),
				strconv.FormatInt(int64(s.Date.Year()), 10),
				s.Tenure,
			})
		}
	}

	f, err := os.Create("house-data.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func writeEncodedCSVFile(houses []house) {
	var encodedPostcodeCount int
	encodedPostcode := make(map[string]int, 0)
	decodePostcode := make(map[int]string, 0)
	var encodedHouseTypeCount int
	encodedHouseType := make(map[string]int, 0)
	decodeHouseType := make(map[int]string, 0)
	var encodedTenureCount int
	encodedTenure := make(map[string]int, 0)
	decodeTenure := make(map[int]string, 0)

	records := [][]string{}

	records = append(records, []string{
		"postcode_first_part",
		"num_bedrooms",
		"house_type",
		"num_sale_records",
		"sale_amount",
		"sale_year",
		"sale_tenure",
	})

	for _, h := range houses {
		for _, s := range h.Sales {
			_, ok := encodedPostcode[h.Postcode.FirstPart]
			if !ok {
				encodedPostcode[h.Postcode.FirstPart] = encodedPostcodeCount
				decodePostcode[encodedPostcodeCount] = h.Postcode.FirstPart
				encodedPostcodeCount++
			}

			_, ok = encodedHouseType[h.HouseType]
			if !ok {
				encodedHouseType[h.HouseType] = encodedHouseTypeCount
				decodeHouseType[encodedHouseTypeCount] = h.HouseType
				encodedHouseTypeCount++
			}

			_, ok = encodedTenure[s.Tenure]
			if !ok {
				encodedTenure[s.Tenure] = encodedTenureCount
				decodeTenure[encodedTenureCount] = s.Tenure
				encodedTenureCount++
			}

			records = append(records, []string{
				strconv.FormatInt(int64(encodedPostcode[h.Postcode.FirstPart]), 10),
				strconv.FormatInt(int64(h.NumBedrooms), 10),
				strconv.FormatInt(int64(encodedHouseType[h.HouseType]), 10),
				strconv.FormatInt(int64(len(h.Sales)), 10),
				strconv.FormatInt(int64(s.Amount), 10),
				strconv.FormatInt(int64(s.Date.Year()), 10),
				strconv.FormatInt(int64(encodedTenure[s.Tenure]), 10),
			})
		}
	}

	f, err := os.Create("encoded-house-data.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	writeJsonFile(encodedPostcode, "encodedPostcode.json")
	writeJsonFile(encodedHouseType, "encodedHouseType.json")
	writeJsonFile(encodedTenure, "encodedTenure.json")
	writeJsonFile(decodePostcode, "decodePostcode.json")
	writeJsonFile(decodeHouseType, "decodeHouseType.json")
	writeJsonFile(decodeTenure, "decodeTenure.json")
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
										addr := s.Text()
										pc := postcodeRegex.FindString(addr)
										spc := strings.Split(pc, " ")
										h.Postcode = postcode{
											FirstPart:  strings.TrimSpace(spc[0]),
											SecondPart: strings.TrimSpace(spc[1]),
										}
										h.Address = addr
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

										h.NumBedrooms = uint(nbn)
										if h.Address == "Apartment 14, Forest Court, Union Street, Chester, Cheshire West And Chester CH1 1AB" {
											fmt.Println(s.Text())
										}
										h.HouseType = strings.TrimSpace(nbb[1])
									}

									if a == "subTitle " {
										nb := s.Text()
										h.NumBedrooms = 1
										h.HouseType = strings.TrimSpace(nb)
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
															sp.Amount = uint(pa)
														}
														// date-sold
														if a == "date-sold" {
															ds := s.Text()
															ds = strings.Replace(ds, "(New Build)", "", 1)
															parseTime, err := time.Parse("2 Jan 2006", ds)
															if err != nil {
																panic(err)
															}
															sp.Date = parseTime
														}
														// Freehold or leasehold
														if a == "table-extra tenure" {
															sp.Tenure = s.Text()
														}
													})

												h.Sales = append(h.Sales, sp)
												sp = salePrice{}
											})
									}
								})

							results = append(results, h)
							h = house{}
						})
				})
		})

	return results
}
