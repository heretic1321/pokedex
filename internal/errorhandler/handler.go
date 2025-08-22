package errorhandler

import (
	"fmt"
	"log"
)

func Handle(err error) error  {
	fmt.Printf("❌ Error: %v\n", err)
	return err
}

func Must(err error) {
	log.Fatalf("❌ Fatal error: %v", err)
	panic(err)
}
