package models

type StockDecreaseLog struct {
	Id           int64 `json:"id" gorm:"primarykey;auto_increment"`
	OrderId      int64 `json:"orderid"`
	ProductRefer int64 `json:"productrefer"`
}

type Product struct {
	Id                int64            `json:"id" gorm:"primarykey;auto_increment"`
	Name              string           `json:"name"`
	Stock             int64            `json:"stock"`
	Price             int64            `json:"price"`
	StockDecreaseLogs StockDecreaseLog `gorm:"foreignKey:ProductRefer"`
}
