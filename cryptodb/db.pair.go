package cryptodb

import (
	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) CreatePair(p *Pair) (err error) {
	result := db.gorm.Create(p)

	return result.Error
}

func (db *Database) GetPairs() (pairs []Pair, err error) {
	result := db.gorm.Order("ID ASC").Find(&pairs)

	return pairs, result.Error
}

func (db *Database) GetPair(id uint) (pair Pair, err error) {
	result := db.gorm.Where("ID = ?", id).First(&pair)

	return pair, result.Error
}

func (db *Database) GetPairByName(s string) (pair Pair, err error) {
	result := db.gorm.Where("Name = ?", s).First(&pair)

	return pair, result.Error
}

// find Pair.Name containing s
func (db *Database) FindPairNames(s string) (pairs []string, err error) {
	result := db.gorm.Model(&Pair{}).Select("name").Where("name LIKE ?", "%"+s+"%").Find(&pairs)

	return pairs, result.Error
}
