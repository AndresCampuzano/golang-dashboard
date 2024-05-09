package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

func (s *PostgresStore) GetEarnings() ([]*Earnings, error) {
	rows, err := s.db.Query(`

	WITH monthly_expenses AS (
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
		SELECT DISTINCT
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
			SUM(CASE WHEN e.currency = 'COP' THEN e.price ELSE 0 END) AS total_cop_expense
		FROM
			expenses e
		GROUP BY
			month
	),
	sales_count AS (
		SELECT
			DATE_TRUNC('month', s.created_at) AS month,
			COUNT(*) AS total_sales_in_month
		FROM
			sales s
		GROUP BY
			DATE_TRUNC('month', s.created_at)
	),
	total_product_variations AS (
	SELECT
		DATE_TRUNC('month', pv.created_at) AS month,
		COUNT(*) AS total_variations
	FROM
		product_variations pv
	GROUP BY
		DATE_TRUNC('month', pv.created_at)
     ),
	-- Generate a distinct list of months from all sources
	distinct_months AS (
		SELECT DISTINCT month FROM (
			SELECT month FROM monthly_expenses
			UNION
			SELECT month FROM monthly_income
			UNION
			SELECT month FROM cop_expenses
			UNION
			SELECT month FROM sales_count
		) AS all_months
	)
	SELECT
		dm.month AS sort_by_month,
		COALESCE(
			(
				SELECT JSON_AGG(jsonb_build_object('currency', me.currency, 'value', me.total_expense))
				FROM (
					SELECT DISTINCT ON (currency) currency, total_expense
					FROM monthly_expenses
					WHERE month = dm.month
					ORDER BY currency
				) AS me
			),
			'[]'
		) AS expenses_summary,
		COALESCE(
			(
				SELECT JSON_AGG(jsonb_build_object(
					'id', aem.id,
					'name', aem.name,
					'price', aem.price,
					'type', aem.type,
					'description', aem.description,
					'currency', aem.currency,
					'created_at', aem.created_at,
					'updated_at', aem.updated_at
				))
				FROM all_expenses_in_month aem
				WHERE DATE_TRUNC('month', aem.created_at) = dm.month
			),
			'[]'
		) AS all_expenses_in_month,
		COALESCE(mi.total_income, 0) AS income,
		COALESCE(ce.total_cop_expense, 0) AS cop_expense,
		CASE WHEN COALESCE(mi.total_income, 0) - COALESCE(ce.total_cop_expense, 0) < 0 THEN 0
		 ELSE COALESCE(mi.total_income, 0) - COALESCE(ce.total_cop_expense, 0)
		END AS earnings,
		COALESCE(sc.total_sales_in_month, 0) AS total_sales_in_month,
		COALESCE(tpv.total_variations, 0) AS total_product_variations_in_month
	FROM
		distinct_months dm
	LEFT JOIN
		monthly_income mi ON dm.month = mi.month
	LEFT JOIN
		cop_expenses ce ON dm.month = ce.month
	LEFT JOIN
		sales_count sc ON dm.month = sc.month
	LEFT JOIN
		total_product_variations tpv ON dm.month = tpv.month
	ORDER BY
		dm.month;

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
		&earning.TotalSalesInMonth,
		&earning.TotalProductVariationsInMonth,
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
