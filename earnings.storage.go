package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

func (s *PostgresStore) GetEarnings() ([]*Earnings, error) {
	rows, err := s.db.Query(`WITH monthly_expenses AS (
		SELECT
			DATE_TRUNC('month', e.created_at) AS month,
			e.currency,
			SUM(e.price) AS total_expense
		FROM
			expenses e
		GROUP BY
			month, e.currency
		),
		all_expenses_in_month AS (
			SELECT
				e.id,
				e.name,
				e.price,
				e.type,
				e.description,
				e.currency,
				e.created_at,
				e.updated_at
			FROM
				expenses e
		),
		monthly_income AS (
			SELECT
				DATE_TRUNC('month', pv.created_at) AS month,
				SUM(pv.price) AS total_income
			FROM
				product_variations pv
			GROUP BY
				month
		),
		cop_expenses AS (
			SELECT
				DATE_TRUNC('month', e.created_at) AS month,
				SUM(e.price) AS total_cop_expense
			FROM
				expenses e
			WHERE
				e.currency = 'COP'
			GROUP BY
				month
		)
		SELECT
			me.month AS sort_by_month,
			(SELECT json_agg(expenses_summary)
			 FROM (SELECT DISTINCT ON (currency) jsonb_build_object('currency', currency, 'value', total_expense) AS expenses_summary
				   FROM monthly_expenses
				   WHERE month = me.month
				   ORDER BY currency) AS distinct_expenses) AS expenses_summary,
			json_agg(
				jsonb_build_object(
					'id', aem.id,
					'name', aem.name,
					'price', aem.price,
					'type', aem.type,
					'description', aem.description,
					'currency', aem.currency,
					'created_at', aem.created_at,
					'updated_at', aem.updated_at
				)
			) AS all_expenses_in_month,
			COALESCE(mi.total_income, 0) AS income,
			COALESCE(ce.total_cop_expense, 0) AS cop_expense,
			COALESCE(mi.total_income, 0) - COALESCE(ce.total_cop_expense, 0) AS earnings
		FROM
			monthly_expenses me
		LEFT JOIN
			all_expenses_in_month aem ON DATE_TRUNC('month', aem.created_at) = me.month
		LEFT JOIN
			monthly_income mi ON me.month = mi.month
		LEFT JOIN
			cop_expenses ce ON me.month = ce.month
		GROUP BY
			me.month, mi.total_income, ce.total_cop_expense;
	`)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var earnings []*Earnings
	for rows.Next() {
		earning, err := scanIntoEarnings(rows)
		if err != nil {
			return nil, err
		}

		earnings = append(earnings, earning)
	}

	return earnings, nil
}

func scanIntoEarnings(rows *sql.Rows) (*Earnings, error) {
	earning := new(Earnings)
	var expensesSummaryJSON []byte
	var allExpensesInMonthJSON []byte
	err := rows.Scan(
		&earning.SortByMonth,
		&expensesSummaryJSON,
		&allExpensesInMonthJSON,
		&earning.Income,
		&earning.CopExpense,
		&earning.Earnings,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(expensesSummaryJSON, &earning.ExpensesSummary)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling expensesSummaryJSON JSON: %v", err)
	}

	err = json.Unmarshal(allExpensesInMonthJSON, &earning.AllExpensesInMonth)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling allExpensesInMonthJSON JSON: %v", err)
	}

	return earning, nil
}
