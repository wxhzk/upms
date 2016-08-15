package upms

import (
	"errors"
	"strconv"
	"strings"
)

type Role struct {
	Id        uint64
	Name      string     `gorm:"type:varchar(64);not null;default '';index"`
	Desc      string     `gorm:"type:varchar(128);not null;default '';"`
	Status    uint32     `gorm:"not null;default 0"`
	Resources []Resource `gorm:"many2many:role_resources"`
}

func InitRole() {
	db_upms.AutoMigrate(&Role{})
}

func (self *Role) Save() error {
	return db_upms.Save(self).Error
}

func (self *Role) UpdateName(name string) error {
	self.Name = name
	return db_upms.Model(self).Update("name", name).Error
}

func (self *Role) UpdateDesc(desc string) error {
	self.Desc = desc
	return db_upms.Model(self).Update("desc", desc).Error
}

func (self *Role) UpdateStatus(status uint32) error {
	self.Status = status
	return db_upms.Model(self).Update("status", status).Error
}

func (self *Role) AddResource(resid uint64) error {
	res := LoadResource(resid)
	if res.Id == 0 {
		return errors.New("resource id error!")
	}
	return db_upms.Model(self).Association("Resources").Append(res).Error
}

func (self *Role) DelResource(resid uint64) error {
	res := LoadResource(resid)
	if res.Id == 0 {
		return errors.New("resource id error!")
	}
	return db_upms.Model(self).Association("Resources").Delete(res).Error
}

func (self *Role) ClearResource() error {
	return db_upms.Model(self).Association("Resources").Clear().Error
}

func (self *Role) UpdateResource(resids string) error {
	resids = strings.TrimSpace(resids)
	res_list := make([]interface{}, 0)
	resid_list := strings.Split(resids, ",")
	for _, resid := range resid_list {
		tmpid, err := strconv.ParseUint(resid, 10, 64)
		if err != nil {
			continue
		}
		res := LoadResource(tmpid)
		if res.Id == 0 {
			continue
		}
		res_list = append(res_list, res)
	}
	return db_upms.Model(self).Association("Resources").Replace(res_list...).Error
}

func (self *Role) LoadResource() error {
	return db_upms.Model(self).Association("Resources").Find(&self.Resources).Error
}

func LoadRole(roleid uint64) Role {
	var role Role
	db_upms.Where("id = ?", roleid).First(&role)
	return role
}

func LoadRoles() []Role {
	role_list := make([]Role, 0)
	db_upms.Where("status = ?", 0).Find(&role_list)
	return role_list
}

func NewRole(name, desc string, status uint32, resids string) *Role {
	role := &Role{}
	role.Name = name
	role.Desc = desc
	role.Status = status
	db_upms.Create(role)
	if resids != "" {
		role.UpdateResource(resids)
	}
	return role
}
