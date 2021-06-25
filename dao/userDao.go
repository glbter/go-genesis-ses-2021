package dao

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/glbter/go-genesis-ses-2021/model"
)

type UserDao struct {
	File string
}

func (dao UserDao) Create(user model.UserLocal) {
	csvFile, err := os.OpenFile(dao.File, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatal("unable to open file", err)
	}
	defer csvFile.Close()

	// w := csv.NewWriter(csvFile)
	w := csv.NewWriter(csvFile)
	defer w.Flush()

	var usr = []string{user.Id, user.Name, user.Email, user.Password}
	err = w.Write(usr)

	if err != nil {
		log.Fatal(err)
	}
}

func (dao UserDao) GetById(id string) model.UserLocal {

	csvFile, err := os.Open(dao.File)
	if err != nil {
		log.Fatal("unable to open file")
	}
	defer csvFile.Close()

	var user model.UserLocal
	r := csv.NewReader(csvFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if record[0] == id {
			user = model.UserLocal{record[0], record[1], record[2], record[3]}
			return user
		}
	}
	return user
}

func (dao UserDao) GetByEmail(email string) model.UserLocal {
	csvFile, err := os.Open(dao.File)
	if err != nil {
		log.Fatal("unable to open file")
	}
	defer csvFile.Close()

	var user model.UserLocal
	r := csv.NewReader(csvFile)
	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if record[2] == email {
			user = model.UserLocal{record[0], record[1], record[2], record[3]}
			return user
		}
	}
	return user
}

func (dao UserDao) GetAll() []model.UserLocal {

	csvFile, err := os.Open(dao.File)
	if err != nil {
		log.Fatal("unable to open file")
		return nil
	}
	defer csvFile.Close()

	//r := csv.NewReader(strings.NewReader(csvFile)).ReadAll()
	r := csv.NewReader(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	var users []model.UserLocal

	csvLines, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range csvLines {
		user := model.UserLocal{line[0], line[1], line[2], line[3]}
		users = append(users, user)
	}

	return users
}
