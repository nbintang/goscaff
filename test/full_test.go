package test

import "testing"

func TestFull_Postgres_Builds(t *testing.T) {
	GenerateAndBuild(t, Case{
		Name:   "full-postgres",
		Preset: "full",
		DB:     "postgres",
	})
}

func TestFull_MySQL_Builds(t *testing.T) {
	GenerateAndBuild(t, Case{
		Name:   "full-mysql",
		Preset: "full",
		DB:     "mysql",
	})
}
