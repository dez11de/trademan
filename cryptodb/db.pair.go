package cryptodb

import (
	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) GetPairByName(s string) (pair Pair, err error) {
	result := db.Where("name = ?", s).Find(&pair)

	return pair, result.Error
}

func (db *Database) CrupdatePair(p *Pair) (err error) {
	pair, err := db.GetPairByName(p.Name)
	if err != nil {
		return err
	}

	if pair.ID == 0 {
		result := db.Create(&p)
		return result.Error
	} else {
		result := db.Save(&p)
		return result.Error
	}
}

// TODO: should only return Active pairs. See GORM api documentation.
func (db *Database) GetPairs() (pairs []Pair, err error) {
	result := db.Where("status = ?", "Trading").Order("ID ASC").Find(&pairs)

	return pairs, result.Error
}

func (db *Database) GetPair(id uint64) (pair Pair, err error) {
	result := db.Where("ID = ?", id).First(&pair)

	return pair, result.Error
}
