package main

import (
	"fmt"
	"math"
)

// Константи
const dollar_rate = 38.5
const euro_rate = 42.1

func main() {

	// Оголошення всіх змінних

	var initial_amount float64 = 5000.0                                   // Початкова сума
	var monthly_savings float64 = 1500.0                                  // Щомісячні накопичення
	var annual_interest_rate float64 = 15.0                               // Річна відсоткова ставка банку
	var years int = 2                                                     // Термін накопичень в роках
	var months int = years * 12                                           // Кількість місяців
	var monthly_interest_rate float64 = annual_interest_rate / (100 * 12) // Місячна_ставка

	var totalContributions float64 = initial_amount + (monthly_savings * float64(months))
	// Загальна сума внесків: початкова сума + всі щомісячні внески

	var compoundInitial float64 = initial_amount * math.Pow((1+monthly_interest_rate), float64(months))
	// Нарощення початкової суми зі складними відсотками

	var contributionEffect float64 = (monthly_savings * (math.Pow((1+monthly_interest_rate), float64(months)) - 1) / monthly_interest_rate)
	// Майбутня вартість ануїтету (щомісячних внесків) з урахуванням складних відсотків

	var finalAmount float64 = compoundInitial + contributionEffect
	// Фінальна сума: нарощена початкова сума + нарощені щомісячні внески

	var interest float64 = finalAmount - totalContributions
	// Нараховані відсотки: різниця між фінальною сумою та сумою внесків

	// Виведення у консоль
	fmt.Println("=== КАЛЬКУЛЯТОР НАКОПИЧЕНЬ ===")
	fmt.Println("\nПочаткові дані:")
	fmt.Printf("- Початкова сума: %.1f грн\n", initial_amount)
	fmt.Printf("- Щомісячні накопичення: %.2f грн\n", monthly_savings)
	fmt.Printf("- Річна ставка: %.1f%%\n", annual_interest_rate)
	fmt.Printf("- Термін: %d роки (%d місяців)\n", years, months)

	fmt.Println("\nРезультати:")
	fmt.Printf("- Загальна сума внесків: %.2f грн\n", totalContributions)
	fmt.Printf("- Нараховані відсотки: %.2f грн\n", interest)
	fmt.Printf("- Фінальна сума: %.2f грн\n", finalAmount)

	fmt.Println("\nУ доларах США:")
	fmt.Printf("- Загальна сума внесків: $%.2f\n", totalContributions/dollar_rate)
	fmt.Printf("- Фінальна сума: $%.2f\n", finalAmount/dollar_rate)

	fmt.Println("\nУ евро:")
	fmt.Printf("- Загальна сума внесків: €%.2f\n", totalContributions/euro_rate)
	fmt.Printf("- Фінальна сума: €%.2f\n", finalAmount/euro_rate)
}
