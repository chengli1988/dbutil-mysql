package dbutil

import "testing"

func init() {
	InitPool("root", "root", "127.0.0.1", 3306, "demo", "utf8mb4")
}

func TestInitPool(t *testing.T) {
}
