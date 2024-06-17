package dataimport

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/wbergg/beerel-roll/db"
)

type Systemet struct {
	All []struct {
		ProductNumber            string  `json:"productNumber"`
		ProductNameBold          string  `json:"productNameBold"`
		RestrictedParcelQuantity int     `json:"restrictedParcelQuantity"`
		AssortmentText           string  `json:"assortmentText"`
		Volume                   float64 `json:"volume"`
	} `json:"all"`
	Store []struct {
		ProductNumber            string  `json:"productNumber"`
		ProductNameBold          string  `json:"productNameBold"`
		RestrictedParcelQuantity int     `json:"restrictedParcelQuantity"`
		AssortmentText           string  `json:"assortmentText"`
		Volume                   float64 `json:"volume"`
	} `json:"store"`
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

	// Insert All slice into DB
	for _, product := range data.All {
		if product.RestrictedParcelQuantity == 0 {
			err := d.Insert(product.ProductNameBold, product.ProductNumber, product.Volume)
			if err != nil {
				fmt.Println(err, product.ProductNameBold, product.ProductNumber, product.Volume)
			}
		}
	}

	// Update inventory with store only items
	for _, product := range data.Store {
		err := d.Insert(product.ProductNameBold, product.ProductNumber, product.Volume)
		if err != nil {
			fmt.Println(err, product.ProductNameBold, product.ProductNumber, product.Volume)
		}
		// Set non orderable products to false
		d.UpdateOrderable(product.ProductNumber)
	}

	return nil
}
