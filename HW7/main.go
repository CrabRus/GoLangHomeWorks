package main

import (
	"errors"
	"fmt"
	"time"
)

const (
	minAmount float64 = 10.00
	maxAmount float64 = 50000.00
)

// ---------- Інтерфейс ----------

type PaymentProcessor interface {
	Name() string
	Validate(amount float64) error
	CalculateFee(amount float64) float64
	Process(amount float64, system *PaymentSystem) (Receipt, error)
}

// ---------- Основні структури ----------

type Account struct {
	Owner   string
	Balance float64
}

func NewAccount(owner string, balance float64) *Account {
	return &Account{Owner: owner, Balance: balance}
}

func (a *Account) AddFunds(funds float64) error {
	if funds <= 0 {
		return fmt.Errorf("сума не повинна бути менше за 0")
	}
	a.Balance += funds
	return nil
}

type Receipt struct {
	Method       string
	Amount       float64
	Fee          float64
	TotalDebited float64
	Success      bool
	Message      string
	Time         time.Time
}

// ---------- Кредитна картка ----------

type CreditCard struct{}

func (c CreditCard) Name() string { return "Кредитна картка" }

func (c CreditCard) Validate(amount float64) error {
	if amount < minAmount || amount > maxAmount {
		return fmt.Errorf("сума повинна бути від %.2f до %.2f грн", minAmount, maxAmount)
	}
	return nil
}

func (c CreditCard) CalculateFee(amount float64) float64 {
	return amount * 0.015
}

func (c CreditCard) Process(amount float64, system *PaymentSystem) (Receipt, error) {
	acct := system.Account
	var r Receipt
	r.Method = c.Name()
	r.Amount = amount
	r.Time = time.Now()

	err := c.Validate(amount)
	if err != nil {
		r.Message = err.Error()
		system.AddReceipt(r)
		return r, err
	}

	fee := c.CalculateFee(amount)
	total := amount + fee
	if acct.Balance < total {
		err = errors.New("недостатньо коштів на рахунку")
		r.Message = err.Error()
		system.AddReceipt(r)
		return r, err
	}

	acct.Balance -= total
	r.Fee = fee
	r.TotalDebited = total
	r.Success = true
	r.Message = "Платіж успішно оброблено!"
	system.AddReceipt(r)
	return r, nil
}

// ---------- PayPal ----------

type PayPal struct{}

func (p PayPal) Name() string { return "PayPal" }

func (p PayPal) Validate(amount float64) error {
	if amount < minAmount || amount > maxAmount {
		return fmt.Errorf("сума повинна бути від %.2f до %.2f грн", minAmount, maxAmount)
	}
	return nil
}

func (p PayPal) CalculateFee(amount float64) float64 {
	return amount * 0.035
}

func (p PayPal) Process(amount float64, system *PaymentSystem) (Receipt, error) {
	acct := system.Account
	var r Receipt
	r.Method = p.Name()
	r.Amount = amount
	r.Time = time.Now()

	err := p.Validate(amount)
	if err != nil {
		r.Message = err.Error()
		system.AddReceipt(r)
		return r, err
	}

	fee := p.CalculateFee(amount)
	total := amount + fee
	if acct.Balance < total {
		err = errors.New("недостатньо коштів на рахунку")
		r.Message = err.Error()
		system.AddReceipt(r)
		return r, err
	}

	acct.Balance -= total
	r.Fee = fee
	r.TotalDebited = total
	r.Success = true
	r.Message = "Платіж успішно оброблено!"
	system.AddReceipt(r)
	return r, nil
}

// ---------- Готівка ----------

type Cash struct{}

func (c Cash) Name() string { return "Готівка" }

func (c Cash) Validate(amount float64) error {
	if amount < minAmount || amount > maxAmount {
		return fmt.Errorf("сума повинна бути від %.2f до %.2f грн", minAmount, maxAmount)
	}
	return nil
}

func (c Cash) CalculateFee(amount float64) float64 { return 0.0 }

func (c Cash) Process(amount float64, system *PaymentSystem) (Receipt, error) {
	var r Receipt
	r.Method = c.Name()
	r.Amount = amount
	r.Time = time.Now()

	if err := c.Validate(amount); err != nil {
		r.Message = err.Error()
		system.AddReceipt(r)
		return r, err
	}

	r.Success = true
	r.Message = "Оплата готівкою — прийнята"
	r.TotalDebited = amount
	system.AddReceipt(r)
	return r, nil
}

// ---------- Банківський переказ ----------

type BankTransfer struct{}

func (b BankTransfer) Name() string { return "Банківський переказ" }

func (b BankTransfer) Validate(amount float64) error {
	if amount < minAmount || amount > maxAmount {
		return fmt.Errorf("сума повинна бути від %.2f до %.2f грн", minAmount, maxAmount)
	}
	return nil
}

func (b BankTransfer) CalculateFee(amount float64) float64 {
	return amount * 0.02
}

func (b BankTransfer) Process(amount float64, system *PaymentSystem) (Receipt, error) {
	acct := system.Account
	var r Receipt
	r.Method = b.Name()
	r.Amount = amount
	r.Time = time.Now()

	err := b.Validate(amount)
	if err != nil {
		r.Message = err.Error()
		system.AddReceipt(r)
		return r, err
	}

	fee := b.CalculateFee(amount)
	total := amount + fee
	if acct.Balance < total {
		err = errors.New("недостатньо коштів на рахунку")
		r.Message = err.Error()
		system.AddReceipt(r)
		return r, err
	}

	acct.Balance -= total
	r.Fee = fee
	r.TotalDebited = total
	r.Success = true
	r.Message = "Банківський переказ успішно виконано!"
	system.AddReceipt(r)
	return r, nil
}

// ---------- Платіжна система ----------

type PaymentSystem struct {
	Account *Account
	Checks  []Receipt
}

func NewPaymentSystem(owner string, balance float64) *PaymentSystem {
	return &PaymentSystem{
		Account: NewAccount(owner, balance),
		Checks:  []Receipt{},
	}
}

func (ps *PaymentSystem) AddReceipt(r Receipt) {
	ps.Checks = append(ps.Checks, r)
}

func (ps *PaymentSystem) ShowAllReceipts() {
	fmt.Println("\n--- Список усіх чеків ---")
	if len(ps.Checks) == 0 {
		fmt.Println("Чеки відсутні.")
		return
	}
	for i, r := range ps.Checks {
		status := "❌"
		if r.Success {
			status = "✅"
		}
		fmt.Printf("%d) [%s] %.2f грн | Комісія: %.2f | %s | %v | %s\n",
			i+1, r.Method, r.Amount, r.Fee, r.Message,
			r.Time.Format("15:04:05"), status)
	}
}

func (ps *PaymentSystem) ShowStats() {
	total := len(ps.Checks)
	success := 0
	sum := 0.0
	fees := 0.0
	for _, r := range ps.Checks {
		if r.Success {
			success++
			sum += r.Amount
			fees += r.Fee
		}
	}
	fmt.Println("\n--- Статистика ---")
	fmt.Printf("Всього платежів: %d\n", total)
	fmt.Printf("Успішних: %d | Неуспішних: %d\n", success, total-success)
	fmt.Printf("Загальна сума: %.2f грн\n", sum)
	fmt.Printf("Комісії: %.2f грн\n", fees)
	fmt.Printf("Поточний баланс: %.2f грн\n", ps.Account.Balance)
}

// ---------- Ввід ----------

func getMenuChoice(max int) int {
	var choice int
	fmt.Print("\nВаш вибір: ")
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > max {
		fmt.Println("Невірний вибір, спробуйте ще раз.")
		return -1
	}
	return choice
}

func getRequiredFloat(prompt string) (float64, error) {
	var f float64
	fmt.Print(prompt)
	_, err := fmt.Scanln(&f)
	if err != nil {
		return 0, fmt.Errorf("необхідно ввести число")
	}
	if f <= 0 {
		return 0, fmt.Errorf("число повинно бути більше 0")
	}
	return f, nil
}

func printPaymentHeader() {
	header := `
 ######     ##     ##  ##   ##   ##  #######  ##   ##  ######
  ##  ##   ####    ##  ##   ### ###   ##   #  ###  ##  # ## #
  ##  ##  ##  ##   ##  ##   #######   ## #    #### ##    ##
  #####   ##  ##    ####    #######   ####    ## ####    ##
  ##      ######     ##     ## # ##   ## #    ##  ###    ##
  ##      ##  ##     ##     ##   ##   ##   #  ##   ##    ##
 ####     ##  ##    ####    ##   ##  #######  ##   ##   ####
`
	fmt.Println(header)
}

func main() {
	system := NewPaymentSystem("Руслан", 2000.00)
	card := CreditCard{}
	paypal := PayPal{}
	cash := Cash{}
	bank := BankTransfer{}

	printPaymentHeader()

	for {
		fmt.Println("\n=== Система платежів ===")
		fmt.Println("\nДоступні методи оплати:")
		fmt.Println("1. Кредитна картка")
		fmt.Println("2. PayPal")
		fmt.Println("3. Готівка")
		fmt.Println("4. Банківський переказ")
		fmt.Println("5. Поповнити рахунок")
		fmt.Println("6. Переглянути баланс")
		fmt.Println("7. Показати всі чеки")
		fmt.Println("8. Показати статистику")
		fmt.Println("9. Вихід")

		choice := getMenuChoice(9)
		if choice == -1 {
			continue
		}

		switch choice {
		case 1, 2, 3, 4:
			amount, err := getRequiredFloat("\nВведіть суму платежу: ")
			if err != nil {
				fmt.Println("Помилка:", err)
				continue
			}

			fmt.Printf("\nОбробляємо платіж на суму %.2f грн...\n\n", amount)

			var processor PaymentProcessor
			switch choice {
			case 1:
				processor = card
			case 2:
				processor = paypal
			case 3:
				processor = cash
			case 4:
				processor = bank
			}

			r, _ := processor.Process(amount, system)
			fmt.Printf("💳 %s\n", r.Method)
			if r.Success {
				fmt.Println("✅", r.Message)
				fmt.Printf("💰 Сума: %.2f грн\n", r.Amount)
				fmt.Printf("💸 Комісія: %.2f грн\n", r.Fee)
				fmt.Printf("📊 До списання: %.2f грн\n", r.TotalDebited)
				fmt.Println("\nДякуємо за покупку!")
			} else {
				fmt.Println("❌ Помилка:", r.Message)
			}
		case 5:
			funds, err := getRequiredFloat("\nВведіть суму для поповнення: ")
			if err != nil {
				fmt.Println("Помилка:", err)
				continue
			}
			system.Account.Balance += funds
			fmt.Printf("\nРахунок поповнено на %.2f грн!\n", funds)
		case 6:
			fmt.Printf("\nРахунок %s: %.2f грн\n", system.Account.Owner, system.Account.Balance)
		case 7:
			system.ShowAllReceipts()
		case 8:
			system.ShowStats()
		case 9:
			fmt.Println("\nДякуємо, що скористались системою!")
			return
		}
	}
}
