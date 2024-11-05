//go:generate mockgen -source ${GOFILE} -package mock -destination mock/generator_mock.go

package guid

import "github.com/google/uuid"

type IGenerator interface {
	Generate() string
}

type Generator string

func (g Generator) Generate() string {
	return uuid.NewString()
}
