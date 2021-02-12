package mt5client

const (
	DealActionBuy          = 0
	DealActionSell         = 1
	DealActionBalance      = 2
	DealEntryIn            = 0
	DealEntryOut           = 1
	TradeActionDealerFirst = 200
	OrderFillingFillFOK    = 0
	DealerRequestDone      = 10009
)

type ApiData struct {
	AppID       uint64  `json:"AppID,string"`
	ID          uint64  `json:"ID,string"`
	ValueInt    int64   `json:"ValueInt,string"`
	ValueUInt   uint64  `json:"ValueUInt,string"`
	ValueDouble float64 `json:"ValueDouble,string"`
}

type Order struct {
	Order            uint64    `json:"Order,string" db:"Order"`
	ExternalID       string    `json:"ExternalID" db:"ExternalID"`
	Login            string    `json:"Login" db:"Login"`
	Dealer           uint64    `json:"Dealer,string" db:"Dealer"`
	Symbol           string    `json:"Symbol" db:"Symbol"`
	Digits           uint8     `json:"Digits,string" db:"Digits"`
	DigitsCurrency   uint8     `json:"DigitsCurrency,string" db:"DigitsCurrency"`
	ContractSize     float64   `json:"ContractSize,string" db:"ContractSize"`
	State            uint8     `json:"State,string" db:"State"`
	Reason           uint8     `json:"Reason,string" db:"Reason"`
	TimeSetup        int64     `json:"TimeSetup,string" db:"TimeSetup"`
	TimeExpiration   int64     `json:"TimeExpiration,string" db:"TimeExpiration"`
	TimeDone         int64     `json:"TimeDone,string" db:"TimeDone"`
	TimeSetupMsc     int64     `json:"TimeSetupMsc,string" db:"TimeSetupMsc"`
	TimeDoneMsc      int64     `json:"TimeDoneMsc,string" db:"TimeDoneMsc"`
	ModifyFlags      uint8     `json:"ModifyFlags,string" db:"ModifyFlags"`
	Type             uint8     `json:"Type,string" db:"Type"`
	TypeFill         uint8     `json:"TypeFill,string" db:"TypeFill"`
	TypeTime         uint8     `json:"TypeTime,string" db:"TypeTime"`
	PriceOrder       float64   `json:"PriceOrder,string" db:"PriceOrder"`
	PriceTrigger     float64   `json:"PriceTrigger,string" db:"PriceTrigger"`
	PriceCurrent     float64   `json:"PriceCurrent,string" db:"PriceCurrent"`
	PriceSL          float64   `json:"PriceSL,string" db:"PriceSL"`
	PriceTP          float64   `json:"PriceTP,string" db:"PriceTP"`
	VolumeInitial    uint64    `json:"VolumeInitial,string" db:"VolumeInitial"`
	VolumeInitialExt uint64    `json:"VolumeInitialExt,string" db:"VolumeInitialExt"`
	VolumeCurrent    uint64    `json:"VolumeCurrent,string" db:"VolumeCurrent"`
	VolumeCurrentExt uint64    `json:"VolumeCurrentExt,string" db:"VolumeCurrentExt"`
	ExpertID         uint64    `json:"ExpertID,string" db:"ExpertID"`
	PositionID       uint64    `json:"PositionID,string" db:"PositionID"`
	PositionByID     uint64    `json:"PositionByID,string" db:"PositionByID"`
	Comment          string    `json:"Comment" db:"Comment"`
	RateMargin       float64   `json:"RateMargin,string" db:"RateMargin"`
	ActivationMode   uint8     `json:"ActivationMode,string" db:"ActivationMode"`
	ActivationTime   int64     `json:"ActivationTime,string" db:"ActivationTime"`
	ActivationPrice  float64   `json:"ActivationPrice,string" db:"ActivationPrice"`
	ActivationFlags  uint8     `json:"ActivationFlags,string" db:"ActivationFlags"`
	ApiData          []ApiData `json:"ApiData"`
}

type OrdersTotalResponse struct {
	Total int
}

type OrdersResponse struct {
	Orders []Order
}

type Deal struct {
	Deal            uint64    `json:"Deal,string" db:"Deal"`
	ExternalID      string    `json:"ExternalID" db:"ExternalID"`
	Login           uint64    `json:"Login,string" db:"Login"`
	Dealer          uint64    `json:"Dealer,string" db:"Dealer"`
	Order           uint64    `json:"Order,string" db:"Order"`
	Action          uint8     `json:"Action,string" db:"Action"`
	Entry           uint8     `json:"Entry,string" db:"Entry"`
	Reason          uint8     `json:"Reason,string" db:"Reason"`
	Digits          uint8     `json:"Digits,string" db:"Digits"`
	DigitsCurrency  uint8     `json:"DigitsCurrency,string" db:"DigitsCurrency"`
	ContractSize    float32   `json:"ContractSize,string" db:"ContractSize"`
	Time            int64     `json:"Time,string" db:"Time"`
	TimeMsc         int64     `json:"TimeMsc,string" db:"TimeMsc"`
	Symbol          string    `json:"Symbol" db:"Symbol"`
	Price           float32   `json:"Price,string" db:"Price"`
	Volume          uint64    `json:"Volume,string" db:"Volume"`
	VolumeExt       uint64    `json:"VolumeExt,string" db:"VolumeExt"`
	Profit          float32   `json:"Profit,string" db:"Profit"`
	Storage         float32   `json:"Storage,string" db:"Storage"`
	Commission      float32   `json:"Commission,string" db:"Commission"`
	Fee             float32   `json:"Fee,string" db:"Fee"`
	RateProfit      float32   `json:"RateProfit,string" db:"RateProfit"`
	RateMargin      float32   `json:"RateMargin,string" db:"RateMargin"`
	ExpertID        uint64    `json:"ExpertID,string" db:"ExpertID"`
	PositionID      uint64    `json:"PositionID,string" db:"PositionID"`
	Comment         string    `json:"Comment" db:"Comment"`
	ProfitRaw       float32   `json:"ProfitRaw,string" db:"ProfitRaw"`
	PricePosition   float32   `json:"PricePosition,string" db:"PricePosition"`
	PriceSL         float32   `json:"PriceSL,string" db:"PriceSL"`
	PriceTP         float32   `json:"PriceTP,string" db:"PriceTP"`
	VolumeClosed    uint64    `json:"VolumeClosed,string" db:"VolumeClosed"`
	VolumeClosedExt uint64    `json:"VolumeClosedExt,string" db:"VolumeClosedExt"`
	TickValue       float32   `json:"TickValue,string" db:"TickValue"`
	TickSize        float32   `json:"TickSize,string" db:"TickSize"`
	Flags           uint8     `json:"Flags,string" db:"Flags"`
	Gateway         string    `json:"Gateway" db:"Gateway"`
	PriceGateway    float32   `json:"PriceGateway,string" db:"PriceGateway"`
	ModifyFlags     uint8     `json:"ModifyFlags,string" db:"ModifyFlags"`
	Value           float32   `json:"Value,string" db:"Value"`
	ApiData         []ApiData `json:"ApiData"`
}

type DealsTotalResponse struct {
	Total int
}

type DealsResponse struct {
	Deals []Deal
}

type Position struct {
	Position         uint64    `json:"Position,string" db:"Position"`
	ExternalID       string    `json:"ExternalID" db:"ExternalID"`
	Login            string    `json:"Login" db:"Login"`
	Dealer           uint64    `json:"Dealer,string" db:"Dealer"`
	Symbol           string    `json:"Symbol" db:"Symbol"`
	Action           uint8     `json:"Action,string" db:"Action"`
	Digits           uint8     `json:"Digits,string" db:"Digits"`
	DigitsCurrency   uint8     `json:"DigitsCurrency,string" db:"DigitsCurrency"`
	Reason           uint8     `json:"Reason,string" db:"Reason"`
	ContractSize     float32   `json:"ContractSize,string" db:"ContractSize"`
	TimeCreate       int64     `json:"TimeCreate,string" db:"TimeCreate"`
	TimeUpdate       int64     `json:"TimeUpdate,string" db:"TimeUpdate"`
	TimeCreateMsc    int64     `json:"TimeCreateMsc,string" db:"TimeCreateMsc"`
	TimeUpdateMsc    int64     `json:"TimeUpdateMsc,string" db:"TimeUpdateMsc"`
	ModifyFlags      uint8     `json:"ModifyFlags,string" db:"ModifyFlags"`
	PriceOpen        float32   `json:"PriceOpen,string" db:"PriceOpen"`
	PriceCurrent     float32   `json:"PriceCurrent,string" db:"PriceCurrent"`
	PriceSL          float32   `json:"PriceSL,string" db:"PriceSL"`
	PriceTP          float32   `json:"PriceTP,string" db:"PriceTP"`
	Volume           uint64    `json:"Volume,string" db:"Volume"`
	VolumeExt        uint64    `json:"VolumeExt,string" db:"VolumeExt"`
	Profit           float32   `json:"Profit,string" db:"Profit"`
	Storage          float32   `json:"Storage,string" db:"Storage"`
	RateProfit       float32   `json:"RateProfit,string" db:"RateProfit"`
	RateMargin       float32   `json:"RateMargin,string" db:"RateMargin"`
	ExpertID         uint64    `json:"ExpertID,string" db:"ExpertID"`
	ExpertPositionID uint64    `json:"ExpertPositionID,string" db:"ExpertPositionID"`
	Comment          string    `json:"Comment" db:"Comment"`
	ActivationMode   uint8     `json:"ActivationMode,string" db:"ActivationMode"`
	ActivationTime   int64     `json:"ActivationTime,string" db:"ActivationTime"`
	ActivationPrice  float32   `json:"ActivationPrice,string" db:"ActivationPrice"`
	ActivationFlags  uint8     `json:"ActivationFlags,string" db:"ActivationFlags"`
	ApiData          []ApiData `json:"ApiData"`
}

type PositionsTotalResponse struct {
	Total int
}

type PositionsResponse struct {
	Positions []Position
}

type ClientsResponse struct {
	Ids []string
}

type User struct {
	Login                  string  `json:"Login"`
	Group                  string  `json:"Group"`
	CertSerialNumber       uint64  `json:"CertSerialNumber,string"`
	Rights                 uint16  `json:"Rights,string"`
	MQID                   string  `json:"MQID"`
	Registration           int64   `json:"Registration,string"`
	LastAccess             int64   `json:"LastAccess,string"`
	LastPassChange         int64   `json:"LastPassChange,string"`
	LastIP                 string  `json:"LastIP"`
	Name                   string  `json:"Name"`
	FirstName              string  `json:"FirstName"`
	LastName               string  `json:"LastName"`
	MiddleName             string  `json:"MiddleName"`
	Company                string  `json:"Company"`
	Account                string  `json:"Account"`
	Country                string  `json:"Country"`
	Language               int     `json:"Language,string"`
	City                   string  `json:"City"`
	State                  string  `json:"State"`
	ZipCode                string  `json:"ZipCode"`
	Address                string  `json:"Address"`
	Phone                  string  `json:"Phone"`
	Email                  string  `json:"Email"`
	ID                     string  `json:"ID"`
	Status                 string  `json:"Status"`
	Comment                string  `json:"Comment"`
	Color                  uint32  `json:"Color,string"`
	PhonePassword          string  `json:"PhonePassword"`
	Leverage               int     `json:"Leverage,string"`
	Agent                  uint64  `json:"Agent,string"`
	Balance                float64 `json:"Balance,string"`
	Credit                 float64 `json:"Credit,string"`
	InterestRate           float64 `json:"InvestRate,string"`
	CommissionDaily        float64 `json:"CommissionDaily,string"`
	CommissionMonthly      float64 `json:"CommissionMonthly,string"`
	CommissionAgentDaily   float64 `json:"CommissionAgentDaily,string"`
	CommissionAgentMonthly float64 `json:"CommissionAgentMonthly,string"`
	BalancePrevDay         float64 `json:"BalancePrevDay,string"`
	BalancePrevMonth       float64 `json:"BalancePrevMonth,string"`
	EquityPrevDay          float64 `json:"EquityPrevDay,string"`
	EquityPrevMonth        float64 `json:"EquityPrevMonth,string"`
	TradeAccounts          string  `json:"TradeAccounts"`
}

type UserAccount struct {
	Login       string  `json:"Login"`
	Balance     float32 `json:"Balance,string"`
	Margin      float32 `json:"Margin,string"`
	MarginFree  float32 `json:"MarginFree,string"`
	MarginLevel float32 `json:"MarginLevel,string"`
	Equity      float32 `json:"Equity,string"`
}

type DealerUpdatesResult struct {
	Id                 string  `json:"ID"`
	DealId             string  `json:"DealID"`
	OrderId            string  `json:"OrderID"`
	PositionExternalId string  `json:"PositionExternalID"`
	Retcode            uint16  `json:"Retcode,string"`
	ExternalRetcode    uint16  `json:"ExternalRetcode,string"`
	Volume             uint64  `json:"Volume,string"`
	VolumeExt          uint64  `json:"VolumeExt,string"`
	Price              float64 `json:"Price,string"`
	PriceGateway       float64 `json:"PriceGateway,string"`
	TickBid            float64 `json:"TickBid,string"`
	TickAsk            float64 `json:"TickAsk,string"`
	TickLast           float64 `json:"TickLast,string"`
	Comment            string  `json:"Comment"`
	Flags              uint16  `json:"Flags,string"`
}

type DealerUpdatesAnswer struct {
	Id                      uint64  `json:"ID,string"`
	IdClient                uint64  `json:"IDClient,string"`
	Login                   string  `json:"Login"`
	SourceLogin             string  `json:"SourceLogin"`
	ExternalAccount         string  `json:"ExternalAccount"`
	Ip                      string  `json:"IP"`
	Group                   string  `json:"Group"`
	Simbol                  string  `json:"Symbol"`
	Digits                  uint8   `json:"Digits,string"`
	Action                  uint16  `json:"Action,string"`
	TimeExpiration          int64   `json:"TimeExpiration,string"`
	Type                    uint8   `json:"Type,string"`
	TypeFill                uint8   `json:"TypeFill,string"`
	TypeTime                int64   `json:"TypeTime,string"`
	Flags                   uint16  `json:"Flags,string"`
	Volume                  uint64  `json:"Volume,string"`
	VolumeExt               uint64  `json:"VolumeExt,string"`
	Order                   uint64  `json:"Order,string"`
	OrderExternalId         string  `json:"OrderExternalID"`
	PriceOrder              float64 `json:"PriceOrder,string"`
	PriceTrigger            float64 `json:"PriceTrigger,string"`
	PriceSL                 float64 `json:"PriceSL,string"`
	PriceTP                 float64 `json:"PriceTP,string"`
	PriceDeviation          uint64  `json:"PriceDeviation,string"`
	PriceDeviationTop       float64 `json:"PriceDeviationTop,string"`
	PriceDeviationTopBottom float64 `json:"PriceDeviationBottom,string"`
	Position                uint64  `json:"Position,string"`
	PositionExternalId      string  `json:"PositionExternalID"`
	PositionBy              string  `json:"PositionBy"`
	PositionByExternalID    string  `json:"PositionByExternalID"`
	Comment                 string  `json:"Comment"`
	ResultRetcode           uint16  `json:"ResultRetcode,string"`
	ResultRetcodeDealer     uint64  `json:"ResultRetcodeDealer,string"`
	ResultDeal              uint64  `json:"ResultDeal,string"`
	ResultOrder             uint64  `json:"ResultOrder,string"`
	ResultVolume            uint64  `json:"ResultVolume,string"`
	ResultVolumeExt         uint64  `json:"ResultVolumeExt,string"`
	ResultPrice             float64 `json:"ResultPrice,string"`
	ResultDealerBid         float64 `json:"ResultDealerBid,string"`
	ResultDealerAsk         float64 `json:"ResultDealerAsk,string"`
	ResultDealerLast        float64 `json:"ResultDealerLast,string"`
	ResultMarketBid         float64 `json:"ResultMarketBid,string"`
	ResultMarketAsk         float64 `json:"ResultMarketAsk,string"`
	ResultMarketLast        float64 `json:"ResultMarketLast,string"`
	ResultComment           string  `json:"ResultComment"`
}

type DealerUpdates struct {
	Result *DealerUpdatesResult `json:"result"`
	Answer *DealerUpdatesAnswer `json:"answer"`
}
