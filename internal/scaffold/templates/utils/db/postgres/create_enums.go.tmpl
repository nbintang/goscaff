package main

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type EnumOptions struct {
	Name   string
	Values []string
	Schema string
}

func CreateEnums(opts EnumOptions) func(ctx context.Context, db *gorm.DB) error {
	return func(ctx context.Context, db *gorm.DB) error {
		return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

			schema := opts.Schema
			if schema == "" {
				schema = "public"
			}
			enumName := strings.ToLower(opts.Name)
			quoted := make([]string, len(opts.Values))
			for i, v := range opts.Values {
				quoted[i] = fmt.Sprintf("'%s'", v)
			}

			enumValues := strings.Join(quoted, ", ")

			query := fmt.Sprintf(`
							DO $$ BEGIN
								CREATE TYPE %s.%s AS ENUM (%s);
							EXCEPTION
								WHEN duplicate_object THEN null;
							END $$;`, schema, enumName, enumValues)

			if err := tx.Exec(query).Error; err != nil {
				return err
			}

			return nil
		})
	}
}