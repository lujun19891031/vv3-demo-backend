package models

import (
	"go_dev/gogin/common"
)

type Produce struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	CreatedTime string  `json:"created_time"`
	Img         string  `json:"img"`
	Desc        string  `json:"desc"`
}

var produces Produce

func Init() {
	err := common.Db.AutoMigrate(&Produce{})
	if err != nil {
		panic(produces.TableName() + err.Error())
	}
}

// 获取数据库名
func (p *Produce) TableName() string {
	return "produce"
}

// 获取所有商品
func (p *Produce) GetAll() ([]Produce, error) {
	var produces []Produce
	err := common.Db.Find(&produces).Error
	return produces, err
}

// 根据id获取商品，并按照时间降序排列
func (p *Produce) GetById(id int) (Produce, error) {
	var produce Produce
	err := common.Db.Where("id = ?", id).First(&produce).Error
	return produce, err
}
