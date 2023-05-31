package main

import (
	_ "github.com/go-sql-driver/mysql"
	"mysqldump_some"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/simking_lucy_427?charset=utf8mb4&parseTime=True&loc=UTC&multiStatements=true"
	dumper, err := mysqldump_some.NewDumper(dsn, "simking_lucy_427")
	if err != nil {
		panic(err)
		return
	}

	dumper.DumpOptions = []string{
		"--user=root",
		"--password=123456",
		"--host=127.0.0.1",
		"--port=3306",
	}

	dumper.ExcludedTables = []string{
		"agent_api_call_histories",
		"amqp_messages",
		"announcements",
		"cities",
		"data_export_tasks",
		"districts",
		"cities",
		"gorp_migrations",
		"notifications",
		"new_agent_push_records",
		"platform_api_call_histories",
		"provinces",
		"purchase_orders",
		"sms_channels",
		"tasks",
		"telecom_reports",
	}

	dumper.TableConditions = []mysqldump_some.DumpTableCondition{
		{
			TableName: "agents",
			Condition: "id = 1 or id = 11",
		},
		{
			TableName: "accounts",
			Condition: "user_id = 9 or user_id = 10",
		},
		{
			TableName: "balance_statistics",
			Condition: "user_id = 9 or user_id = 10",
		},
		{
			TableName: "sim_cards",
			Condition: "top_agent_id = 1 or top_agent_id = 11",
		},
	}

	dumper.MatchAnyConditions = []mysqldump_some.DumpMatchCondition{
		{
			Column: "agent_id",
			Values: []string{"1", "11"},
		},
		{
			Column: "account_id",
			Values: []string{"1", "11"},
		},

		{
			Column: "month",
			Values: []string{"202305", "202306"},
		},
	}

	err = dumper.Dump()
	if err != nil {
		panic(err)
		return
	}

}
