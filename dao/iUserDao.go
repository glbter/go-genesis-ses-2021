package dao

import(
	"github.com/glbter/go-genesis-ses-2021/model"
)

type iUserDao interface {
	create() model.UserLocal
	getById() model.UserLocal
	getByEmail() model.UserLocal
	getAll() []model.UserLocal
}