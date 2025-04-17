package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type Expense struct {
	Id          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Category    string    `json:"category"`
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage expense-tracker <operation> --flags")
	}
	expenses := readFile("expense.json")

	switch os.Args[1] {
	case "add":
		var description string
		var amount int
		var category string

		pflag.StringVar(&description, "description", "", "Test")
		pflag.IntVar(&amount, "amount", 0, "Amount Test")
		pflag.StringVar(&category, "category", "", "category Test")
		pflag.Parse()

		if description == "" || amount <= 0 || category == "" {
			fmt.Println("Description or amount or category invalid")
			return
		}

		add(expenses, description, amount, category)
	case "list":
		var category string
		pflag.StringVar(&category, "category", "", "category Test")
		pflag.Parse()
		if category == "" {
			list(expenses)
			return
		}
		listByCategory(expenses, category)

	case "summary":
		var month int
		var category string

		pflag.IntVar(&month, "month", 0, "test")
		pflag.StringVar(&category, "category", "", "category Test")
		pflag.Parse()

		if month < 0 || month > 12 {
			fmt.Println("Month is not valid")
			return
		}
		//Burası tam olarak doğru değil
		if month == 0 && category == "" {
			summary(expenses)
			return
		} else if category != "" {
		}
		if month != 0 && category == "" {
			summaryWithMonth(expenses, month)
		} else if month == 0 && category != "" {
			summaryWithCategory(expenses, category)
		} else {
			summaryWithMonthAndCategory(expenses, month, category)
		}
	case "delete":
		var id int
		pflag.IntVar(&id, "id", 0, "Id Test")
		pflag.Parse()

		if id <= 0 {
			fmt.Println("Id is not valid")
			return
		}
		delete(expenses, id)
	}
}

func add(expenses []Expense, description string, amount int, category string) {
	newExpense := Expense{}
	if len(expenses) == 0 {
		newExpense.Id = 1
	} else {
		newExpense.Id = expenses[len(expenses)-1].Id + 1
	}
	newExpense.Amount = amount
	newExpense.Description = description
	newExpense.Date = time.Now()
	newExpense.Category = category

	expenses = append(expenses, newExpense)
	writeFile(expenses)
	fmt.Printf("Expense added successfully (ID: %d)\n", newExpense.Id)
}

func summary(expenses []Expense) {
	var total int
	for _, e := range expenses {
		total += e.Amount
	}
	fmt.Printf("Total expenses: $%d", total)
}
func summaryWithMonth(expenses []Expense, month int) {

	total := 0
	for _, e := range expenses {
		expensesMonth := e.Date.Month()
		if int(expensesMonth) == month {
			total += e.Amount
		}
	}
	fmt.Printf("Total expenses for %s: $%d", time.Month(month), total)
}

func summaryWithCategory(expenses []Expense, category string) {
	total := 0
	for _, e := range expenses {
		if e.Category == category {
			total += e.Amount
		}
	}
	fmt.Printf("Total expenses for %s: $%d", category, total)
}
func summaryWithMonthAndCategory(expenses []Expense, month int, category string) {
	total := 0
	for _, e := range expenses {
		expensesMonth := e.Date.Month()
		if int(expensesMonth) == month && e.Category == category {
			total += e.Amount
		}
	}
	fmt.Printf("Total expenses for %s and category %s: $%d", time.Month(month), category, total)
}

func delete(expenses []Expense, id int) {
	index := -1
	for i, e := range expenses {
		if e.Id == id {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Expenses id not found")
		return
	}
	expenses = append(expenses[:index], expenses[index+1:]...)
	err := writeFile(expenses)
	if err != nil {
		fmt.Println("Error : ", err)
		return
	}
	fmt.Println("Expense deleted successfully")
}

func readFile(Filename string) []Expense {
	var expenses []Expense
	context, err := os.ReadFile(Filename)
	if err != nil {
		fmt.Println("Error :", err)
		return nil
	}
	if len(context) == 0 {
		fmt.Println("file is empty")
		return nil
	}
	err = json.Unmarshal(context, &expenses)
	if err != nil {
		fmt.Println("json unmarshal error : ", err)
		return nil
	}
	return expenses
}

func writeFile(expenses []Expense) error {
	arr, err := json.Marshal(expenses)
	if err != nil {
		return fmt.Errorf("Marshall error")

	}
	err = os.WriteFile("expense.json", arr, 0644)
	if err != nil {
		return fmt.Errorf("File writting error")
	}
	return nil
}

func list(expenses []Expense) {
	if len(expenses) != 0 {
		fmt.Printf("ID           Date          Description       Category    Amount\n")

		for _, expense := range expenses {
			fmt.Printf("%d          %s          %s         %s         $%d\n",
				expense.Id,
				expense.Date.Format("2006-01-02"),
				expense.Description,
				expense.Category,
				expense.Amount)
		}
	}
}

func listByCategory(expenses []Expense, category string) {
	if len(expenses) != 0 {
		fmt.Printf("ID           Date          Description       Category    Amount\n")
		for _, expense := range expenses {
			if expense.Category == category {
				fmt.Printf("%d          %s          %s         %s         $%d\n",
					expense.Id,
					expense.Date.Format("2006-01-02"),
					expense.Description,
					expense.Category,
					expense.Amount)
			}
		}
	}
}
