package dto

// ---- Entrada ----

type SaleDetailInput struct {
	Product  uint `json:"product" binding:"required"`
	Quantity uint `json:"quantity" binding:"required"`
}

type SalePaymentInput struct {
	Method string `json:"method" binding:"required"`
	Amount uint   `json:"amount" binding:"required"`
}

type SaleCreateInput struct {
	Details  []SaleDetailInput  `json:"details" binding:"required,dive"`
	Payments []SalePaymentInput `json:"payments" binding:"required,dive"`
}

// ---- Salida ----

type SaleDetailResponse struct {
	ID          uint   `json:"id"`
	Product     uint   `json:"product"`
	ProductName string `json:"product_name"`
	UnitPrice   uint   `json:"unit_price"`
	Quantity    uint   `json:"quantity"`
	Subtotal    uint   `json:"subtotal"`
}

type SalePaymentResponse struct {
	ID     uint   `json:"id"`
	Method string `json:"method"`
	Amount uint   `json:"amount"`
}

type SaleResponse struct {
	ID           uint                  `json:"id"`
	Total        uint                  `json:"total"`
	CreatedAt    string                `json:"created_at"`
	Items        []SaleDetailResponse  `json:"items"`
	PaymentsList []SalePaymentResponse `json:"payments_list"`
}

type SaleVoidInput struct {
	Reason string `json:"reason"` // opcional
}
