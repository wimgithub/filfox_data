package mysql

import (
	"filfox_data/models"
	"fmt"
	"strings"
)

func (s *Store) AddFilData(data []*model.Data) error {
	if len(data) == 0 {
		return nil
	}
	var valueStrings []string
	for _, tx := range data {
		valueString := fmt.Sprintf("(%v, '%v', %v,'%v', '%v', '%v', '%v')",
			tx.Time, tx.FilFrom, tx.Height, tx.Message, tx.FilTo, tx.Type, tx.Value)
		valueStrings = append(valueStrings, valueString)
	}
	sql := fmt.Sprintf("INSERT IGNORE INTO fil_data (time,fil_from,height,message,fil_to,type,value) VALUES %s",
		strings.Join(valueStrings, ","))
	return s.db.Exec(sql).Error
}

func (s *Store) GetFilFoxCount() (count int64, err error) {
	err = s.db.Model(&model.Data{}).Count(&count).Error
	return
}
