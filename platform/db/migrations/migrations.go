package migrations

import "embed"

//go:embed postgres/*.sql
var EmbedMigrations embed.FS
