package GlobalDNS

type Error struct {
	Result string `json:"result"`
	Error  struct {
		ErrorCode  string `json:"errorCode"`
		Message    string `json:"message"`
		DevMessage string `json:"devMessage"`
	} `json:"error"`
}

type Sort struct {
	Sorted   bool `json:"sorted"`
	Unsorted bool `json:"unsorted"`
	Empty    bool `json:"empty"`
}

type Pageable struct {
	Sort       Sort `json:"sort"`
	Offset     int  `json:"offset"`
	PageNumber int  `json:"pageNumber"`
	PageSize   int  `json:"pageSize"`
	Paged      bool `json:"paged"`
	Unpaged    bool `json:"unpaged"`
}

type DomainCreateRequest struct {
	Name     string `json:"name"`
	Comments string `json:"comments"`
}

type Domain struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	CompleteYn bool   `json:"completeYn"`
	DnssecYn   bool   `json:"dnssecYn"`
}

type Domains struct {
	Content       []Domain `json:"content"`
	TotalElements int      `json:"totalElements"`
	TotalPages    int      `json:"totalPages"`
	First         bool     `json:"first"`
	Last          bool     `json:"last"`
	Empty         bool     `json:"empty"`
	Size          int      `json:"size"`
	Number        int      `json:"number"`
	Pageable      Pageable `json:"pageable"`
	Sort          Sort     `json:"sort"`
}

type RecordCreateRequest struct {
	Host    string `json:"host"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
}

type RecordUpdateRequest struct {
	Id      int    `json:"id"`
	Host    string `json:"host"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
}

type Record struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Host         string `json:"host"`
	Type         string `json:"type"`
	Content      string `json:"content"`
	Ttl          int    `json:"ttl"`
	AliasId      int    `json:"aliasId"`
	AliasYn      bool   `json:"aliasYn"`
	LbId         int    `json:"lbId"`
	LbYn         bool   `json:"lbYn"`
	LbRegionCode string `json:"lbRegionCode"`
	LbPlatform   string `json:"lbPlatform"`
	DelYn        bool   `json:"delYn"`
	DomainName   string `json:"domainName"`
	CreateDate   int    `json:"createDate"`
	ModifiedDate int    `json:"modifiedDate"`
	ApplyYn      bool   `json:"applyYn"`
	DefaultYn    bool   `json:"defaultYn"`
}

type Records struct {
	Content       []Record `json:"content"`
	TotalElements int      `json:"totalElements"`
	TotalPages    int      `json:"totalPages"`
	First         bool     `json:"first"`
	Last          bool     `json:"last"`
	Empty         bool     `json:"empty"`
	Size          int      `json:"size"`
	Number        int      `json:"number"`
	Pageable      Pageable `json:"pageable"`
	Sort          Sort     `json:"sort"`
}
