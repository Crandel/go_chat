package migrations

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samonzeweb/godb"
)

const (
	migrationFolder = "migrations"

	migrationBegin = `BEGIN TRANSACTION;
`

	migrationSuffix = `
COMMIT;`
)

func RunMigrations(db *godb.DB) error {
	files, err := ioutil.ReadDir(migrationFolder)
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

		b, err := ioutil.ReadFile(filepath.Join(migrationFolder, fileName))
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
		fmt.Printf("\nRawSQL\n%#v\n", sqlQueryRaw)
		sqlQueryRaw = strings.ReplaceAll(sqlQueryRaw, "\n", " ")
		fmt.Printf("\nSQL\n%#v\n", sqlQueryRaw)
		_, err = db.CurrentDB().Exec(sqlQueryRaw)
	}
	return err
}
