package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var (
	currentProductID  int = 0
	currentCustomerID int = 0
)

func main() {
	store := initializeStore()
	displayMainMenu(store)
}

// ---------- Структури ----------

type Store struct {
	Name      string
	Products  []Product
	Customers []Customer
	Orders    []Order
	Carts     map[int]Cart
}

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Category    string
	Stock       int
	IsActive    bool
}

type Customer struct {
	ID      int
	Name    string
	Email   string
	Phone   string
	Address string
}

type CartItem struct {
	ProductID int
	Name      string
	UnitPrice float64
	Quantity  int
}

type Cart struct {
	CustomerID int
	Items      map[int]CartItem
	Discount   float64
}

type ShippingInfo struct {
	Address      string
	Method       string
	Cost         float64
	TrackingCode string
}

type OrderItem struct {
	ProductID int
	Name      string
	UnitPrice float64
	Quantity  int
	LineTotal float64
}

type Order struct {
	ID         int
	CustomerID int
	Items      []OrderItem
	Subtotal   float64
	Discount   float64
	Shipping   float64
	Total      float64
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ShippingTo *ShippingInfo
}

// ---------- Ініціалізація ----------

func initializeStore() *Store {
	return &Store{
		Name:      "TechStore",
		Products:  make([]Product, 0),
		Customers: make([]Customer, 0),
		Orders:    make([]Order, 0),
		Carts:     make(map[int]Cart),
	}
}

// ---------- Головне меню ----------

func displayMainMenu(store *Store) {
	for {
		fmt.Printf("\n=== Онлайн-магазин \"%s\" ===\n", store.Name)
		fmt.Println("1. Управління товарами")
		fmt.Println("2. Управління клієнтами")
		fmt.Println("3. Кошик покупок")
		fmt.Println("4. Замовлення")
		fmt.Println("5. Статистика магазину")
		fmt.Println("6. Вихід")

		choice := getMenuChoice(6)
		switch choice {
		case 1:
			displayProductsMenu(store)
		case 2:
			displayCustomersMenu(store)
		case 3:
			displayCartsMenu(store)
		case 4:
			displayOrdersMenu(store)
		case 5:
			displayStatistics(store)
		case 6:
			fmt.Println("До побачення!")
			return
		default:
			continue
		}
	}
}

// ---------- Меню товарів ----------

func displayProductsMenu(store *Store) {
	for {
		fmt.Println("\n=== МЕНЮ ТОВАРІВ ===")
		fmt.Println("1. Додати товар")
		fmt.Println("2. Переглянути всі товари")
		fmt.Println("3. Знайти товар за ID")
		fmt.Println("4. Пошук за категорією")
		fmt.Println("5. Оновити або видалити товар")
		fmt.Println("6. Повернутися до головного меню")

		choice := getMenuChoice(6)
		switch choice {
		case 1:
			addProduct(store)
		case 2:
			viewAllProducts(store)
		case 3:
			viewProductByID(store)
		case 4:
			viewProductsByCategory(store)
		case 5:
			editOrDeleteProduct(store)
		case 6:
			return
		default:
			continue
		}
	}
}

// ---------- Функції для товарів ----------

func NewProduct(name, description string, price float64, category string, stock int) *Product {
	currentProductID++
	return &Product{
		ID:          currentProductID,
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
		Stock:       stock,
		IsActive:    true,
	}
}

func addProduct(store *Store) {
	fmt.Println("\n--- Додавання товару ---")

	name, err := getRequiredString("Введіть назву товару: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}
	description, _ := getRequiredString("Введіть опис: ")
	price, err := getRequiredFloat("Введіть ціну: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}
	category, _ := getRequiredString("Введіть категорію: ")
	stock, err := getRequiredInt("Введіть кількість на складі: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}

	product := NewProduct(name, description, price, category, stock)
	store.Products = append(store.Products, *product)

	fmt.Printf("Товар \"%s\" успішно додано!\n", product.Name)
}

func viewAllProducts(store *Store) {
	if len(store.Products) == 0 {
		fmt.Println("Каталог порожній.")
		return
	}

	fmt.Println("\n--- Каталог товарів ---")
	for _, p := range store.Products {
		if p.IsActive {
			fmt.Printf("ID: %d | %s | %.2f € | Категорія: %s | На складі: %d\n",
				p.ID, p.Name, p.Price, p.Category, p.Stock)
		}
	}
}

func viewProductByID(store *Store) {
	id, err := getRequiredInt("Введіть ID товару: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}
	product := findProductByID(store, id)
	if product == nil || !product.IsActive {
		fmt.Println("Товар не знайдено.")
		return
	}
	fmt.Printf("ID: %d | %s | %.2f € | Категорія: %s | Опис: %s | На складі: %d\n",
		product.ID, product.Name, product.Price, product.Category, product.Description, product.Stock)
}

func viewProductsByCategory(store *Store) {
	category, err := getRequiredString("Введіть категорію: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}
	products := findProductsByCategory(store, category)
	if len(products) == 0 {
		fmt.Println("Товари не знайдені.")
		return
	}
	fmt.Printf("\n--- Товари у категорії \"%s\" ---\n", category)
	for _, p := range products {
		fmt.Printf("ID: %d | %s | %.2f € | На складі: %d\n", p.ID, p.Name, p.Price, p.Stock)
	}
}

func findProductByID(store *Store, id int) *Product {
	for i := range store.Products {
		if store.Products[i].ID == id {
			return &store.Products[i]
		}
	}
	return nil
}

func findProductsByCategory(store *Store, category string) []Product {
	var result []Product
	for _, p := range store.Products {
		if p.IsActive && strings.EqualFold(p.Category, category) {
			result = append(result, p)
		}
	}
	return result
}

func editOrDeleteProduct(store *Store) {
	id, err := getRequiredInt("Введіть ID товару: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}
	product := findProductByID(store, id)
	if product == nil || !product.IsActive {
		fmt.Println("Товар не знайдено.")
		return
	}

	action, _ := getRequiredString("Введіть 'update' для оновлення або 'delete' для видалення: ")
	switch strings.ToLower(action) {
	case "update":
		updateProductFields(product)
	case "delete":
		product.IsActive = false
		fmt.Printf("Товар \"%s\" видалено!\n", product.Name)
	default:
		fmt.Println("Невідома дія.")
	}
}

func updateProductFields(product *Product) {
	fmt.Println("Оновлення даних товару (Enter — пропустити):")

	if name, _ := getOptionalString("Нова назва: "); name != "" {
		product.Name = name
	}
	if price, _ := getOptionalFloat("Нова ціна: "); price > 0 {
		product.Price = price
	}
	if stock, _ := getOptionalInt("Нова кількість: "); stock > 0 {
		product.Stock = stock
	}
	fmt.Println("Товар оновлено.")
}

// ----- Управління клієнтами -----

func NewCustomer(name, phone, email, address string) *Customer {
	currentCustomerID++
	return &Customer{
		ID:      currentCustomerID,
		Name:    name,
		Phone:   phone,
		Email:   email,
		Address: address,
	}
}

// ==== Меню управління клієнтами ====
func displayCustomersMenu(store *Store) {
	for {
		fmt.Println("\n=== УПРАВЛІННЯ КЛІЄНТАМИ ===")
		fmt.Println("1. Зареєструвати нового клієнта")
		fmt.Println("2. Переглянути інформацію про клієнта")
		fmt.Println("3. Оновити контактні дані клієнта")
		fmt.Println("4. Повернутися до головного меню")

		choice := getMenuChoice(4)
		if choice == -1 {
			continue
		}

		switch choice {
		case 1:
			registerCustomer(store)
		case 2:
			viewCustomers(store)
		case 3:
			updateCustomer(store)
		case 4:
			fmt.Println("Повернення до головного меню...")
			return
		}
	}
}

func registerCustomer(store *Store) {
	fmt.Println("\n--- Реєстрація нового клієнта ---")

	name, err := getRequiredString("Введіть ім’я клієнта: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}

	phone, err := getRequiredString("Введіть номер телефону: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}

	valid, errs_phone := validatePhone(phone)
	if !valid {
		fmt.Println("Помилка у номері телефону:")
		for _, e := range errs_phone {
			fmt.Println("-", e)
		}
		return
	}

	email, err := getRequiredString("Введіть e-mail: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}
	valid, errs_email := validateEmail(email)
	if !valid {
		fmt.Println("Помилка у e-mail:")
		for _, e := range errs_email {
			fmt.Println("-", e)
		}
		return
	}

	address, err := getRequiredString("Введіть адресу: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}

	customer := NewCustomer(name, phone, email, address)
	store.Customers = append(store.Customers, *customer)

	fmt.Printf("Клієнта %s успішно зареєстровано!\n", customer.Name)
}

func viewCustomers(store *Store) {
	if len(store.Customers) == 0 {
		fmt.Println("\nНемає зареєстрованих клієнтів.")
		return
	}
	fmt.Println("\n--- Список клієнтів ---")
	for _, c := range store.Customers {
		fmt.Printf("[%d] %s | Телефон: %s | Email: %s | Адреса: %s\n",
			c.ID, c.Name, c.Phone, c.Email, c.Address)
	}
}

func updateCustomer(store *Store) {
	if len(store.Customers) == 0 {
		fmt.Println("\nНемає клієнтів для оновлення.")
		return
	}

	viewCustomers(store)
	id, err := getRequiredInt("\nВведіть ID клієнта для оновлення: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}

	index := -1
	for i, c := range store.Customers {
		if c.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("Клієнта не знайдено.")
		return
	}

	fmt.Println("\n--- Оновлення даних клієнта ---")
	newName, _ := getOptionalString("Нове ім’я (Enter — без змін): ")
	newPhone, _ := getOptionalString("Новий телефон (Enter — без змін): ")
	newEmail, _ := getOptionalString("Новий e-mail (Enter — без змін): ")
	newAddress, _ := getOptionalString("Нова адреса (Enter — без змін): ")

	if newName != "" {
		store.Customers[index].Name = newName
	}
	if newPhone != "" {
		store.Customers[index].Phone = newPhone
	}
	if newEmail != "" {
		store.Customers[index].Email = newEmail
	}
	if newAddress != "" {
		store.Customers[index].Address = newAddress
	}

	fmt.Println("Дані клієнта оновлено успішно!")
}

// ----- КОШИК ПОКУПОК -----

func displayCartsMenu(store *Store) {
	fmt.Println("\n=== МЕНЮ КОШИКА ===")

	customerID, err := getRequiredInt("Введіть ID клієнта: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}

	customer := findCustomerByID(store, customerID)
	if customer == nil {
		fmt.Println("Клієнта не знайдено.")
		return
	}

	// Якщо кошика ще немає — створюємо
	if _, exists := store.Carts[customerID]; !exists {
		store.Carts[customerID] = *NewCart(customerID)
	}

	for {
		fmt.Printf("\n--- Кошик клієнта %s ---\n", customer.Name)
		fmt.Println("1. Додати товар до кошика")
		fmt.Println("2. Видалити товар з кошика")
		fmt.Println("3. Переглянути вміст кошика")
		fmt.Println("4. Застосувати знижку")
		fmt.Println("5. Очистити кошик")
		fmt.Println("6. Повернутися до головного меню")

		choice := getMenuChoice(6)
		switch choice {
		case 1:
			addToCart(store, customerID)
		case 2:
			removeFromCart(store, customerID)
		case 3:
			viewCart(store, customerID)
		case 4:
			applyDiscount(store, customerID)
		case 5:
			clearCart(store, customerID)
		case 6:
			return
		default:
			continue
		}
	}
}

// ---- Допоміжні функції ----

func findCustomerByID(store *Store, id int) *Customer {
	for i := range store.Customers {
		if store.Customers[i].ID == id {
			return &store.Customers[i]
		}
	}
	return nil
}

// Створення нового кошика
func NewCart(customerID int) *Cart {
	return &Cart{
		CustomerID: customerID,
		Items:      make(map[int]CartItem),
		Discount:   0,
	}
}

func addToCart(store *Store, customerID int) {
	id, err := getRequiredInt("Введіть ID товару: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}

	product := findProductByID(store, id)
	if product == nil || !product.IsActive {
		fmt.Println("Товар не знайдено або неактивний.")
		return
	}

	qty, err := getRequiredInt("Введіть кількість: ")
	if err != nil || qty <= 0 {
		fmt.Println("Некоректна кількість.")
		return
	}
	if qty > product.Stock {
		fmt.Printf("На складі лише %d шт.\n", product.Stock)
		return
	}

	cart := store.Carts[customerID]
	item, exists := cart.Items[id]
	if exists {
		item.Quantity += qty
	} else {
		item = CartItem{
			ProductID: id,
			Name:      product.Name,
			UnitPrice: product.Price,
			Quantity:  qty,
		}
	}
	cart.Items[id] = item
	store.Carts[customerID] = cart

	fmt.Printf("Товар \"%s\" (%d шт.) додано до кошика!\n", product.Name, qty)
}

func removeFromCart(store *Store, customerID int) {
	cart := store.Carts[customerID]
	if len(cart.Items) == 0 {
		fmt.Println("Кошик порожній.")
		return
	}

	viewCart(store, customerID)
	id, err := getRequiredInt("Введіть ID товару для видалення: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}

	if _, exists := cart.Items[id]; exists {
		delete(cart.Items, id)
		store.Carts[customerID] = cart
		fmt.Println("Товар видалено з кошика.")
	} else {
		fmt.Println("Товар із таким ID не знайдено у кошику.")
	}
}

func viewCart(store *Store, customerID int) {
	cart := store.Carts[customerID]
	if len(cart.Items) == 0 {
		fmt.Println("Кошик порожній.")
		return
	}

	customer := findCustomerByID(store, customerID)
	fmt.Printf("\n--- Кошик клієнта: %s ---\n", customer.Name)

	total := 0.0
	for _, item := range cart.Items {
		lineTotal := item.UnitPrice * float64(item.Quantity)
		total += lineTotal
		fmt.Printf("%d. %s x%d - %.2f €\n", item.ProductID, item.Name, item.Quantity, lineTotal)
	}

	discounted := total * (1 - cart.Discount/100)
	fmt.Printf("\nЗнижка: %.0f%%\n", cart.Discount)
	fmt.Printf("Загальна сума: %.2f €\n", discounted)
}

func applyDiscount(store *Store, customerID int) {
	cart := store.Carts[customerID]
	if len(cart.Items) == 0 {
		fmt.Println("Кошик порожній, знижку застосувати неможливо.")
		return
	}

	discount, err := getRequiredFloat("Введіть знижку у % (0–50): ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}
	if discount < 0 || discount > 50 {
		fmt.Println("Знижка повинна бути в межах 0–50%.")
		return
	}

	cart.Discount = discount
	store.Carts[customerID] = cart
	fmt.Printf("Знижку %.0f%% застосовано!\n", discount)
}

func clearCart(store *Store, customerID int) {
	cart := store.Carts[customerID]
	if len(cart.Items) == 0 {
		fmt.Println("Кошик уже порожній.")
		return
	}
	cart.Items = make(map[int]CartItem)
	cart.Discount = 0
	store.Carts[customerID] = cart
	fmt.Println("Кошик очищено.")
}

// ---------- Система замовлень ----------

func displayOrdersMenu(store *Store) {
	for {
		fmt.Println("\n=== МЕНЮ ЗАМОВЛЕНЬ ===")
		fmt.Println("1. Створити замовлення з кошика")
		fmt.Println("2. Переглянути всі замовлення")
		fmt.Println("3. Переглянути замовлення клієнта")
		fmt.Println("4. Змінити статус замовлення")
		fmt.Println("5. Повернутися до головного меню")

		choice := getMenuChoice(5)
		switch choice {
		case 1:
			createOrderFromCart(store)
		case 2:
			viewAllOrders(store)
		case 3:
			viewOrdersByCustomer(store)
		case 4:
			updateOrderStatus(store)
		case 5:
			return
		default:
			continue
		}
	}
}

func createOrderFromCart(store *Store) {
	if len(store.Customers) == 0 {
		fmt.Println("Немає клієнтів. Зареєструйте хоча б одного.")
		return
	}

	id, err := getRequiredInt("Введіть ID клієнта: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}

	cart, exists := store.Carts[id]
	if !exists || len(cart.Items) == 0 {
		fmt.Println("Кошик порожній або не знайдено.")
		return
	}

	subtotal := 0.0
	for _, item := range cart.Items {
		subtotal += float64(item.Quantity) * item.UnitPrice
	}
	discount := subtotal * cart.Discount / 100
	shipping := 150.00
	total := subtotal - discount + shipping

	fmt.Printf("\n--- Оформлення замовлення ---\n")
	fmt.Printf("Проміжна сума: %.2f грн\n", subtotal)
	fmt.Printf("Знижка: %.2f грн (%.0f%%)\n", discount, cart.Discount)
	fmt.Printf("Доставка: %.2f грн\n", shipping)
	fmt.Printf("Загальна сума до сплати: %.2f грн\n", total)

	answer, _ := getRequiredString("Підтвердити замовлення? (y/n): ")
	if strings.ToLower(answer) != "y" {
		fmt.Println("Замовлення скасовано.")
		return
	}

	order := Order{
		ID:         len(store.Orders) + 1,
		CustomerID: id,
		Items:      make([]OrderItem, 0),
		Subtotal:   subtotal,
		Discount:   discount,
		Shipping:   shipping,
		Total:      total,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ShippingTo: &ShippingInfo{
			Address: store.Customers[id-1].Address,
			Method:  "Стандартна доставка",
			Cost:    shipping,
		},
	}

	for _, item := range cart.Items {
		order.Items = append(order.Items, OrderItem{
			ProductID: item.ProductID,
			Name:      item.Name,
			UnitPrice: item.UnitPrice,
			Quantity:  item.Quantity,
			LineTotal: float64(item.Quantity) * item.UnitPrice,
		})
	}

	store.Orders = append(store.Orders, order)
	delete(store.Carts, id)

	fmt.Printf("Замовлення #%d створено успішно! Статус: %s\n", order.ID, order.Status)
}

func viewAllOrders(store *Store) {
	if len(store.Orders) == 0 {
		fmt.Println("Немає замовлень.")
		return
	}

	fmt.Println("\n--- Список усіх замовлень ---")
	for _, o := range store.Orders {
		fmt.Printf("Замовлення #%d | Клієнт #%d | Статус: %s | Сума: %.2f грн | Дата: %s\n",
			o.ID, o.CustomerID, o.Status, o.Total, o.CreatedAt.Format("2006-01-02 15:04"))
	}
}

func viewOrdersByCustomer(store *Store) {
	id, err := getRequiredInt("Введіть ID клієнта: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}

	found := false
	for _, o := range store.Orders {
		if o.CustomerID == id {
			if !found {
				fmt.Printf("\n--- Замовлення клієнта #%d ---\n", id)
				found = true
			}
			fmt.Printf("Замовлення #%d | Статус: %s | Сума: %.2f грн | Дата: %s\n",
				o.ID, o.Status, o.Total, o.CreatedAt.Format("2006-01-02 15:04"))
		}
	}

	if !found {
		fmt.Println("Цей клієнт не має замовлень.")
	}
}

func updateOrderStatus(store *Store) {
	if len(store.Orders) == 0 {
		fmt.Println("Немає замовлень для оновлення.")
		return
	}

	viewAllOrders(store)
	id, err := getRequiredInt("\nВведіть ID замовлення: ")
	if err != nil {
		fmt.Println("Помилка:", err)
		return
	}

	index := -1
	for i, o := range store.Orders {
		if o.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("Замовлення не знайдено.")
		return
	}

	fmt.Println("Можливі статуси: pending, shipped, delivered, canceled")
	status, _ := getRequiredString("Введіть новий статус: ")

	validStatuses := []string{"pending", "completed", "cancelled"}
	if !contains(validStatuses, strings.ToLower(status)) {
		fmt.Println("Невідомий статус.")
		return
	}

	store.Orders[index].Status = strings.ToLower(status)
	store.Orders[index].UpdatedAt = time.Now()

	fmt.Printf("Статус замовлення #%d оновлено на '%s'.\n", id, status)
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// ---------- Статистика магазину ----------

func displayStatistics(store *Store) {
	fmt.Println("\n=== СТАТИСТИКА МАГАЗИНУ ===")

	// 1. Товари
	totalProducts := len(store.Products)
	activeProducts := 0
	totalStockValue := 0.0
	for _, p := range store.Products {
		if p.IsActive {
			activeProducts++
			totalStockValue += float64(p.Stock) * p.Price
		}
	}

	fmt.Printf("Усього товарів: %d\n", totalProducts)
	fmt.Printf("Активних товарів: %d\n", activeProducts)
	fmt.Printf("Загальна вартість складу: %.2f €\n", totalStockValue)

	// 2. Клієнти
	totalCustomers := len(store.Customers)
	fmt.Printf("\nЗареєстровано клієнтів: %d\n", totalCustomers)

	// 3. Замовлення
	totalOrders := len(store.Orders)
	pending, completed, cancelled := 0, 0, 0
	for _, o := range store.Orders {
		switch strings.ToLower(o.Status) {
		case "pending":
			pending++
		case "completed":
			completed++
		case "cancelled":
			cancelled++
		}
	}

	fmt.Printf("\nЗамовлень створено: %d\n", totalOrders)
	fmt.Printf(" - Очікують: %d\n", pending)
	fmt.Printf(" - Виконані: %d\n", completed)
	fmt.Printf(" - Скасовані: %d\n", cancelled)

	// 4. Кошики
	fmt.Printf("\nАктивних кошиків: %d\n", len(store.Carts))

	fmt.Println("\n(Натисніть Enter для повернення)")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
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

func getRequiredString(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("рядок не може бути порожнім")
	}
	return input, nil
}

func getRequiredInt(prompt string) (int, error) {
	var n int
	fmt.Print(prompt)
	_, err := fmt.Scanln(&n)
	if err != nil {
		return 0, fmt.Errorf("необхідно ввести число")
	}
	return n, nil
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

// Опціональні (для оновлення)
func getOptionalString(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input), nil
}

func getOptionalInt(prompt string) (int, error) {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	if input == "" {
		return 0, nil
	}
	var n int
	fmt.Sscan(input, &n)
	return n, nil
}

func getOptionalFloat(prompt string) (float64, error) {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	if input == "" {
		return 0, nil
	}
	var f float64
	fmt.Sscan(input, &f)
	return f, nil
}

// Валідація пошти та номеру телефона
func validateEmail(email string) (bool, []string) {
	var errors []string

	if strings.Count(email, "@") == 0 {
		errors = append(errors, "Немає символа '@'")
	} else if strings.Count(email, "@") > 1 {
		errors = append(errors, "Більше ніж один символ '@'")
	}

	parts := strings.SplitN(email, "@", 2)
	if len(parts) == 2 {
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
		if !regexp.MustCompile(`^[A-Za-z0-9._-]+$`).MatchString(local) {
			errors = append(errors, "Недозволені символи у локальній частині")
		}
		if !strings.Contains(domain, ".") {
			errors = append(errors, "У доменній частині немає крапки '.'")
		} else {
			domainParts := strings.Split(domain, ".")
			tld := domainParts[len(domainParts)-1]
			if len(tld) < 2 || len(tld) > 6 {
				errors = append(errors, "Після останньої крапки має бути 2–6 символів")
			}
		}
	}

	return len(errors) == 0, errors
}

func validatePhone(phone string) (bool, []string) {
	var errors []string
	digits := 0

	if !strings.HasPrefix(phone, "+") {
		errors = append(errors, "Номер має починатися з '+'")
	}

	for _, ch := range phone {
		if unicode.IsDigit(ch) {
			digits++
		}
	}

	if digits < 10 || digits > 15 {
		errors = append(errors, "Кількість цифр має бути від 10 до 15")
	}

	return len(errors) == 0, errors
}
