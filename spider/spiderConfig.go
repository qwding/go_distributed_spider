package spider

import (
	"time"
)

const Url = "/base"
const Share = 3
const AllFDefault = "all.txt"
const MatchFDefault = "match.txt"

var DefaultDailTimeout time.Duration = time.Duration(2) * time.Second
