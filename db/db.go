package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DBobject struct {
	db    *sqlx.DB
	Setup int
}

type Inventory struct {
	ID            int             `db:"id"`
	Name          sql.NullString  `db:"name"`
	ProductNumber sql.NullInt64   `db:"product_number"`
	Volume        sql.NullFloat64 `db:"volume"`
	Orderable     sql.NullBool    `db:"orderable"`
	Consumed      sql.NullBool    `db:"consumed"`
	ConsumedBy    sql.NullString  `db:"consumed_by"`
	ConsumedTime  sql.NullTime    `db:"consumed_time"`
	ImageURL      sql.NullString  `db:"image_url"`
	Rejected      sql.NullString  `db:"rejected"`
}

type Beer struct {
	Name          string         `json:"name"`
	ProductNumber int64          `json:"product_number"`
	ConsumedTime  sql.NullTime   `json:"consumed_time"`
	ConsumedBy    sql.NullString `json:"consumed_by"`
	Volume        float64        `json:"volume"`
	ImageURL      string         `json:"image_url"`
	Rejected      string         `json:"rejected"`
}

func Open(dbPath string) DBobject {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	setup := initDB(db)

	return DBobject{
		db:    db,
		Setup: setup,
	}
}

func initDB(db *sqlx.DB) int {

	// Check if table exists
	dbQuery1 := `SELECT COUNT(name) as P FROM sqlite_master WHERE type='table' AND name='inventory';`

	var s1 int
	err := db.Get(&s1, dbQuery1)
	if err != nil {
		panic(err)
	}

	// If tables doesn't exist, create it from the schema
	if s1 == 0 {
		dat, err := os.ReadFile("db/db.schema")
		if err != nil {
			panic(err)
		}
		db.MustExec(string(dat))
	}
	return s1
}

func (d *DBobject) GetAllItems() (Inventory, error) {
	databaseResp := Inventory{}

	dbQuery := "SELECT * FROM inventory"

	err := d.db.Select(&databaseResp, dbQuery)
	if err != nil {
		panic(err)
	}

	return databaseResp, nil
}

func (d *DBobject) GetRandBeer() (Beer, error) {
	databaseResp := Inventory{}

	dbQuery := "SELECT * FROM inventory WHERE consumed = 'false' ORDER BY RANDOM() LIMIT 1"

	err := d.db.Get(&databaseResp, dbQuery)
	if err != nil {
		return Beer{}, err
	}

	if !databaseResp.Name.Valid {
		return Beer{}, fmt.Errorf("beer name wasn't valid?")
	}

	return Beer{
		Name:          databaseResp.Name.String,
		ProductNumber: databaseResp.ProductNumber.Int64,
		ConsumedTime:  databaseResp.ConsumedTime,
		Volume:        databaseResp.Volume.Float64,
		ImageURL:      databaseResp.ImageURL.String,
		Rejected:      databaseResp.Rejected.String,
	}, nil
}

func (d *DBobject) GetNLastConsumed(n int64) ([]Beer, error) {
	databaseResp := []Inventory{}

	dbQuery := "SELECT * FROM inventory WHERE consumed = 'true' ORDER BY consumed_time DESC LIMIT $1"

	err := d.db.Select(&databaseResp, dbQuery, n)
	if err != nil {
		return []Beer{}, err
	}

	beerResp := []Beer{}

	for _, beer := range databaseResp {
		beerResp = append(beerResp, Beer{
			Name:          beer.Name.String,
			ProductNumber: beer.ProductNumber.Int64,
			ConsumedTime:  beer.ConsumedTime,
			ConsumedBy:    beer.ConsumedBy,
			Volume:        beer.Volume.Float64,
			ImageURL:      beer.ImageURL.String,
		})
	}

	return beerResp, nil
}

func (d *DBobject) ConsumeBeer(product_number int64, consumer string) error {

	dbQuery := `UPDATE inventory SET consumed = 'true', consumed_by = $1, consumed_time = CURRENT_TIMESTAMP WHERE product_number = $2`

	_, err := d.db.Exec(dbQuery, consumer, product_number)

	return err

}

func (d *DBobject) RejectBeer(product_number int64, consumer string) error {

	dbQuery := `UPDATE inventory SET rejected = $1 WHERE product_number = $2`

	_, err := d.db.Exec(dbQuery, consumer, product_number)

	return err

}

func (d *DBobject) GetRowItemByPid(pid string) error {

	var result string

	dbQuery := "SELECT * FROM inventory WHERE id=?"

	row := d.db.QueryRow(dbQuery, pid)
	err := row.Scan(&result)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return err
}

func (d *DBobject) Insert(name string, pid int64, volume float64, imageUrl string) error {
	dbQuery := `
		INSERT INTO inventory (name, product_number, volume, consumed, orderable, image_url) VALUES (
			$1, $2, $3, $4, $5, $6
		)`

	_, err := d.db.Exec(dbQuery,
		name,
		pid,
		volume,
		"false",
		"true",
		imageUrl,
	)

	return err
}

func (d *DBobject) UpdateOrderable(pid int64) error {

	dbQuery := `UPDATE inventory SET orderable = 'false' WHERE product_number = $1`

	_, err := d.db.Exec(dbQuery, pid)

	return err

}
