package models

const ConfigCollection = "configs"
const SettingCollection = "settings"

const (
	GeneralTypeSetting = "general"
	PaymentTypeSetting = "payment"
	StorageTypeSetting = "storage"
	MailTypeSetting    = "mail"
)

type Setting struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Code        string                 `json:"code"`
	Description string                 `json:"description"`
	Setting     map[string]interface{} `json:"setting"`
	Type        string                 `json:"type"` // general, payment, etc
	Status      bool                   `json:"status"`
	IsDefault   bool                   `json:"isDefault" bson:"isDefault"`
}
