package models

import (
	"time"

	"gorm.io/gorm"
)

// AssetAssignmentAction represents assignment lifecycle actions.
type AssetAssignmentAction string

const (
	AssetAssignmentActionAssign   AssetAssignmentAction = "assign"
	AssetAssignmentActionReassign AssetAssignmentAction = "reassign"
	AssetAssignmentActionReturn   AssetAssignmentAction = "return"
)

// AssetAssignmentHistory stores immutable assignment/reassignment/return audit records.
type AssetAssignmentHistory struct {
	ID             uint                  `json:"id" gorm:"primaryKey"`
	TenantID       *uint                 `json:"tenant_id,omitempty" gorm:"index"`
	AssetID        uint                  `json:"asset_id" gorm:"index;not null"`
	AssetCode      string                `json:"asset_code" gorm:"size:50;index;not null"`
	Action         AssetAssignmentAction `json:"action" gorm:"type:varchar(20);index;not null"`
	FromEmployeeID *string               `json:"from_employee_id,omitempty" gorm:"size:50;index"`
	ToEmployeeID   *string               `json:"to_employee_id,omitempty" gorm:"size:50;index"`
	FromEmployee   *string               `json:"from_employee,omitempty" gorm:"size:255"`
	ToEmployee     *string               `json:"to_employee,omitempty" gorm:"size:255"`
	Department     *string               `json:"department,omitempty" gorm:"size:100"`
	Notes          *string               `json:"notes,omitempty" gorm:"type:text"`
	OccurredAt     time.Time             `json:"occurred_at" gorm:"index;not null"`
	ActorUserID    *uint                 `json:"actor_user_id,omitempty" gorm:"index"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	DeletedAt      gorm.DeletedAt        `json:"-" gorm:"index"`
}

// TableName specifies the table name.
func (AssetAssignmentHistory) TableName() string {
	return "asset_assignment_histories"
}

