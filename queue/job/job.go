package job

import (
	"context"
	"fmt"
	"time"
)

type Job interface {
	fmt.Stringer
	Describe() []string
	RunAt() time.Time
	Run(context.Context) ([]Job, error)
}
