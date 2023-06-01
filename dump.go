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
		"--create-options",   // 在 CREATE TABLE 语句中包括所有MySQL特性选项
		"--extended-insert",  // 使用具有多个VALUES列的INSERT语法。这样使导出文件更小，并加速导入时的速度。默认为打开状态
		"--add-drop-table",   // 每个数据表创建之前添加drop数据表语句
		"--disable-keys",     // 这样可以更快地导入dump出来的文件，因为它是在插入所有行后创建索引的
		"--skip-lock-tables", // 不锁定所有表
		"--set-gtid-purged=OFF",
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
					buf.WriteString(" or ")
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
