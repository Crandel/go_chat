package migrations

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samonzeweb/godb"

	lg "github.com/Crandel/go_chat/internal/logging"
)

var log = lg.InitLogger()

const (
	migrationBegin = `BEGIN TRANSACTION;
`

	migrationSuffix = `
COMMIT;`
)

func RunMigrations(db *godb.DB, migrationFolder string) error {
	files, err := os.ReadDir(migrationFolder)
	if err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool {
		return strings.Compare(files[i].Name(), files[j].Name()) < 0
	})

	var sb strings.Builder
	count := 0
	for _, f := range files {
		fileName := f.Name()
		ext := path.Ext(fileName)
		if strings.ToLower(ext) != ".sql" {
			continue
		}
		var b []byte
		b, err = os.ReadFile(filepath.Join(migrationFolder, fileName))
		if err != nil {
			return err
		}
		if count > 0 {
			sb.WriteString(" ")
		}

		_, err = sb.WriteString(migrationBegin + string(b) + migrationSuffix + " ")

		if err != nil {
			return err
		}

		count++
	}
	if sb.Len() > 0 {
		sqlQueryRaw := sb.String()
		sqlQueryRaw = strings.ReplaceAll(sqlQueryRaw, "\n", " ")
		log.Debugln(sqlQueryRaw)
		_, err = db.CurrentDB().Exec(sqlQueryRaw)
	}
	return err
}
