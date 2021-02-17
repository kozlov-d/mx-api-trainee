package data

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kozlov-d/mx-api-trainee/common"
)

//Data wraps data objects and implements access methods
type Data struct {
	db    *pgxpool.Pool
	Cache *Cache
}

//InitializeDB initializes DB connection pool with provided config
func (d *Data) InitializeDB(dbc common.DBConfig) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbc.Host, dbc.Port, dbc.User, dbc.Password, dbc.Name)
	pool, err := pgxpool.Connect(context.TODO(), connString)
	if err != nil {
		log.Fatal(err)
	}
	d.db = pool
}

func (d *Data) GetMerchants(oID, mID int, sub string) []*common.Merchant {
	sel := `SELECT * FROM offers WHERE 
	($1 = 0 OR offerid=$1) AND 
	($2 = 0 OR merchantid=$2) AND 
	($3 = '' OR offername LIKE ('%' || $3 || '%'))`
	rows, err := d.db.Query(context.TODO(), sel, oID, mID, sub)
	merchants := []*common.Merchant{}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}
	defer rows.Close()
	merch := &common.Merchant{}
	for rows.Next() {
		offer := common.Offer{}
		if err := rows.Scan(&offer.OfferID, &offer.OfferName,
			&offer.Price, &offer.Quantity, &offer.MerchantID); err != nil {
			log.Fatal(err)
		}
		if merch.MerchantID == 0 {
			merch.MerchantID = offer.MerchantID
		}
		if merch.MerchantID == offer.MerchantID {
			merch.Offers = append(merch.Offers, offer)
		}
		if merch.MerchantID != offer.MerchantID {
			merch = &common.Merchant{}
			merch.Offers = append(merch.Offers, offer)
		}
		if !contains(merch.MerchantID, merchants) {
			merchants = append(merchants, merch)
		}
	}
	return merchants
}

func contains(ID int, ms []*common.Merchant) bool {
	if ID == 0 {
		return true
	}
	for _, val := range ms {
		if val.MerchantID == ID {
			return true
		}
	}
	return false
}

func (d *Data) UpsertMerchant(ID int) (*common.Merchant, error) {
	ups := `INSERT INTO merchants(merchantid) VALUES ($1) 
	ON CONFLICT (merchantid) 
	DO UPDATE SET merchantid=EXCLUDED.merchantid 
	RETURNING *`
	m := common.Merchant{}
	err := d.db.QueryRow(context.TODO(), ups, ID).Scan(&m.MerchantID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

//DeleteOffer with given offer and merchant IDs
func (d *Data) DeleteOffer(oID, mID int) (int64, error) {
	del := `DELETE FROM offers WHERE offerid=$1 AND merchantid=$2`
	tag, err := d.db.Exec(context.TODO(), del, oID, mID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (d *Data) SelectOffer(oID, mID int) (*common.Offer, error) {
	sel := `SELECT * FROM offers WHERE offerid=$1 AND merchantid=$2`
	o := &common.Offer{}
	err := d.db.QueryRow(context.TODO(), sel, oID, mID).
		Scan(&o.OfferID, &o.OfferName, &o.Price, &o.Quantity, &o.MerchantID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return o, nil
}

func (d *Data) InsertOffer(o *common.Offer, mID int) error {
	ins := `INSERT INTO 
	offers(offerid, offername, price, quantity, merchantid) 
	VALUES ($1, $2, $3, $4, $5) RETURNING *`
	scan := &common.Offer{} //mb should be returned
	err := d.db.QueryRow(context.TODO(), ins,
		o.OfferID, o.OfferName, o.Price, o.Quantity, mID).
		Scan(&scan.OfferID, &scan.OfferName, &scan.Price, &scan.Quantity, &scan.MerchantID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

func (d *Data) UpdateOffer(o *common.Offer) error {
	upd := `UPDATE offers SET offername=$2, price=$3, quantity=$4 
	WHERE offerid=$1 AND merchantid=$5`
	err := d.db.QueryRow(context.TODO(), upd,
		o.OfferID, o.OfferName, o.Price, o.Quantity, o.MerchantID).Scan()
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}
