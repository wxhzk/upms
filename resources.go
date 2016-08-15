package upms

type Resource struct {
	Id           uint64
	Name         string     `gorm:"type:varchar(64);not null;default '';unique_index;"`
	Class        string     `gorm:"type:varchar(64);not null;default '';"`
	Attr         string     `gorm:"type:varchar(64);not null;default '';"`
	Type         uint32     `gorm:"index;not null;default 0"`
	ParentId     uint64     `gorm:"index;not null;default 0"`
	Order        uint32     `gorm:"not null;default 0"`
	Status       uint32     `gorm:"not null;default 0"`
	SubResources []Resource `gorm:foreign_key:ParentId`
}

func InitResource() {
	db_upms.AutoMigrate(&Resource{})
}

func NewResource(name, class, attr string, restype, order, status uint32, parentid uint64) *Resource {
	res := &Resource{}
	res.Name = name
	res.Class = class
	res.Attr = attr
	res.Type = restype
	res.Order = order
	res.Status = status
	res.ParentId = parentid
	db_upms.Create(res)
	return res
}

func (self *Resource) UpdateName(name string) error {
	self.Name = name
	return db_upms.Model(self).Update("name", name).Error
}

func (self *Resource) UpdateClass(resclass string) error {
	self.Class = resclass
	return db_upms.Model(self).Update("class", resclass).Error
}

func (self *Resource) UpdateAttr(attr string) error {
	self.Attr = attr
	return db_upms.Model(self).Update("attr", attr).Error
}

func (self *Resource) UpdateType(restype uint32) error {
	self.Type = restype
	return db_upms.Model(self).Update("type", restype).Error
}

func (self *Resource) UpdateOrder(order uint32) error {
	self.Order = order
	return db_upms.Model(self).Update("order", order).Error
}

func (self *Resource) UpdateStatus(status uint32) error {
	self.Status = status
	return db_upms.Model(self).Update("status", status).Error
}

func (self *Resource) Save() error {
	return db_upms.Save(self).Error
}

func (self *Resource) Update(updict map[string]interface{}) error {
	return db_upms.Model(self).Updates(updict).Error
}

func (self *Resource) AddSubresource(res *Resource) error {
	return db_upms.Model(self).Association("SubResources").Append(*res).Error
}

func (self *Resource) DelSubresource(res *Resource) error {
	return db_upms.Model(self).Association("SubResources").Delete(*res).Error
}

func (self *Resource) ClearSubresource() error {
	return db_upms.Model(self).Association("SubResources").Clear().Error
}

func (self *Resource) LoadSubresource() error {
	if self.SubResources == nil {
		self.SubResources = make([]Resource, 0)
	}
	return db_upms.Model(self).Association("SubResources").Find(&self.SubResources).Error
}

func LoadResource(resid uint64) Resource {
	var res Resource
	db_upms.Where("id = ?", resid).First(&res)
	return res
}

func LoadResources() []Resource {
	res_list := make([]Resource, 0)
	db_upms.Where("status = ?", 0).Find(&res_list)
	return res_list
}
