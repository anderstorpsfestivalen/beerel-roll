package dataimport

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/wbergg/beerel-roll/db"
)

type Systemet struct {
	Ordered []Item `json:"ordered"`
	Store   []Item `json:"store"`
}

type Item struct {
	ProductNumber            string  `json:"productNumber"`
	ProductId                string  `json:"productId"`
	ProductNameBold          string  `json:"productNameBold"`
	ProductNameThin          string  `json:"productNameThin"`
	RestrictedParcelQuantity int     `json:"restrictedParcelQuantity"`
	AssortmentText           string  `json:"assortmentText"`
	Volume                   float64 `json:"volume"`
	Images                   []Image `json:"images"`
}

type Image struct {
	ImageURL string      `json:"imageUrl"`
	FileType interface{} `json:"fileType"`
	Size     interface{} `json:"size"`
}

func DbSetup(d *db.DBobject) error {
	// Open the JSON file
	jsonFile, err := os.Open("dataimport/all.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// Read the file's content
	byteValue, _ := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return err
	}

	// Unmarshal the byteArray into the Systemet struct
	var data Systemet

	if err := json.Unmarshal(byteValue, &data); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	// this is a fuuuuuuugly hack but whatever
	baseImageURL := "https://product-cdn.systembolaget.se/productimages/"

	// Insert All slice into DB
	for _, product := range data.Ordered {
		if product.RestrictedParcelQuantity == 0 {
			pid, _ := strconv.ParseInt(product.ProductNumber, 10, 64)
			err := d.Insert(formatProductName(product.ProductNameBold, product.ProductNameThin), pid, product.Volume, baseImageURL+product.ProductId+"/"+product.ProductId+"_300.png")
			if err != nil {
				fmt.Println(err, formatProductName(product.ProductNameBold, product.ProductNameThin), product.ProductNumber, product.Volume)
			}
		}
	}

	// Update inventory with store only items
	for _, product := range data.Store {

		pid, _ := strconv.ParseInt(product.ProductNumber, 10, 64)
		err := d.Insert(formatProductName(product.ProductNameBold, product.ProductNameThin), pid, product.Volume, product.Images[0].ImageURL+"_300.png")
		if err != nil {
			fmt.Println(err, formatProductName(product.ProductNameBold, product.ProductNameThin), product.ProductNumber, product.Volume)
		}
		// Set non orderable products to false
		d.UpdateOrderable(pid)
	}

	return nil
}

func formatProductName(pNameBold string, pNameThin string) string {
	if pNameThin != "" {
		return pNameBold + " " + pNameThin
	}
	return pNameBold
}
