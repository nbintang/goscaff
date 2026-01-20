package test

import "testing"

func TestBase_Postgres_Builds(t *testing.T) {
	GenerateAndBuild(t, Case{
		Name:   "base-postgres",
		Preset: "base",
		DB:     "postgres",
	})
}

func TestBase_MySQL_Builds(t *testing.T) {
	GenerateAndBuild(t, Case{
		Name:   "base-mysql",
		Preset: "base",
		DB:     "mysql",
	})
}
