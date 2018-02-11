package bwpool

/* Container */
type PoolData struct {
	Bogus   Bogus
	Workers Workers
}

type Bogus struct {
	Bogus int `json:"Bogus"`
}

type WorkerResponse struct {
	UserName string   `json:userName`
	Total    int64    `json:total`
	Workers  []Worker `json:"worker"`
}

/* name:Worker */
type Workers map[string]Worker

/*
 {
	"name": "worker1", // miner's name
	"hashrate": 123456.123, // miner's calculated
	"accepted": 123456, // accepted
	"rejected": 1234 // Reject count
}
*/
type Worker struct {
	Name     string  `json:"name,omit"`
	HashRate float64 `json:"hashrate"`
	Accepted int64   `json:"accepted"`
	Rejected int64   `json:"rejected"`
}
