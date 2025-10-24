package main

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

var (
	players       []string
	scores        map[string]int
	matchesPlayed map[string]int
	wins          map[string]int
	losses        map[string]int
)

func main() {
	initializeSystem()
	displayMenu()
}

// ======Обов'язкові функції для роботи з гравцями: ======

// Реєструє нового гравця з початковим рейтингом 1000
func registerPlayer(nickname string) error {
	if !playerExists(nickname) {
		players = append(players, nickname)
		scores[nickname] = 1000
		matchesPlayed[nickname] = 0
		wins[nickname] = 0
		losses[nickname] = 0
		fmt.Printf("Гравець \"%s\" успішно зареєстрований з рейтингом 1000!", nickname)
		return nil
	} else {
		return fmt.Errorf("гравець %s вже існує", nickname)
	}

}

// Видаляє гравця з системи
func removePlayer(nickname string) error {
	if playerExists(nickname) {
		if index := findPlayerIndex(nickname); index != -1 {
			players = slices.Delete(players, index, index+1)
		}
		delete(scores, nickname)
		delete(matchesPlayed, nickname)
		delete(wins, nickname)
		delete(losses, nickname)
		fmt.Printf("Гравець \"%s\" успішно видалений!", nickname)
		return nil

	}
	return fmt.Errorf("гравця %s не існує", nickname)
}

// Знаходить індекс гравця в слайсі
func findPlayerIndex(nickname string) int {
	for i, n := range players {
		if n == nickname {
			return i
		}
	}
	return -1
}

// Перевіряє чи існує гравець
func playerExists(nickname string) bool {
	return slices.Contains(players, nickname)
}

// Відображає всіх гравців з їх статистикою
func displayAllPlayers() {
	fmt.Println("--- Всі гравці ---")
	for i, nickname := range players {
		fmt.Printf("%d. ", i+1)
		displayPlayerStats(nickname)
	}
}

// ====== Функції для роботи з рейтингом: ======

// Оновлює рейтинг після матчу (true = перемога, false = поразка)
func updateRating(nickname string, won bool, pointsChange int) {
	if won {
		scores[nickname] += pointsChange
		wins[nickname]++
	} else {
		scores[nickname] -= pointsChange
		losses[nickname]++
	}
	updateMatchesPlayed(nickname)
}

func updateMatchesPlayed(nickname string) {
	matchesPlayed[nickname] = wins[nickname] + losses[nickname]
}

// Повертає топ-10 гравців за рейтингом
func getTopPlayers(count int) []string {
	sorted, err := sortPlayersByRating()
	if err != nil {
		fmt.Println("Помилка:", err)
		return nil
	}

	// Якщо кількість гравців менша за count — обрізаємо до доступної довжини
	if count > len(sorted) {
		count = len(sorted)
	}

	return sorted[:count]
}

// Знаходить гравців у діапазоні рейтингу
func findPlayersByRatingRange(minRating, maxRating int) []string {
	var foundPlayers []string
	for k, v := range scores {
		if v >= minRating && v <= maxRating {
			foundPlayers = append(foundPlayers, k)
		}
	}
	return foundPlayers
}

// Розраховує середній рейтинг всіх гравців
func calculateAverageRating() float64 {
	if len(scores) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range scores {
		sum += float64(v)
	}
	return sum / float64(len(scores))
}

// ====== Функції для статистики: ======

// Знаходить гравця з найвищим рейтингом
func getBestPlayer() string {
	sorted, err := sortPlayersByRating()
	if err != nil {
		fmt.Println("Помилка:", err)
		return ""
	} else {
		var firstNickname string
		for _, nickname := range sorted {
			firstNickname = nickname
			break
		}
		return firstNickname
	}
}

// Знаходить гравця з найнижчим рейтингом
func getWorstPlayer() string {
	sorted, err := sortPlayersByRating()
	if err != nil {
		fmt.Println("Помилка:", err)
		return ""
	} else {
		var lastNickname string
		for _, nickname := range sorted {
			lastNickname = nickname
		}
		return lastNickname
	}
}

// Розраховує відсоток перемог гравця
func calculateWinRate(nickname string) (float64, error) {
	if !playerExists(nickname) {
		return 0, fmt.Errorf("гравця %s не існує", nickname)
	}

	total := matchesPlayed[nickname]
	if total == 0 {
		return 0, nil // ще не грав — відсоток перемог 0
	}

	winRate := (float64(wins[nickname]) / float64(total)) * 100
	winRate = math.Round(winRate*100) / 100 // округлення до сотих

	return winRate, nil
}

// Відображає детальну статистику гравця
func displayPlayerStats(nickname string) {
	fmt.Printf("%s - Рейтинг: %d | Матчів: %d | Перемог: %d | Поразок: %d\n",
		nickname, scores[nickname], matchesPlayed[nickname], wins[nickname], losses[nickname])
}

// Показує загальну статистику системи
func displaySystemStats() {
	fmt.Println("--- Загальна статистика ---")
	if len(players) == 0 {
		fmt.Println("Немає зареєстрованих гравців для статистики.")
		return
	}
	fmt.Printf("Кількість гравців: %d\n", len(players))
	fmt.Printf("Середній рейтинг: %.2f\n", calculateAverageRating())

	best := getBestPlayer()
	worst := getWorstPlayer()

	if best != "" {
		fmt.Printf("Гравець з найвищим рейтингом: %s (%d)\n", best, scores[best])
	}

	if worst != "" {
		fmt.Printf("Гравець з найнижчим рейтингом: %s (%d)\n", worst, scores[worst])
	}
	topPlayers := getTopPlayers(10)
	fmt.Println("Топ-10 гравців:")
	for i, nickname := range topPlayers {
		fmt.Printf("%d. ", i+1)
		displayPlayerStats(nickname)
	}

}

// ====== Допоміжні функції: ======

// Відображає меню
func displayMenu() {
	for {
		fmt.Println("=== Система рейтингу гравців ===")
		fmt.Println("")
		fmt.Println("Виберіть опцію:")
		fmt.Println("1. Зареєструвати гравця")
		fmt.Println("2. Видалити гравця")
		fmt.Println("3. Оновити рейтинг після матчу")
		fmt.Println("4. Знайти гравця")
		fmt.Println("5. Показати всіх гравців")
		fmt.Println("6. Топ-10 гравців")
		fmt.Println("7. Пошук за діапазоном рейтингу")
		fmt.Println("8. Статистика гравця")
		fmt.Println("9. Загальна статистика")
		fmt.Println("10. Вихід")

		var choice int
		fmt.Print("\nВаш вибір: ")
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Помилка введення. Спробуйте ще раз.")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("\n--- Реєстрація гравця ---")
			nickname, err := getStringInput("Введіть нікнейм гравця: ")
			if err != nil {
				fmt.Println("Помилка:", err)
				break
			}
			err = registerPlayer(nickname)
			if err != nil {
				fmt.Println("Помилка:", err)
			}

		case 2:
			fmt.Println("\n--- Видалення гравця ---")
			nickname, err := getStringInput("Введіть нікнейм гравця: ")
			if err != nil {
				fmt.Println("Помилка:", err)
				break
			}
			err = removePlayer(nickname)
			if err != nil {
				fmt.Println("Помилка:", err)
			}

		case 3:
			fmt.Println("\n--- Оновлення рейтингу ---")
			nickname, err := getStringInput("Введіть нікнейм гравця: ")
			if err != nil {
				fmt.Println("Помилка:", err)
				break
			}
			if !playerExists(nickname) {
				fmt.Println("Гравця не знайдено!")
				break
			}
			var winInt int
			fmt.Print("Гравець переміг чи програв? (1 - перемога, 0 - поразка): ")
			_, err = fmt.Scanln(&winInt)
			if err != nil || (winInt != 0 && winInt != 1) {
				fmt.Println("Введіть 1 або 0!")
				break
			}
			won := winInt == 1
			pointsChange, err := getIntInput("Введіть зміну рейтингу: ")
			if err != nil {
				fmt.Println("Помилка:", err)
				break
			}
			updateRating(nickname, won, pointsChange)
			fmt.Printf("\nРейтинг гравця \"%s\" оновлено! Новий рейтинг: %d\n", nickname, scores[nickname])

		case 4:
			fmt.Println("\n--- Пошук гравця ---")
			nickname, err := getStringInput("Введіть нікнейм гравця: ")
			if err != nil {
				fmt.Println("Помилка:", err)
				break
			}
			if playerExists(nickname) {
				fmt.Printf("Гравець \"%s\" знайдений!\n", nickname)
				displayPlayerStats(nickname)
			} else {
				fmt.Println("Гравця не знайдено!")
			}

		case 5:
			displayAllPlayers()

		case 6:
			fmt.Println("\n--- Топ-10 гравців ---")
			count, err := getIntInput("Введіть кількість гравців для відображення: ")
			if err != nil || count <= 0 {
				fmt.Println("Некоректне число!")
				break
			}
			topPlayers := getTopPlayers(count)
			if len(topPlayers) == 0 {
				fmt.Println("Немає гравців для відображення.")
				break
			}
			fmt.Println("\nТОП ГРАВЦІ:")
			for i, nickname := range topPlayers {
				fmt.Printf("%d. %s - %d очок\n", i+1, nickname, scores[nickname])
			}

		case 7:
			fmt.Println("\n--- Пошук за діапазоном рейтингу ---")
			minRating, err1 := getIntInput("Введіть мінімальний рейтинг: ")
			maxRating, err2 := getIntInput("Введіть максимальний рейтинг: ")
			if err1 != nil || err2 != nil || minRating > maxRating {
				fmt.Println("Некоректний діапазон!")
				break
			}
			found := findPlayersByRatingRange(minRating, maxRating)
			if len(found) == 0 {
				fmt.Println("Гравців у цьому діапазоні не знайдено.")
			} else {
				fmt.Println("Гравці у діапазоні:")
				for _, nickname := range found {
					displayPlayerStats(nickname)
				}
			}

		case 8:
			fmt.Println("\n--- Статистика гравця ---")
			nickname, err := getStringInput("Введіть нікнейм гравця: ")
			if err != nil {
				fmt.Println("Помилка:", err)
				break
			}
			if !playerExists(nickname) {
				fmt.Println("Гравця не знайдено!")
				break
			}
			fmt.Printf("\nСтатистика гравця \"%s\":\n", nickname)
			fmt.Printf("Поточний рейтинг: %d\n", scores[nickname])
			fmt.Printf("Зіграно матчів: %d\n", matchesPlayed[nickname])
			fmt.Printf("Перемоги: %d\n", wins[nickname])
			fmt.Printf("Поразки: %d\n", losses[nickname])
			winRate, err := calculateWinRate(nickname)
			if err != nil {
				fmt.Println("Помилка:", err)
			} else {
				fmt.Printf("Відсоток перемог: %.1f%%\n", winRate)
			}

		case 9:
			displaySystemStats()

		case 10:
			fmt.Println("\nДо побачення!")
			return

		default:
			fmt.Println("Невірний вибір. Спробуйте ще раз.")
		}

		fmt.Println("")
	}
}

// Отримує текстове введення
func getStringInput(prompt string) (string, error) {
	var input string
	fmt.Print(prompt)
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", fmt.Errorf("невірне введення")
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("рядок не може бути порожнім")
	}
	if strings.Contains(input, " ") {
		return "", fmt.Errorf("рядок не може містити пробіли")
	}

	return input, nil
}

// Отримує числове введення
func getIntInput(prompt string) (int, error) {
	var input int
	fmt.Print(prompt)
	_, err := fmt.Scanln(&input)
	if err != nil {
		return 0, fmt.Errorf("необхідно ввести число")
	}
	return input, nil
}

// Сортує гравців за рейтингом (для топ-списку)
func sortPlayersByRating() ([]string, error) {
	// Перевірка: чи є взагалі гравці
	if len(players) == 0 {
		return nil, fmt.Errorf("немає жодного зареєстрованого гравця")
	}

	// Копія слайсу, щоб не змінювати оригінал
	sorted := slices.Clone(players)

	// Сортування за рейтингом (від більшого до меншого)
	slices.SortFunc(sorted, func(a, b string) int {
		// Якщо рейтинги рівні — сортуємо за алфавітом
		if scores[a] == scores[b] {
			return strings.Compare(a, b)
		}
		// Більший рейтинг має бути вище
		if scores[a] > scores[b] {
			return -1
		}
		return 1
	})

	return sorted, nil
}

func initializeSystem() {
	players = make([]string, 0)
	scores = make(map[string]int)
	matchesPlayed = make(map[string]int)
	wins = make(map[string]int)
	losses = make(map[string]int)
}
