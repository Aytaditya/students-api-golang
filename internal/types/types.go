package types

type Student struct {
	Id    int64
	Name  string `validate:"required"`
	Email string `validate:"required"`
	Age   int    `validate:"gte=0,lte=130,required"`
}

// gte = greater than equal to
// lte = less than equal to
