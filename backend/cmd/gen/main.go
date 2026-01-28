package main

import (
	"os"

	"github.com/MisakaTAT/GTerm/backend/dal/model"
	"gorm.io/gen"
)

func models() []any {
	return []any{
		model.Connection{},
		model.Credential{},
		model.Group{},
		model.Metadata{},
	}
}

func main() {
	_, err := os.Stat("main.go")
	if err != nil {
		panic("请在项目根目录执行 go run backend/cmd/gen/main.go")
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./backend/dal/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
	})

	g.ApplyBasic(models()...)
	g.Execute()
}
