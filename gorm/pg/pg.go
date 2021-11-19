package pg

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Postgres struct {
	Name   string
	Schema string
	Args   string
	Model  interface{}
	Db     *gorm.DB
}

func (p *Postgres) GetDb() error {
	var err error
	if p.Db == nil {
		p.Db, err = p.openDb()
		return err
	}
	if err := p.Db.DB().Ping(); err != nil {
		log.Println("ping error")
		p.Db.Close()
		p.Db, err = p.openDb()
		return err
	}
	return nil
}

func (p *Postgres) openDb() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", p.Args+"&search_path="+p.Schema)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(25)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)
	if p.Model != nil {
		if p.Schema != "" && p.Schema != "public" {
			if err := db.Exec("create schema if not exists " + p.Schema).Error; err != nil {
				return nil, err
			}
		}
		db = db.Table(p.Name)
		if db.HasTable(p.Name) {
			if err = db.AutoMigrate(p.Model).Error; err != nil {
				return nil, err
			}
		} else {
			if err = db.CreateTable(p.Model).Error; err != nil {
				return nil, err
			}
		}
	}
	return db, nil
}
