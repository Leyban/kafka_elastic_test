package model

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type HitsTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type HitsItem struct {
	Index  string        `json:"_index"`
	ID     string        `json:"_id"`
	Score  *int64        `json:"_score"`
	Source DepositHit    `json:"_source"`
	Sort   []interface{} `json:"sort"`
}

type DepositHit struct {
	PaymentMethodCode string  `json:"payment_method_code"`
	UpdatedBy         string  `json:"updated_by"`
	Status            int     `json:"status"`
	Currency          string  `json:"currency"`
	LoginName         string  `json:"login_name"`
	GrossAmount       float64 `json:"gross_amount"`
	NetAmount         float64 `json:"net_amount"`
	MemberID          int     `json:"member_id"`
	UpdatedAt         string  `json:"updated_at"`
	CreatedAt         string  `json:"created_at"`
	CreatedBy         string  `json:"created_by"`
	TransactionID     string  `json:"transaction_id"`
	ChargeAmount      float64 `json:"charge_amount"`
	RefCode           string  `json:"ref_code"`
	ID                int     `json:"id"`
	BankAccountID     int     `json:"bank_account_id"`
}

type Hits struct {
	Total    HitsTotal  `json:"total"`
	Hits     []HitsItem `json:"hits"`
	MaxScore *float64   `json:"max_score"`
}

type Aggregations struct {
	FilterStatus FilterStatus `json:"filter_by_status"`
}

type FilterStatus struct {
	DocCount        int             `json:"doc_count"`
	GroupByMemberID GroupByMemberID `json:"group_by_member_id"`
}

type GroupByMemberID struct {
	DocCountErrorUpperBound int      `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int      `json:"sum_other_doc_count"`
	Buckets                 []Bucket `json:"buckets"`
}

type TotalGrossAmount struct {
	Value float64 `json:"value"`
}

type Bucket struct {
	Key              int              `json:"key"`
	DocCount         int              `json:"doc_count"`
	TotalGrossAmount TotalGrossAmount `json:"total_gross_amount"`
}

type Response struct {
	Took         int          `json:"took"`
	TimedOut     bool         `json:"timed_out"`
	Shards       Shards       `json:"_shards"`
	Hits         Hits         `json:"hits"`
	Aggregations Aggregations `json:"aggregations"`
}
