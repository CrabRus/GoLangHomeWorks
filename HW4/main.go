package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// === КОНСТАНТИ === //
const (
	// Базові тарифи (в умові тариф впливає лише на експрес: додається половина суми)
	standardRate float64 = 0   // стандартна доставка без доплати
	expressRate  float64 = 0.5 // експрес — +50% до базової вартості

	// Вартість пакувальних матеріалів за м²
	standardPackagingRate   float64 = 40.0
	reinforcedPackagingRate float64 = 50.0
	premiumPackagingRate    float64 = 75.0
)

func main() {

	for {
		// === ГОЛОВНЕ МЕНЮ === //
		fmt.Println("=== Калькулятор доставки посилок ===")
		fmt.Println("1. Розрахунок вартості доставки")
		fmt.Println("2. Оцінка часу доставки")
		fmt.Println("3. Розрахунок пакувальних матеріалів")
		fmt.Println("4. Вихід")

		choice := getIntInput("\nВаш вибір: ")

		switch choice {
		// === 1. РОЗРАХУНОК ВАРТОСТІ === //
		case 1:
			fmt.Println("--- Розрахунок вартості доставки ---")

			// Отримання ваги
			weight, err := getNumberInput("\nВведіть вагу посилки (кг): ")
			if err != nil {
				fmt.Println("Помилка:", err)
				continue
			}

			// Тип доставки
			fmt.Println("Виберіть тип доставки:")
			fmt.Println("1. Стандартна")
			fmt.Println("2. Експрес")
			deliveryType := getIntInput("")
			switch deliveryType {
			case 1, 2:
			default:
				fmt.Println("Такої опції немає")
				continue
			}

			// Відстань
			distance, err := getNumberInput("\nВведіть відстань доставки (км): ")
			if err != nil {
				fmt.Println("Помилка:", err)
				continue
			}

			// Статус клієнта
			fmt.Println("Виберіть статус клієнта:")
			fmt.Println("1. Звичайний")
			fmt.Println("2. Постійний")
			clientStatus := getIntInput("")
			switch clientStatus {
			case 1, 2:
			default:
				fmt.Println("Такої опції немає")
				continue
			}

			// === РОЗРАХУНОК === //
			basePrice := calculateBasePrice(weight, distance)
			additionalPrice := calculateDeliveryTypePrice(basePrice, deliveryType)
			discount := calculateDiscount(basePrice, additionalPrice, clientStatus)
			finalPrice := calculateFinalPrice(basePrice, additionalPrice, discount)

			// === ВИВІД === //
			fmt.Println("Результати розрахунку вартості:")
			fmt.Printf("\nБазова вартість: %.02f", basePrice)
			switch deliveryType {
			case 1:
				fmt.Printf("\nДодаткова вартість (Стандарт): %.02f", additionalPrice)
			case 2:
				fmt.Printf("\nДодаткова вартість (Експрес): %.02f", additionalPrice)
			}
			fmt.Printf("\nЗнижка (Постійний клієнт): %.02f", discount)
			fmt.Printf("\nЗагальна вартість: %.02f", finalPrice)

			// Повернення до меню
			if getYesNoInput("\nБажаєте повернутися до головного меню? (так/ні): ") {
				fmt.Println("Повертаємось до головного меню...")
				continue
			} else {
				fmt.Println("Програма завершена.")
				return
			}

		// === 2. ОЦІНКА ЧАСУ ДОСТАВКИ === //
		case 2:
			fmt.Println("--- Оцінка часу доставки ---")

			distance, err := getNumberInput("\nВведіть відстань доставки (км): ")
			if err != nil {
				fmt.Println("Помилка:", err)
				continue
			}

			fmt.Println("Виберіть тип доставки:")
			fmt.Println("1. Стандартна")
			fmt.Println("2. Експрес")
			deliveryType := getIntInput("")
			switch deliveryType {
			case 1, 2:
			default:
				fmt.Println("Такої опції немає")
				continue
			}

			fmt.Println("Виберіть погодні умови:")
			fmt.Println("1. Хороші")
			fmt.Println("2. Задовільні")
			fmt.Println("3. Погані")
			weatherCondition := getIntInput("")
			switch weatherCondition {
			case 1, 2, 3:
			default:
				fmt.Println("Такої опції немає")
				continue
			}

			// Вихідний день
			isWeekend := getYesNoInput("\nСьогодні вихідний день? (так/ні): ")

			// Розрахунок часу
			baseDeliveryTime := calculateBaseDeliveryTime(distance, deliveryType)
			weatherDelay := addWeatherDelay(weatherCondition)
			finalDeliveryTime := calculateFinalDeliveryTime(baseDeliveryTime, weatherDelay, isWeekend)
			deliveryDate := calculateDeliveryDate(finalDeliveryTime)

			// Вивід
			fmt.Println("Результати оцінки часу доставки:")
			switch deliveryType {
			case 1:
				fmt.Printf("Базовий час доставки (Стандартна): %.02f днів", baseDeliveryTime)
			case 2:
				fmt.Printf("Базовий час доставки (Експрес): %.02f днів", baseDeliveryTime)
			}
			fmt.Printf("\nЗатримка через погодні умови: %.02f днів", weatherDelay)
			fmt.Printf("\nЗагальний орієнтовний час доставки: %.02f днів", finalDeliveryTime)
			fmt.Printf("\nОрієнтовна дата прибуття: %s", deliveryDate)

			if getYesNoInput("\nБажаєте повернутися до головного меню? (так/ні): ") {
				fmt.Println("Повертаємось до головного меню...")
				continue
			} else {
				fmt.Println("Програма завершена.")
				return
			}

		// === 3. ПАКУВАЛЬНІ МАТЕРІАЛИ === //
		case 3:
			fmt.Println("--- Розрахунок пакувальних матеріалів ---")

			length, err := getNumberInput("\nВведіть довжину посилки (см): ")
			if err != nil {
				fmt.Println("Помилка:", err)
				continue
			}
			width, err := getNumberInput("\nВведіть ширину посилки (см): ")
			if err != nil {
				fmt.Println("Помилка:", err)
				continue
			}
			height, err := getNumberInput("\nВведіть висоту посилки (см): ")
			if err != nil {
				fmt.Println("Помилка:", err)
				continue
			}

			fmt.Println("Виберіть тип пакувального матеріалу:")
			fmt.Println("1. Стандартний картон")
			fmt.Println("2. Посилений картон з плівкою")
			fmt.Println("3. Преміум пакування")
			materialType := getIntInput("")
			switch materialType {
			case 1, 2, 3:
			default:
				fmt.Println("Такої опції немає")
				continue
			}

			materialAmount := calculatePackagingMaterial(length, width, height)
			packingCost := calculatePackagingCost(materialAmount, materialType)

			fmt.Println("Результати розрахунку пакувальних матеріалів:")
			fmt.Printf("\nНеобхідна кількість пакувального матеріалу: %.02f м²", materialAmount)
			fmt.Printf("\nВартість пакувальних матеріалів: %.02f грн", packingCost)

			if getYesNoInput("\nБажаєте повернутися до головного меню? (так/ні): ") {
				fmt.Println("Повертаємось до головного меню...")
				continue
			} else {
				fmt.Println("Програма завершена.")
				return
			}

		case 4:
			return
		default:
			fmt.Println("Помилка вводу")
		}
	}
}

// === ФУНКЦІЇ РОЗРАХУНКУ === //

// Базова вартість доставки (залежить від ваги і відстані)
func calculateBasePrice(weight, distance float64) float64 {
	return weight * distance * 0.6
}

// Додаткова вартість (експрес або стандарт)
func calculateDeliveryTypePrice(basePrice float64, deliveryType int) float64 {
	switch deliveryType {
	case 1:
		return basePrice * standardRate
	case 2:
		return basePrice * expressRate
	}
	return 0
}

// Знижка для постійних клієнтів
func calculateDiscount(basePrice, additionalPrice float64, clientStatus int) float64 {
	switch clientStatus {
	case 1:
		return 0
	case 2:
		return (basePrice + additionalPrice) * 0.1
	}
	return 0
}

// Фінальна ціна доставки
func calculateFinalPrice(basePrice, additionalPrice, discount float64) float64 {
	return basePrice + additionalPrice - discount
}

// Базовий час доставки (залежить від типу)
func calculateBaseDeliveryTime(distance float64, deliveryType int) float64 {
	baseDeliveryTime := distance / 100
	switch deliveryType {
	case 1:
		return baseDeliveryTime * 0.8
	case 2:
		return baseDeliveryTime * 0.4
	default:
		return baseDeliveryTime
	}
}

// Затримка через погоду
func addWeatherDelay(weatherCondition int) float64 {
	switch weatherCondition {
	case 1:
		return 0
	case 2:
		return 0.5
	case 3:
		return 1
	default:
		return 0
	}
}

// Фінальний час доставки (враховує вихідні)
func calculateFinalDeliveryTime(baseTime, weatherDelay float64, isWeekend bool) float64 {
	if isWeekend {
		return baseTime + weatherDelay + 1
	}
	return baseTime + weatherDelay
}

// Визначає дату прибуття
func calculateDeliveryDate(deliveryDays float64) string {
	days := int(math.Round(deliveryDays)) // класичне округлення
	if days < 1 {
		days = 1
	}
	currentDate := time.Now()
	deliveryDate := currentDate.AddDate(0, 0, days)
	return deliveryDate.Format("02.01.2006")
}

// Площа пакувального матеріалу (у м²)
func calculatePackagingMaterial(length, width, height float64) float64 {
	value := 2 * (length*width + length*height + width*height) / 10000
	roundedValue := math.Round(value*100) / 100
	return roundedValue
}

// Вартість пакувальних матеріалів
func calculatePackagingCost(materialAmount float64, materialType int) float64 {
	switch materialType {
	case 1:
		return materialAmount * standardPackagingRate
	case 2:
		return materialAmount * reinforcedPackagingRate
	case 3:
		return materialAmount * premiumPackagingRate
	default:
		return materialAmount
	}
}

// === ОТРИМАННЯ ДАНИХ ВІД КОРИСТУВАЧА === //

// Ввід числа з перевіркою (float)
func getNumberInput(prompt string) (float64, error) {
	var input float64
	fmt.Print(prompt)
	_, err := fmt.Scanln(&input)
	if err != nil {
		return 0, fmt.Errorf("invalid input")
	}
	if input <= 0 {
		return 0, fmt.Errorf("Число менше за 0")
	}
	return input, nil
}

// Ввід цілого числа (int)
func getIntInput(prompt string) int {
	var input int
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}

// Ввід відповіді "так/ні"
func getYesNoInput(prompt string) bool {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "так" || input == "t" || input == "y" || input == "yes"
}
