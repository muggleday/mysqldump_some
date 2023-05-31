package mysqldump_some

import (
	"fmt"
	"golang.org/x/exp/slices"
	"os/exec"
	"strings"
)

func (d *Dumper) Dump() error {
	tables, err := QueryTables(d.db, d.Database)
	if err != nil {
		return err
	}

	for _, table := range tables {
		bytes, err := d.DumpTable(*table.TABLENAME)
		fmt.Println(string(bytes))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return nil
}

func (d *Dumper) DumpTable(tableName string) ([]byte, error) {
	if slices.Contains(d.ExcludedTables, tableName) {
		return nil, nil
	}

	if len(d.IncludedTables) > 0 && !slices.Contains(d.IncludedTables, tableName) {
		return nil, nil
	}

	resultFilepath := fmt.Sprintf("%s/%s.sql", d.OutDir, tableName)

	args := []string{
		"--skip-add-locks",
		"--skip-lock-tables",
		"--create-options",
		"--extended-insert",
		"--add-drop-table",
		"--disable-keys",
		fmt.Sprintf(`--result-file=%s`, resultFilepath),
		d.Database,
		tableName,
	}

	if len(d.DumpOptions) > 0 {
		args = append(args, d.DumpOptions...)
	}

	where, err := d.buildDumpWhere(tableName)
	if err != nil {
		return nil, err
	}

	if where != "" {
		args = append(args, fmt.Sprintf(`--where=%s`, where))
	}

	cmd := exec.Command("mysqldump", args...)
	fmt.Println(cmd.String())
	return cmd.CombinedOutput()
}

func (d *Dumper) buildDumpWhere(tableName string) (string, error) {
	for _, condition := range d.TableConditions {
		if condition.TableName == tableName {
			return condition.Condition, nil
		}
	}

	for _, condition := range d.MatchAnyConditions {
		columnExists, err := IsColumnExistsInTable(d.db, d.Database, tableName, condition.Column)
		if err != nil {
			return "", err
		}

		if columnExists {
			var buf strings.Builder
			for _, value := range condition.Values {
				if buf.Len() > 0 {
					buf.WriteString(" and ")
				}
				buf.WriteString(condition.Column)
				buf.WriteByte('=')
				buf.WriteString(value)
			}
			return buf.String(), nil
		}
	}

	return "", nil
}
