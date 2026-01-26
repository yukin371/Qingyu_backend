package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Device 设备模型
type Device struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Name      string             `bson:"name"`       // 设备名称（如 "iPhone 13", "Chrome on Windows"）
	Type      string             `bson:"type"`       // 设备类型（mobile, desktop, tablet）
	UserAgent string             `bson:"user_agent"` // 完整User-Agent
	IP        string             `bson:"ip"`         // 登录IP
	LastSeen  time.Time          `bson:"last_seen"`  // 最后活跃时间
	CreatedAt time.Time          `bson:"created_at"`
}

// 设备类型常量
const (
	DeviceTypeMobile  = "mobile"
	DeviceTypeDesktop = "desktop"
	DeviceTypeTablet  = "tablet"
)

// NewDevice 创建新设备
func NewDevice(userID primitive.ObjectID, name, deviceType, userAgent, ip string) *Device {
	now := time.Now()
	return &Device{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Name:      name,
		Type:      deviceType,
		UserAgent: userAgent,
		IP:        ip,
		LastSeen:  now,
		CreatedAt: now,
	}
}

// UpdateLastSeen 更新最后活跃时间
func (d *Device) UpdateLastSeen() {
	d.LastSeen = time.Now()
}
