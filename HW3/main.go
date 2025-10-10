package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

/*
   === АНАЛІЗАТОР ТЕКСТУ ===
   Програма дозволяє аналізувати текст, введений користувачем:
   - Підраховує кількість слів у тексті
   - Знаходить кількість входжень заданого слова (незалежно від регістру)
   - Визначає найдовше слово та його довжину (у символах)
   - Знаходить перше слово, що починається на задану літеру
   - Підтримує повторний аналіз нового тексту
*/

func main() {
	for {
		fmt.Println("=== АНАЛІЗАТОР ТЕКСТУ ===")
		fmt.Println("Введіть текст для аналізу:")

		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Помилка читання:", err)
			return
		}

		text = strings.TrimSpace(text) // Видаляємо зайві пробіли та \n

		// Перевірка на порожній текст
		if len(text) == 0 {
			fmt.Println("Текст порожній! Спробуйте ще раз.")
			continue
		}

		text = strings.ToLower(text) // Переводимо текст у нижній регістр

		// Замінюємо всі розділові знаки на пробіли
		punctuations := ",.!?;:-\"'()[]{}/\\|+="
		cleanedText := ""
		for _, r := range text {
			if strings.ContainsRune(punctuations, r) {
				cleanedText += " "
			} else {
				cleanedText += string(r)
			}
		}
		text = cleanedText

		// Запитуємо слово для пошуку
		fmt.Print("Введіть слово або літеру для пошуку: ")
		wordR, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Помилка читання:", err)
			return
		}
		word := strings.TrimSpace(strings.ToLower(wordR)) // Видаляємо пробіли та \n

		// Розбиваємо текст на слова
		words := strings.Fields(text)

		// Перевірка на порожній список слів
		if len(words) == 0 {
			fmt.Println("У тексті немає слів для аналізу.")
			continue
		}

		var largestWord = ""
		var count int

		// Проходимо по кожному слову:
		for _, w := range words {
			// Рахуємо кількість входжень шуканого слова/літери
			if w == word {
				count++
			}
			// Шукаємо найдовше слово (за кількістю символів)
			if utf8.RuneCountInString(w) > utf8.RuneCountInString(largestWord) {
				largestWord = w
			}
		}

		// Вивід результатів пошуку слова
		if word == "" {
			fmt.Println("Слово для пошуку не введено.")
		} else if count == 0 {
			fmt.Printf("Слово або літера \"%s\" не знайдено у тексті.\n", strings.TrimSpace(wordR))
		} else {
			fmt.Printf("Слово або літера \"%s\" зустрічається %d раз(и).\n", strings.TrimSpace(wordR), count)
		}

		// Вивід загальної кількості слів та найдовшого слова
		fmt.Printf("У тексті всього %d слів.\n", len(words))
		fmt.Printf("Найдовше слово тексту: %s (%d символів)\n", largestWord, utf8.RuneCountInString(largestWord))

		// Запитуємо літеру для пошуку першого слова, що починається з неї
		var letter string
		fmt.Println("Введіть літеру для пошуку першого слова:")
		fmt.Scanln(&letter)
		letter = strings.ToLower(strings.TrimSpace(letter))

		var firstLetterWord string
		for _, w := range words {
			// Перевіряємо, чи слово починається на введену літеру
			if len(w) > 0 && letter == string([]rune(w)[0]) {
				firstLetterWord = w
				break // Знайшли перше слово — виходимо з циклу
			}
		}
		if letter == "" {
			fmt.Println("Літера для пошуку не введена.")
		} else if firstLetterWord == "" {
			fmt.Printf("Слів, що починаються на \"%s\", не знайдено.\n", letter)
		} else {
			fmt.Printf("Перше слово, що починається на \"%s\": %s\n", letter, firstLetterWord)
		}

		// Запит на повторний аналіз
		var choice string
		fmt.Println("Хочете проаналізувати інший текст? (yes/no)")
		fmt.Scanln(&choice)

		switch strings.ToLower(choice) {
		case "yes":
			continue
		case "no":
			fmt.Println("Вихід...")
			return
		default:
			return
		}
	}
}
