package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Кольори для консолі (ANSI escape-коди)
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

func main() {
	fmt.Println(Cyan + "===Валідатор даних===" + Reset)
	fmt.Println("Виберіть опцію:")
	fmt.Println("1. Перевірка email-адреси")
	fmt.Println("2. Перевірка надійності пароля")
	fmt.Println("3. Перевірка телефонного номера")
	fmt.Println("4. Перевірка IP-адреси")
	fmt.Println("5. Перевірка URL-адреси")
	fmt.Println("0. Вихід")

	var choice int
	fmt.Print("\nВаш вибір: ")
	fmt.Scanln(&choice)

	switch choice {

	//Перевірка email-адреси
	case 1:
		var email string
		fmt.Print("Введіть email-адресу: ")
		fmt.Scanln(&email)
		result := true
		var errors []string

		if strings.Count(email, "@") == 0 {
			errors = append(errors, "Немає символа '@'")
		} else if strings.Count(email, "@") > 1 {
			errors = append(errors, "Більше ніж один символ '@'")
		}
		parts := strings.SplitN(email, "@", 2)
		local := parts[0]
		domain := parts[1]

		if len(local) < 1 {
			errors = append(errors, "Порожня локальна частина (до @)")
		}
		if len(domain) < 1 {
			errors = append(errors, "Порожня доменна частина (після @)")
		}
		if strings.Contains(email, " ") {
			errors = append(errors, "У адресі є пробіли")
		}
		if !regexp.MustCompile(`^[A-Za-z0-9._-]+$`).MatchString(email) {
			errors = append(errors, "Недозволені символи") // (у локальній частині)
		}
		if !strings.Contains(domain, ".") {
			errors = append(errors, "У доменній частині немає крапки '.'")
		} else {
			parts := strings.Split(domain, ".")
			tld := parts[len(parts)-1]
			if len(tld) < 2 || len(tld) > 6 {
				errors = append(errors, "Після останньої крапки має бути 2–6 символів")
			}
		}

		// Вивід результату з кольорами
		if len(errors) > 0 {
			result = false
		}

		if !result {
			fmt.Println(Red + "Результат: Невалідно! Причини:" + Reset)
			fmt.Println(Red + strings.Join(errors, "\n") + Reset)
		} else {
			fmt.Println(Green + "Результат: Валідно!" + Reset)
		}

	//Перевірка пароля
	case 2:
		var password string
		fmt.Print("Введіть пароль: ")
		fmt.Scanln(&password)
		result := true
		var errors []string

		if len(password) < 8 {
			errors = append(errors, "Пароль занадто короткий (мінімум 8 символів)")
		}
		if !regexp.MustCompile(`[a-z]`).MatchString(password) {
			errors = append(errors, "Немає малої літери")
		}
		if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
			errors = append(errors, "Немає великої літери")
		}
		if !regexp.MustCompile(`[0-9]`).MatchString(password) {
			errors = append(errors, "Немає цифри")
		}
		if !regexp.MustCompile(`[!@#\$%\^&\*\(\)\-_=+\[\]\{\}\|;:'",.<>/?]`).MatchString(password) {
			errors = append(errors, "Немає спецсимволу (!@#$%^&*)")
		}
		if strings.Contains(password, " ") {
			errors = append(errors, "Пароль містить пробіли")
		}

		// Вивід результату з кольорами
		if len(errors) > 0 {
			result = false
		}

		if !result {
			fmt.Println(Red + "Результат: Невалідно! Причини:" + Reset)
			fmt.Println(Red + strings.Join(errors, "\n") + Reset)
		} else {
			fmt.Println(Green + "Результат: Валідно!" + Reset)
		}

	//Перевірка телефонного номера
	case 3:
		var phone string
		fmt.Print("Введіть номер телефону: ")
		fmt.Scanln(&phone)
		result := true
		var errors []string
		digits := 0

		if !strings.HasPrefix(phone, "+") {
			errors = append(errors, "Номер має починатися з '+' (міжнародний формат)")
		}

		for _, ch := range phone {
			if unicode.IsDigit(ch) {
				digits++
			} else if strings.ContainsRune("+-() ", ch) {
				continue
			} else {
				errors = append(errors, fmt.Sprintf("Недозволений символ: %q", ch))
			}
		}

		if digits < 10 || digits > 15 {
			errors = append(errors, "Кількість цифр має бути від 10 до 15")
		}

		clean := ""
		for _, ch := range phone {
			if unicode.IsDigit(ch) {
				clean += string(ch)
			}
		}

		if strings.HasPrefix(phone, "+380") {
			if len(clean) != 12 {
				errors = append(errors, "Український номер повинен мати 12 цифр")
			}
			operator := clean[3:6]
			validOperators := []string{"050", "063", "066", "067", "068", "091", "092", "093", "094", "095", "096", "097", "098", "099"}
			found := false
			for _, op := range validOperators {
				if operator == op {
					found = true
					break
				}
			}
			if !found {
				errors = append(errors, "Невідомий код оператора: "+operator)
			}
		}

		// Вивід результату з кольорами
		if len(errors) > 0 {
			result = false
		}

		if !result {
			fmt.Println(Red + "Результат: Невалідно! Причини:" + Reset)
			fmt.Println(Red + strings.Join(errors, "\n") + Reset)
		} else {
			fmt.Println(Green + "Результат: Валідно!" + Reset)
		}

	//Перевірка IP-адреси
	case 4:
		var ip string
		fmt.Print("Введіть IP-адресу: ")
		fmt.Scanln(&ip)
		result := true
		var errors []string

		parts := strings.Split(ip, ".")
		if len(parts) != 4 {
			errors = append(errors, "Не відповідає формату X.X.X.X")
		}
		if strings.Contains(ip, " ") {
			errors = append(errors, "IP-адреса містить пробіли")
		}

		for _, p := range parts {
			num, err := strconv.Atoi(p)
			if err != nil {
				errors = append(errors, "Частина не є числом: "+p)
				continue
			}
			if num < 0 || num > 255 {
				errors = append(errors, fmt.Sprintf("Частина виходить за межі 0–255: %d", num))
			}
		}

		// Вивід результату з кольорами
		if len(errors) > 0 {
			result = false
		}

		if !result {
			fmt.Println(Red + "Результат: Невалідно! Причини:" + Reset)
			fmt.Println(Red + strings.Join(errors, "\n") + Reset)
		} else {
			fmt.Println(Green + "Результат: Валідно!" + Reset)
		}

	//Перевірка URL-адреси
	case 5:
		var url string
		fmt.Print("Введіть URL: ")
		fmt.Scanln(&url)
		result := true
		var errors []string

		if strings.HasPrefix(url, "http://") {
			url = strings.TrimPrefix(url, "http://")
		} else if strings.HasPrefix(url, "https://") {
			url = strings.TrimPrefix(url, "https://")
		} else {
			errors = append(errors, "Відсутній протокол (http:// або https://)")
		}

		if idx := strings.IndexAny(url, "/?#"); idx != -1 {
			url = url[:idx]
		}

		if !strings.Contains(url, ".") {
			errors = append(errors, "У доменній частині немає крапки '.'")
		} else {
			parts := strings.Split(url, ".")
			tld := parts[len(parts)-1]
			if len(tld) < 2 || len(tld) > 6 {
				errors = append(errors, "Після останньої крапки має бути 2–6 символів")
			}
		}

		if !regexp.MustCompile(`^[A-Za-z0-9._-]+$`).MatchString(url) {
			errors = append(errors, "Домен містить недозволені символи")
		}

		if strings.Contains(url, " ") {
			errors = append(errors, "Домен містить пробіли")
		}

		// Вивід результату з кольорами
		if len(errors) > 0 {
			result = false
		}

		if !result {
			fmt.Println(Red + "Результат: Невалідно! Причини:" + Reset)
			fmt.Println(Red + strings.Join(errors, "\n") + Reset)
		} else {
			fmt.Println(Green + "Результат: Валідно!" + Reset)
		}

	case 0:
		fmt.Println(Yellow + "Вихід із програми." + Reset)
		return

	default:
		fmt.Println(Red + "Невірний вибір опції!" + Reset)
	}
}
