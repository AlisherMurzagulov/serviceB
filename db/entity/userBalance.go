package entity

type UserBalance struct {
	UserId  string  `json:"userId"gorm:"unique; not null; column:user_id; primaryKey; type:character varying(100)"`
	Balance float64 `json:"balance"`
}

func (UserBalance) TableName() string {
	return "public.QuoteMarkets"
}
