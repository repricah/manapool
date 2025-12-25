package manapool

import "net/url"

// PricesMeta describes price export metadata.
type PricesMeta struct {
	AsOf Timestamp `json:"as_of"`
}

// SinglesPricesList represents the singles prices export.
type SinglesPricesList struct {
	Meta PricesMeta           `json:"meta"`
	Data []SinglePriceListing `json:"data"`
}

// SinglePriceListing represents a single price listing.
type SinglePriceListing struct {
	URL                    string  `json:"url"`
	Name                   string  `json:"name"`
	SetCode                string  `json:"set_code"`
	Number                 string  `json:"number"`
	MultiverseID           *string `json:"multiverse_id"`
	ScryfallID             string  `json:"scryfall_id"`
	AvailableQuantity      int     `json:"available_quantity"`
	PriceCents             *int    `json:"price_cents"`
	PriceCentsLPPlus       *int    `json:"price_cents_lp_plus"`
	PriceCentsNM           *int    `json:"price_cents_nm"`
	PriceCentsFoil         *int    `json:"price_cents_foil"`
	PriceCentsLPPlusFoil   *int    `json:"price_cents_lp_plus_foil"`
	PriceCentsNMFoil       *int    `json:"price_cents_nm_foil"`
	PriceCentsEtched       *int    `json:"price_cents_etched"`
	PriceCentsLPPlusEtched *int    `json:"price_cents_lp_plus_etched"`
	PriceCentsNMEtched     *int    `json:"price_cents_nm_etched"`
}

// VariantPricesList represents the variant prices export.
type VariantPricesList struct {
	Meta PricesMeta            `json:"meta"`
	Data []VariantPriceListing `json:"data"`
}

// VariantPriceListing represents a variant price listing.
type VariantPriceListing struct {
	URL                string  `json:"url"`
	ProductType        string  `json:"product_type"`
	ProductID          string  `json:"product_id"`
	SetCode            string  `json:"set_code"`
	Number             string  `json:"number"`
	Name               string  `json:"name"`
	ScryfallID         string  `json:"scryfall_id"`
	TCGPlayerProductID *int    `json:"tcgplayer_product_id"`
	LanguageID         string  `json:"language_id"`
	ConditionID        *string `json:"condition_id"`
	FinishID           *string `json:"finish_id"`
	LowPrice           int     `json:"low_price"`
	AvailableQuantity  int     `json:"available_quantity"`
}

// SealedPricesList represents the sealed prices export.
type SealedPricesList struct {
	Meta PricesMeta           `json:"meta"`
	Data []SealedPriceListing `json:"data"`
}

// SealedPriceListing represents a sealed price listing.
type SealedPriceListing struct {
	URL                string `json:"url"`
	ProductType        string `json:"product_type"`
	ProductID          string `json:"product_id"`
	SetCode            string `json:"set_code"`
	Name               string `json:"name"`
	TCGPlayerProductID *int   `json:"tcgplayer_product_id"`
	LanguageID         string `json:"language_id"`
	LowPrice           int    `json:"low_price"`
	AvailableQuantity  int    `json:"available_quantity"`
}

// OptimizerRequest represents a cart optimization request.
type OptimizerRequest struct {
	Cart               []OptimizerCartItem `json:"cart"`
	Model              string              `json:"model,omitempty"`
	DestinationCountry string              `json:"destination_country,omitempty"`
	ExcludeSellerIDs   []string            `json:"exclude_seller_ids,omitempty"`
	AllowSellerIDs     []string            `json:"allow_seller_ids,omitempty"`
	ShipFromCountries  []string            `json:"ship_from_countries,omitempty"`
}

// OptimizerCartItem represents an item requested by the optimizer.
type OptimizerCartItem struct {
	Type                      string   `json:"type"`
	Name                      string   `json:"name,omitempty"`
	SetCode                   string   `json:"set_code,omitempty"`
	CollectorNumber           string   `json:"collector_number,omitempty"`
	IsToken                   *bool    `json:"is_token,omitempty"`
	IncludeNonSanctionedLegal *bool    `json:"include_non_sanctioned_legal,omitempty"`
	MTGJsonID                 *string  `json:"mtgjson_id,omitempty"`
	LanguageIDs               []string `json:"language_ids,omitempty"`
	FinishIDs                 []string `json:"finish_ids,omitempty"`
	ConditionIDs              []string `json:"condition_ids,omitempty"`
	URI                       string   `json:"uri,omitempty"`
	TCGPlayerSKUIds           []int    `json:"tcgplayer_sku_ids,omitempty"`
	ProductType               string   `json:"product_type,omitempty"`
	ProductIDs                []string `json:"product_ids,omitempty"`
	QuantityRequested         int      `json:"quantity_requested"`
	Index                     *int     `json:"index,omitempty"`
}

// OptimizedCart represents an optimized cart response.
type OptimizedCart struct {
	Cart   []OptimizedCartItem `json:"cart"`
	Totals OptimizedCartTotals `json:"totals"`
}

// OptimizedCartItem represents a selected inventory item.
type OptimizedCartItem struct {
	InventoryID      string `json:"inventory_id"`
	QuantitySelected int    `json:"quantity_selected"`
}

// OptimizedCartTotals represents cart totals.
type OptimizedCartTotals struct {
	SubtotalCents int `json:"subtotal_cents"`
	ShippingCents int `json:"shipping_cents"`
	TotalCents    int `json:"total_cents"`
	SellerCount   int `json:"seller_count"`
}

// BuyerOrdersOptions defines filters for buyer orders.
type BuyerOrdersOptions struct {
	Since  *Timestamp
	Limit  int
	Offset int
}

// BuyerOrdersResponse represents a list of buyer orders.
type BuyerOrdersResponse struct {
	Orders []BuyerOrderSummary `json:"orders"`
}

// BuyerOrderSummary represents a summary of a buyer order.
type BuyerOrderSummary struct {
	ID                string                   `json:"id"`
	CreatedAt         Timestamp                `json:"created_at"`
	SubtotalCents     int                      `json:"subtotal_cents"`
	TaxCents          int                      `json:"tax_cents"`
	ShippingCents     int                      `json:"shipping_cents"`
	TotalCents        int                      `json:"total_cents"`
	OrderNumber       string                   `json:"order_number"`
	OrderSellerDetail []BuyerOrderSellerDetail `json:"order_seller_details"`
}

// BuyerOrderSellerDetail represents seller details in a buyer order.
type BuyerOrderSellerDetail struct {
	OrderNumber    string                  `json:"order_number"`
	SellerID       string                  `json:"seller_id"`
	SellerUsername string                  `json:"seller_username"`
	ItemCount      int                     `json:"item_count"`
	Fulfillments   []BuyerOrderFulfillment `json:"fulfillments,omitempty"`
	Items          []BuyerOrderItem        `json:"items,omitempty"`
}

// BuyerOrderFulfillment represents fulfillment details for a buyer order.
type BuyerOrderFulfillment struct {
	Status          *string    `json:"status"`
	TrackingCompany *string    `json:"tracking_company"`
	TrackingNumber  *string    `json:"tracking_number"`
	TrackingURL     *string    `json:"tracking_url"`
	InTransitAt     *Timestamp `json:"in_transit_at"`
}

// BuyerOrderItem represents an item in a buyer order.
type BuyerOrderItem struct {
	PriceCents int               `json:"price_cents"`
	Quantity   int               `json:"quantity"`
	Product    BuyerOrderProduct `json:"product"`
}

// BuyerOrderProduct represents a product in a buyer order item.
type BuyerOrderProduct struct {
	ProductType string            `json:"product_type"`
	ProductID   string            `json:"product_id"`
	Single      *BuyerOrderSingle `json:"single"`
	Sealed      *BuyerOrderSealed `json:"sealed"`
}

// BuyerOrderSingle represents a single product in a buyer order.
type BuyerOrderSingle struct {
	ScryfallID  string `json:"scryfall_id"`
	MTGJsonID   string `json:"mtgjson_id"`
	Name        string `json:"name"`
	Set         string `json:"set"`
	Number      string `json:"number"`
	LanguageID  string `json:"language_id"`
	ConditionID string `json:"condition_id"`
	FinishID    string `json:"finish_id"`
}

// BuyerOrderSealed represents a sealed product in a buyer order.
type BuyerOrderSealed struct {
	MTGJsonID  string `json:"mtgjson_id"`
	Name       string `json:"name"`
	Set        string `json:"set"`
	LanguageID string `json:"language_id"`
}

// BuyerOrderResponse represents a buyer order response.
type BuyerOrderResponse struct {
	Order BuyerOrderDetails `json:"order"`
}

// BuyerOrderDetails represents detailed buyer order data.
type BuyerOrderDetails struct {
	ID                string                   `json:"id"`
	CreatedAt         Timestamp                `json:"created_at"`
	SubtotalCents     int                      `json:"subtotal_cents"`
	TaxCents          int                      `json:"tax_cents"`
	ShippingCents     int                      `json:"shipping_cents"`
	TotalCents        int                      `json:"total_cents"`
	OrderNumber       string                   `json:"order_number"`
	OrderSellerDetail []BuyerOrderSellerDetail `json:"order_seller_details"`
}

// BuyerCredit represents buyer credit.
type BuyerCredit struct {
	UserCreditCents int `json:"user_credit_cents"`
}

// PendingOrderRequest represents a pending order create/update request.
type PendingOrderRequest struct {
	ShippingOverrides map[string]string      `json:"shipping_overrides,omitempty"`
	LineItems         []PendingOrderLineItem `json:"line_items"`
	TaxAddress        *Address               `json:"tax_address,omitempty"`
	ShippingAddress   *Address               `json:"shipping_address,omitempty"`
}

// PendingOrderLineItem represents a pending order line item.
type PendingOrderLineItem struct {
	InventoryID      string `json:"inventory_id"`
	QuantitySelected int    `json:"quantity_selected"`
}

// PendingOrder represents a pending order response.
type PendingOrder struct {
	ID                string                 `json:"id"`
	ShippingOverrides map[string]string      `json:"shipping_overrides,omitempty"`
	LineItems         []PendingOrderLineItem `json:"line_items"`
	Status            string                 `json:"status"`
	Totals            PendingOrderTotals     `json:"totals"`
	Order             *PendingOrderCompleted `json:"order"`
}

// PendingOrderTotals represents totals for a pending order.
type PendingOrderTotals struct {
	SubtotalCents int `json:"subtotal_cents"`
	ShippingCents int `json:"shipping_cents"`
	TaxCents      int `json:"tax_cents"`
	TotalCents    int `json:"total_cents"`
}

// PendingOrderCompleted represents a completed order reference.
type PendingOrderCompleted struct {
	ID string `json:"id"`
}

// PurchasePendingOrderRequest represents a purchase request.
type PurchasePendingOrderRequest struct {
	PaymentMethod   string  `json:"payment_method"`
	BillingAddress  Address `json:"billing_address"`
	ShippingAddress Address `json:"shipping_address"`
	ApplyBuyerFee   *bool   `json:"apply_buyer_fee,omitempty"`
}

// Address represents a shipping/billing address.
type Address struct {
	Name       string  `json:"name,omitempty"`
	Line1      string  `json:"line1"`
	Line2      *string `json:"line2,omitempty"`
	Line3      *string `json:"line3,omitempty"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	PostalCode string  `json:"postal_code"`
	Country    string  `json:"country"`
}

// SellerAccountUpdate represents a seller account update request.
type SellerAccountUpdate struct {
	SinglesLive *bool `json:"singles_live"`
	SealedLive  *bool `json:"sealed_live"`
}

// DeckCreateRequest represents a deck create request.
type DeckCreateRequest struct {
	CommanderNames []string    `json:"commander_names"`
	OtherCards     []OtherCard `json:"other_cards"`
}

// OtherCard represents a non-commander card in a deck.
type OtherCard struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// DeckCreateResponse represents the deck validation response.
type DeckCreateResponse struct {
	Valid   bool                  `json:"valid"`
	BuyURL  string                `json:"buy_url,omitempty"`
	Details DeckValidationDetails `json:"details"`
}

// DeckValidationDetails represents deck validation details.
type DeckValidationDetails struct {
	CommanderCount          int                          `json:"commander_count"`
	TotalCardCount          int                          `json:"total_card_count"`
	AllCardsLegal           bool                         `json:"all_cards_legal"`
	IllegalCards            []string                     `json:"illegal_cards"`
	CardsNotFound           []string                     `json:"cards_not_found"`
	ValidQuantities         bool                         `json:"valid_quantities"`
	QuantityViolations      []DeckQuantityViolation      `json:"quantity_violations"`
	ValidColorIdentity      bool                         `json:"valid_color_identity"`
	ColorIdentityViolations []DeckColorIdentityViolation `json:"color_identity_violations"`
	ValidPartnership        bool                         `json:"valid_partnership"`
	PartnerViolations       []string                     `json:"partner_violations"`
}

// DeckQuantityViolation represents a quantity violation.
type DeckQuantityViolation struct {
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	MaxAllowed int    `json:"max_allowed"`
}

// DeckColorIdentityViolation represents a color identity violation.
type DeckColorIdentityViolation struct {
	Name            string   `json:"name"`
	CardColors      []string `json:"card_colors"`
	CommanderColors []string `json:"commander_colors"`
}

// InventoryItemsResponse represents a response with inventory items.
type InventoryItemsResponse struct {
	Inventory []InventoryItem `json:"inventory"`
}

// InventoryListingResponse represents a response with a single inventory item.
type InventoryListingResponse struct {
	Inventory InventoryItem `json:"inventory"`
}

// InventoryListingsResponse represents inventory listing items.
type InventoryListingsResponse struct {
	InventoryItems []InventoryItem `json:"inventory_items"`
}

// InventoryItemResponse represents inventory item response for /inventory/listings/{id}.
type InventoryItemResponse struct {
	InventoryItem InventoryItem `json:"inventory_item"`
}

// InventoryUpdateRequest represents a request to update inventory.
type InventoryUpdateRequest struct {
	PriceCents int `json:"price_cents"`
	Quantity   int `json:"quantity"`
}

// InventoryBulkItemBySKU represents bulk inventory items by SKU.
type InventoryBulkItemBySKU struct {
	TCGPlayerSKU int `json:"tcgplayer_sku"`
	PriceCents   int `json:"price_cents"`
	Quantity     int `json:"quantity"`
}

// InventoryBulkItemByProduct represents bulk inventory items by product.
type InventoryBulkItemByProduct struct {
	ProductType string `json:"product_type"`
	ProductID   string `json:"product_id"`
	PriceCents  int    `json:"price_cents"`
	Quantity    int    `json:"quantity"`
}

// InventoryBulkItemByScryfall represents bulk inventory items by Scryfall ID.
type InventoryBulkItemByScryfall struct {
	ScryfallID  string `json:"scryfall_id"`
	LanguageID  string `json:"language_id"`
	FinishID    string `json:"finish_id"`
	ConditionID string `json:"condition_id"`
	PriceCents  int    `json:"price_cents"`
	Quantity    int    `json:"quantity"`
}

// InventoryBulkItemByTCGPlayerID represents bulk inventory items by TCGPlayer ID.
type InventoryBulkItemByTCGPlayerID struct {
	TCGPlayerID int     `json:"tcgplayer_id"`
	LanguageID  string  `json:"language_id"`
	FinishID    *string `json:"finish_id"`
	ConditionID *string `json:"condition_id"`
	PriceCents  int     `json:"price_cents"`
	Quantity    int     `json:"quantity"`
}

// InventoryByScryfallOptions defines lookup options by Scryfall ID.
type InventoryByScryfallOptions struct {
	LanguageID  string
	FinishID    string
	ConditionID string
}

func (opts InventoryByScryfallOptions) toParams() url.Values {
	params := url.Values{}
	if opts.LanguageID != "" {
		params.Add("language_id", opts.LanguageID)
	}
	if opts.FinishID != "" {
		params.Add("finish_id", opts.FinishID)
	}
	if opts.ConditionID != "" {
		params.Add("condition_id", opts.ConditionID)
	}
	return params
}

// InventoryByTCGPlayerOptions defines lookup options by TCGPlayer ID.
type InventoryByTCGPlayerOptions struct {
	LanguageID  string
	FinishID    string
	ConditionID string
}

func (opts InventoryByTCGPlayerOptions) toParams() url.Values {
	params := url.Values{}
	if opts.LanguageID != "" {
		params.Add("language_id", opts.LanguageID)
	}
	if opts.FinishID != "" {
		params.Add("finish_id", opts.FinishID)
	}
	if opts.ConditionID != "" {
		params.Add("condition_id", opts.ConditionID)
	}
	return params
}

// OrdersOptions defines filters for order listing endpoints.
type OrdersOptions struct {
	Since           *Timestamp
	IsUnfulfilled   *bool
	IsFulfilled     *bool
	HasFulfillments *bool
	Label           string
	Limit           int
	Offset          int
}

// OrdersResponse represents order summaries.
type OrdersResponse struct {
	Orders []OrderSummary `json:"orders"`
}

// OrderSummary represents order summary information.
type OrderSummary struct {
	ID                      string    `json:"id"`
	CreatedAt               Timestamp `json:"created_at"`
	Label                   string    `json:"label"`
	TotalCents              int       `json:"total_cents"`
	ShippingMethod          string    `json:"shipping_method"`
	LatestFulfillmentStatus *string   `json:"latest_fulfillment_status"`
}

// OrderDetailsResponse represents detailed order response.
type OrderDetailsResponse struct {
	Order OrderDetails `json:"order"`
}

// OrderDetails represents detailed order data.
type OrderDetails struct {
	OrderSummary
	BuyerID         string             `json:"buyer_id"`
	ShippingAddress Address            `json:"shipping_address"`
	Payment         OrderPayment       `json:"payment"`
	Fulfillments    []OrderFulfillment `json:"fulfillments"`
	Items           []OrderItem        `json:"items"`
}

// OrderPayment represents payment details.
type OrderPayment struct {
	SubtotalCents int `json:"subtotal_cents"`
	ShippingCents int `json:"shipping_cents"`
	TotalCents    int `json:"total_cents"`
	FeeCents      int `json:"fee_cents"`
	NetCents      int `json:"net_cents"`
}

// OrderFulfillment represents order fulfillment.
type OrderFulfillment struct {
	Status              *string    `json:"status"`
	TrackingCompany     *string    `json:"tracking_company"`
	TrackingNumber      *string    `json:"tracking_number"`
	TrackingURL         *string    `json:"tracking_url"`
	InTransitAt         *Timestamp `json:"in_transit_at"`
	EstimatedDeliveryAt *Timestamp `json:"estimated_delivery_at"`
	DeliveredAt         *Timestamp `json:"delivered_at"`
}

// OrderItem represents an order item.
type OrderItem struct {
	TCGSKU      *int    `json:"tcgsku"`
	ProductID   string  `json:"product_id"`
	ProductType string  `json:"product_type"`
	Product     Product `json:"product"`
	Quantity    int     `json:"quantity"`
	PriceCents  int     `json:"price_cents"`
}

// OrderFulfillmentRequest represents a fulfillment update request.
type OrderFulfillmentRequest struct {
	Status              *string    `json:"status"`
	TrackingCompany     *string    `json:"tracking_company"`
	TrackingNumber      *string    `json:"tracking_number"`
	TrackingURL         *string    `json:"tracking_url"`
	InTransitAt         *Timestamp `json:"in_transit_at"`
	EstimatedDeliveryAt *Timestamp `json:"estimated_delivery_at"`
	DeliveredAt         *Timestamp `json:"delivered_at"`
}

// OrderFulfillmentResponse represents fulfillment response.
type OrderFulfillmentResponse struct {
	Fulfillment OrderFulfillment `json:"fulfillment"`
}

// OrderReportsResponse represents order reports.
type OrderReportsResponse struct {
	Reports []OrderReport `json:"reports"`
}

// OrderReport represents an order report.
type OrderReport struct {
	ReportID            string              `json:"report_id"`
	OrderID             string              `json:"order_id"`
	OrderReportedIssues OrderReportedIssues `json:"order_reported_issues"`
}

// OrderReportedIssues represents reported issues.
type OrderReportedIssues struct {
	Comment                   *string                    `json:"comment"`
	CreatedAt                 Timestamp                  `json:"created_at"`
	ProposedRemediationMethod *string                    `json:"proposed_remediation_method"`
	ReporterRole              string                     `json:"reporter_role"`
	IsNonDeliveryReport       bool                       `json:"is_nondelivery_report"`
	Rescinded                 bool                       `json:"rescinded"`
	Items                     []OrderReportedItem        `json:"items"`
	Remediations              []OrderReportedRemediation `json:"remediations"`
	Charges                   []OrderReportedCharge      `json:"charges"`
}

// OrderReportedItem represents a reported order item.
type OrderReportedItem struct {
	OrderItemID string `json:"order_item_id"`
	Quantity    int    `json:"quantity"`
}

// OrderReportedRemediation represents a remediation entry.
type OrderReportedRemediation struct {
	RemediationExpenseCents *int      `json:"remediation_expense_cents"`
	Comment                 *string   `json:"comment"`
	CreatedAt               Timestamp `json:"created_at"`
}

// OrderReportedCharge represents a charge entry.
type OrderReportedCharge struct {
	SellerChargeCents *int    `json:"seller_charge_cents"`
	PayoutID          *string `json:"payout_id"`
}

// Webhook represents a webhook registration.
type Webhook struct {
	ID          string `json:"id"`
	Topic       string `json:"topic"`
	CallbackURL string `json:"callback_url"`
}

// WebhooksResponse represents webhooks list response.
type WebhooksResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

// WebhookRegisterRequest represents webhook register request.
type WebhookRegisterRequest struct {
	Topic       string `json:"topic"`
	CallbackURL string `json:"callback_url"`
}

// CardInfoRequest represents a card info request.
type CardInfoRequest struct {
	CardNames []string `json:"card_names"`
}

// CardInfoResponse represents a card info response.
type CardInfoResponse struct {
	Cards    []CardInfo `json:"cards"`
	NotFound []string   `json:"not_found"`
}

// CardInfo represents card metadata.
type CardInfo struct {
	Name              string   `json:"name"`
	SetCode           string   `json:"set_code"`
	SetName           string   `json:"set_name"`
	CardNumber        string   `json:"card_number"`
	Rarity            string   `json:"rarity"`
	FromPriceCents    *int     `json:"from_price_cents"`
	QuantityAvailable int      `json:"quantity_available"`
	ReleaseDate       string   `json:"release_date"`
	LegalFormats      []string `json:"legal_formats"`
	FlavorName        *string  `json:"flavor_name"`
	Layout            *string  `json:"layout"`
	IsToken           bool     `json:"is_token"`
	PromoTypes        []string `json:"promo_types"`
	Finishes          []string `json:"finishes"`
	Text              *string  `json:"text"`
	ColorIdentity     []string `json:"color_identity"`
	EdhrecSaltiness   *string  `json:"edhrecSaltiness"`
	Power             *string  `json:"power"`
	Defense           *string  `json:"defense"`
	ManaCost          *string  `json:"mana_cost"`
	ManaValue         *string  `json:"mana_value"`
}

// JobApplicationRequest represents a job application.
type JobApplicationRequest struct {
	FirstName           string
	LastName            string
	Email               string
	LinkedInURL         string
	GitHubURL           string
	Application         []byte
	ApplicationFilename string
}

// JobApplicationResponse represents a job application response.
type JobApplicationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
