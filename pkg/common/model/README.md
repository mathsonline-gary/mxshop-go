# Package model

This package contains common GORM model fields and types.

## Usage

```go
package user

import "github.com/zycgary/mxshop-go/pkg/common/model"

type User struct {
	model.IncrementID
	model.Timestamps
	model.SoftDelete

	Name string `gorm:"type:varchar(255);not null"`

	Tags model.StringList `gorm:"type:json"`
}
```