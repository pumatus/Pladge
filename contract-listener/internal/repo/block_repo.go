package repo

var last uint64 = 0

func GetLastBlock() uint64 {
	return last
}

func SaveLastBlock(b uint64) {
	if b > last {
		last = b
	}
}
