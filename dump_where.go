package mysqldump_some

// DumpTableCondition is an option for selected records.
type DumpTableCondition struct {
	TableName string
	Condition string
}

type DumpMatchCondition struct {
	Column string
	Values []string
}
