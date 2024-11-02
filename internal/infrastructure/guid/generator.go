package guid

import "github.com/google/uuid"

type Generator string

func (g Generator) Generate() string {
	return uuid.NewString()
}
