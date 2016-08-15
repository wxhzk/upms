package upms

/**
*用户组，每个用户只能属于一个组，但是一个用户可以属于多个组的管理员
**/

import "time"

type Group struct {
	Id        uint64
	Name      string `gorm:"type:varchar(64);not null;default '';unique"`
	Desc      string `gorm:"type:varchar(128);not null;default'';"`
	Admins    []User `gorm:"many2many:group_admins;"`
	Members   []User `gorm:"many2many:group_members;"`
	OwnerId   uint64 `gorm:"unique;not null;default 0;"`
	Status    uint32 `gorm:"not null;default 0;"`
	CreatedAt uint64 `gorm:"not null;default 0;"`
	UpdatedAt uint64 `gorm:"not null;defaullt 0;"`
}

func InitGroup() {
	db_upms.AutoMigrate(&Group{})
}

func NewGroup(name, desc string, owner uint64) *Group {
	group := &Group{}
	group.Name = name
	group.Desc = desc
	group.OwnerId = owner
	group.CreatedAt = uint64(time.Now().Unix())
	db_upms.Create(group)
	return group
}

func (self *Group) Save() error {
	return db_upms.Save(self).Error
}

func (self *Group) UpdateName(name string) error {
	self.Name = name
	return db_upms.Model(self).Updates(map[string]interface{}{"name": name, "updated_at": time.Now().Unix()}).Error
}

func (self *Group) UpdateDesc(desc string) error {
	self.Desc = desc
	return db_upms.Model(self).Updates(map[string]interface{}{"desc": desc, "updated_at": time.Now().Unix()}).Error
}

func (self *Group) UpdateStatus(status uint32) error {
	self.Status = status
	return db_upms.Model(self).Updates(map[string]interface{}{"status": status, "updated_at": time.Now().Unix()}).Error
}

func (self *Group) AddAdmin(u User) error {
	for _, user := range self.Admins {
		if user.Id == u.Id {
			return nil
		}
	}
	return db_upms.Model(self).Association("Admins").Append(u).Error
}

func (self *Group) DelAdmin(u User) error {
	for _, user := range self.Admins {
		if user.Id == u.Id {
			return db_upms.Model(self).Association("Admins").Delete(u).Error
		}
	}
	return nil
}

func (self *Group) LoadAdmin() error {
	return db_upms.Model(self).Association("Admins").Find(&self.Admins).Error
}

func (self *Group) AddMember(u User) error {
	for _, m := range self.Members {
		if m.Id == u.Id {
			return nil
		}
	}
	return db_upms.Model(self).Association("Members").Append(u).Error
}

func (self *Group) DelMember(u User) error {
	return db_upms.Model(self).Association("Members").Delete(u).Error
}

func (self *Group) LoadMembers() error {
	return db_upms.Model(self).Association("Members").Find(&self.Members).Error
}

func LoadGroupById(groupid uint64) Group {
	var group Group
	db_upms.Where("id = ?", groupid).First(&group)
	return group
}

func LoadGroupByMemberId(memid uint64) Group {
	var group Group
	db_upms.Joins("join group_members on groups.id = group_members.group_id").Where("group_members.user_id = ?", memid).First(&group)
	return group
}

func LoadGroupByAdminId(adminid uint64) []Group {
	var group_list []Group
	db_upms.Joins("Join group_admins on groups.id = group_admins.group_id").Where("group_admins.user_id = ?", adminid).Find(&group_list)
	return group_list
}

func LoadGroupByOwnerId(ownerid uint64) []Group {
	group_list := make([]Group, 0)
	db_upms.Where("owner_id = ? and status = 0", ownerid).Find(&group_list)
	return group_list
}

func LoadGroups() []Group {
	group_list := make([]Group, 0)
	db_upms.Where("status = ?", 0).Find(&group_list)
	return group_list
}
