// Package quickbooks implements exporting function to iif format for Quickbooks
package quickbooks

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Package global constants
const (
	ConstErrorModule = "quickbooks"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

var (
	orderFields = [][]string{
		{
			"!TRNS", "TRNSID", "TRNSTYPE", "DATE", "ACCNT", "NAME", "CLASS", "AMOUNT", "DOCNUM",
			"MEMO", "CLEAR", "TOPRINT", "ADDR1", "ADDR2", "ADDR3", "ADDR4", "ADDR5", "DUEDATE",
			"TERMS", "PAID", "PAYMETH", "SHIPVIA", "SHIPDATE", "REP", "FOB", "PONUM", "INVTITLE",
			"INVMEMO", "SADDR1", "SADDR2", "SADDR3", "SADDR4", "SADDR5", "NAMEISTAXABLE",
		},
	}

	dataSeted = [][]interface{}{
		{
			"TRNS", "$_id", "PAYMENT",
			func(record map[string]interface{}) string {
				return utils.InterfaceToTime(record["created_at"]).Format("02/01/06")
			},
			"Orders", "$customer_name", "", "$grand_total", "$increment_id",
			"MEMO", "N", "TOPRINT", "ADDR1", "ADDR2", "ADDR3", "ADDR4", "ADDR5", "DUEDATE",
			"TERMS", "PAID", "PAYMETH", "SHIPVIA", "SHIPDATE", "REP", "FOB", "PONUM", "INVTITLE",
			"INVMEMO", "SADDR1", "SADDR2", "SADDR3", "SADDR4", "SADDR5", "NAMEISTAXABLE",
		},
	}
)

// "$_id", "PAYMENT", "func", "Orders", "$customer_name", "", "$grand_total", "$increment_id", "MEMO", "N", "TOPRINT", "ADDR1", "ADDR2", "ADDR3", "ADDR4", "ADDR5", "DUEDATE", "TERMS", "PAID", "PAYMETH", "SHIPVIA", "SHIPDATE", "REP", "FOB", "PONUM", "INVTITLE", "INVMEMO", "SADDR1", "SADDR2", "SADDR3", "SADDR4", "SADDR5", "NAMEISTAXABLE"
