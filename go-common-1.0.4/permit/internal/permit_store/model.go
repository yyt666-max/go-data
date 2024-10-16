package permit_store

import "time"

type Permit struct {
	Id         int64     `gorm:"column:id;type:bigint(20);not null;primary_key;auto_increment;comment:'主键'" `
	Domain     string    `gorm:"column:domain;type:varchar(255);not null;comment:'权限域';uniqueIndex:domain_access_target;"`
	Access     string    `gorm:"column:access;type:varchar(100);not null;comment:'权限标识';uniqueIndex:domain_access_target;" `
	Target     string    `gorm:"column:target;type:varchar(64);not null;comment:'权限目标';uniqueIndex:domain_access_target;" `
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null;comment:'创建时间'" `
}

func (p *Permit) IdValue() int64 {
	return p.Id
}
func (p *Permit) TableName() string {
	return "permit"
}
