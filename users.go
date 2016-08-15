package upms

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"git.code4.in/mobilegameserver/logging"
)

type User struct {
	Id        uint64
	Name      string     `gorm:"type:varchar(64);unique_index;not null;default '';"`
	Password  string     `gorm:"type:varchar(128);not null;default '';"`
	Salt      string     `gorm:"type:varchar(128);not null;default '';"`
	OwnerId   uint64     `gorm:"not null;default 0;"`
	Code      string     `gorm:"type:varchar(32);default '';index;"`
	Locked    uint32     `gorm:"not null;default 0;"`
	LockIp    string     `gorm:"type:varchar(64);not null;default '';"`
	Roles     []Role     `gorm:"many2many:user_roles;"`
	Resources []Resource `gorm:"many2many:user_resources;"`
	CreatedAt uint64     `gorm:"not null;default 0;"`
	ChangedAt uint64     `gorm:"not null;default 0;"`
}

func InitUser() {
	db_upms.AutoMigrate(&User{})
}

func (self *User) Save() error {
	return db_upms.Save(self).Error
}

func (self *User) IsLocked() bool {
	return self.Locked == 1
}

func (self *User) UnLock() error {
	self.Locked = 0
	return db_upms.Model(self).Updates(map[string]interface{}{"locked": 1, "changed_at": uint64(time.Now().Unix())}).Error
}

func (self *User) Lock() error {
	self.Locked = 1
	return db_upms.Model(self).Updates(map[string]interface{}{"locked": 1, "changed_at": uint64(time.Now().Unix())}).Error
}

func (self *User) AddRole(rid uint64) error {
	res := LoadRole(rid)
	if res.Id == 0 {
		return errors.New("role id error!")
	}
	return db_upms.Model(self).Association("Roles").Append(res).Error
}

func (self *User) DelRole(rid uint64) error {
	res := LoadRole(rid)
	if res.Id == 0 {
		return errors.New("role id error!")
	}
	return db_upms.Model(self).Association("Roles").Delete(res).Error
}

func (self *User) LoadRole() error {
	return db_upms.Model(self).Association("Roles").Find(&self.Roles).Error
}

func (self *User) UpdateRoles(roleids string) error {
	roleids = strings.TrimSpace(roleids)
	role_list := make([]interface{}, 0)
	roleid_list := strings.Fields(roleids)
	for _, roleid := range roleid_list {
		tmpid, err := strconv.ParseUint(roleid, 10, 64)
		if err != nil {
			continue
		}
		res := LoadRole(tmpid)
		if res.Id == 0 {
			continue
		}
		role_list = append(role_list, res)
	}
	return db_upms.Model(self).Association("Resources").Replace(role_list...).Error
}

func (self *User) AddResource(resid uint64) error {
	res := LoadResource(resid)
	if res.Id == 0 {
		return errors.New("resource id error!")
	}
	return db_upms.Model(self).Association("Resources").Append(res).Error
}

func (self *User) DelResource(resid uint64) error {
	res := LoadResource(resid)
	if res.Id == 0 {
		return errors.New("resource id error!")
	}
	return db_upms.Model(self).Association("Resources").Delete(res).Error
}

func (self *User) LoadResource() error {
	return db_upms.Model(self).Association("Resources").Find(&self.Resources).Error
}

func (self *User) UpdateResource(resids string) error {
	resids = strings.TrimSpace(resids)
	res_list := make([]interface{}, 0)
	resid_list := strings.Fields(resids)
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

func (self *User) UpdateLockIp(ip string) error {
	self.LockIp = ip
	return db_upms.Model(self).Updates(&User{LockIp: ip, ChangedAt: uint64(time.Now().Unix())}).Error
}

func (self *User) UpdatePassword(newpasswd, oldpasswd string) error {
	if self.CheckPassword(oldpasswd) {
		self.Salt = GetRandomSalt()
		self.Password = MD5String(newpasswd + self.Salt)
		return db_upms.Model(self).Updates(map[string]interface{}{"password": self.Password, "salt": self.Salt, "changed_at": time.Now().Unix()}).Error
	}
	return errors.New("password error")
}

func (self *User) CheckPassword(passwd string) bool {
	if MD5String(passwd+self.Salt) == self.Password {
		return true
	}
	return false
}

func (self *User) GetMyRoleIds() map[uint64]int {
	retm := make(map[uint64]int)
	for _, role := range self.Roles {
		retm[role.Id] = 1
	}
	return retm
}

func (self *User) GetMyResourceIds() map[uint64]int {
	retm := make(map[uint64]int)
	for _, res := range self.Resources {
		retm[res.Id] = 1
	}
	return retm
}

func (self *User) CreateUser(name, password, roleids, lockip string, groupid uint64) (*User, error) {
	if self.OwnerId != 0 && groupid == 0 {
		return nil, errors.New("groupid error")
	}

	u, err := NewUser(name, password, "", lockip, groupid, self.Id)
	if err != nil {
		logging.Error("CreateUser error:%s", err.Error())
		return u, err
	}

	roleids = strings.TrimSpace(roleids)
	rolemap := self.GetMyRoleIds()
	update := false
	for _, roleid := range strings.Split(roleids, ",") {
		tmpid, err := strconv.ParseUint(roleid, 10, 64)
		if err != nil {
			continue
		}
		if _, ok := rolemap[tmpid]; ok {
			role := LoadRole(tmpid)
			if role.Id != 0 {
				u.Roles = append(u.Roles, role)
				update = true
			}
		}
	}
	if update {
		u.Code = strconv.FormatUint(u.Id, 10) + "_" + strconv.FormatUint(uint64(time.Now().Unix()), 10) + "_" + GetRandomString(8)
		u.Save()
	}
	if groupid != 0 {
		group := LoadGroupById(groupid)
		if group.Id != 0 {
			group.AddMember(*u)
		}
	}
	return u, nil
}

func (self *User) UpdateResourceForUser(userid uint64, resids string) error {
	if self.OwnerId == 0 {
		err := errors.New("user id or resource id error")
		if userid == 0 || resids == "" {
			return err
		}
		u := FindUserById(userid)
		if u.Id == 0 {
			return err
		}
		return u.UpdateResource(resids)
	}
	return errors.New("no permissions")
}

func (self *User) CreateGroup(name, desc string) *Group {
	if self.OwnerId == 0 {
		return NewGroup(name, desc, self.Id)
	}
	return nil
}

func (self *User) CreateRole(name, desc string, status uint32) *Role {
	if self.OwnerId == 0 {
		return NewRole(name, desc, status, "")
	}
	return nil
}

func (self *User) UpdateResourceForRole(roleid uint64, resids string) error {
	if self.OwnerId == 0 {
		err := errors.New("role id or resource id error")
		if roleid == 0 || resids == "" {
			return err
		}
		role := LoadRole(roleid)
		if role.Id == 0 {
			return err
		}
		return role.UpdateResource(resids)
	}
	return errors.New("no permissions")
}

func (self *User) CreateResource(name, class, attr string, restype, order, status uint32, parentid uint64) *Resource {
	if self.OwnerId == 0 {
		return NewResource(name, class, attr, restype, order, status, parentid)
	}
	return nil
}

func NewUser(name, password, roleids, lockip string, groupid, owner uint64) (*User, error) {
	salt := GetRandomSalt()
	password = MD5String(password + salt)
	u := &User{Name: name, Password: password, Salt: salt, LockIp: lockip, OwnerId: owner, CreatedAt: uint64(time.Now().Unix())}
	err := db_upms.Create(u).Error
	if err != nil {
		return nil, err
	}
	roleids = strings.TrimSpace(roleids)
	if roleids != "" {
		u.UpdateRoles(roleids)
	}
	if groupid != 0 {
		group := LoadGroupById(groupid)
		if group.Id != 0 {
			group.AddMember(*u)
		}
	}
	return u, nil
}

func FindUserByName(name string) *User {
	u := &User{}
	db_upms.Where("name = ?", name).First(u)
	return u
}

func FindUserById(id uint64) *User {
	u := &User{}
	db_upms.Where("id = ?", id).First(u)
	return u
}

func UserLogin(name, password string, ip string) (*User, error) {
	u := FindUserByName(name)
	if u == nil {
		return nil, errors.New("账号或密码错误")
	}
	if !u.CheckPassword(password) {
		return nil, errors.New("账号或密码错误")
	}
	if u.LockIp != "" && u.LockIp != ip {
		return nil, errors.New("登陆IP非法")
	}
	if u.IsLocked() {
		return nil, errors.New("账号被封，请联系管理员")
	}
	group := LoadGroupByMemberId(u.Id)
	if group.Status != 0 {
		return nil, errors.New("用户组禁用")
	}
	logging.Info("user:%s, ip:%s, login success!", name, ip)
	return u, nil
}
