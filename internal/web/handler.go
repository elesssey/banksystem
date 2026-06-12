package web

import (
	"banksystem/internal/model"
	"banksystem/internal/service"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Handler struct {
	bankingService service.BankingService
	authService    service.AuthService
	stateStore     *StateStore
	templates      *template.Template
}

func NewHandler(
	bankingService service.BankingService,
	authService service.AuthService,
	stateStore *StateStore,
	templates *template.Template,
) *Handler {
	return &Handler{
		bankingService: bankingService,
		authService:    authService,
		stateStore:     stateStore,
		templates:      templates,
	}
}

func (h Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	err := h.templates.ExecuteTemplate(w, "login_page.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.authService.Login(email, password)
	if err != nil {
		h.templates.ExecuteTemplate(w, "login_page.html", map[string]any{
			"Error": "Неверный email или пароль",
		})
		return
	}

	token, err := generateToken(user.ID, string(user.Role))
	if err != nil {
		http.Error(w, "Ошибка создания токена", http.StatusInternalServerError)
		return
	}

	setTokenCookie(w, token)

	if user.Role == "admin" {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/transaction", http.StatusSeeOther)
	}
}

func (h Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	banks, err := h.bankingService.GetBanks()
	if err != nil {
		http.Error(w, "Ошибка загрузки банков: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.templates.ExecuteTemplate(w, "register_page.html", map[string]any{
		"Banks": banks,
	})
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	middlename := r.FormValue("middlename")
	surname := r.FormValue("surname")
	password := r.FormValue("password")
	passportSeries := r.FormValue("passport_series")
	passportNumber := r.FormValue("passport_number")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	bankIdStr := r.FormValue("bank_id")

	if name == "" || surname == "" || password == "" || email == "" || bankIdStr == "" {
		banks, _ := h.bankingService.GetBanks()
		h.templates.ExecuteTemplate(w, "register_page.html", map[string]any{
			"Error": "Заполните все обязательные поля",
			"Banks": banks,
		})
		return
	}

	bankId, err := strconv.Atoi(bankIdStr)
	if err != nil {
		http.Error(w, "Неверный ID банка", http.StatusBadRequest)
		return
	}

	newUser := model.User{
		Name:           name,
		MiddleName:     middlename,
		Surname:        surname,
		Password:       password,
		PassportSeries: passportSeries,
		PassportNumber: passportNumber,
		Phone:          phone,
		Email:          email,
		Role:           "client",
	}

	err = h.authService.Registrate(&newUser, bankId)
	if err != nil {
		banks, _ := h.bankingService.GetBanks()
		h.templates.ExecuteTemplate(w, "register_page.html", map[string]any{
			"Error": "Ошибка при регистрации: " + err.Error(),
			"Banks": banks,
		})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := getTokenFromCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		claims, err := parseToken(tokenStr)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		reqCsrfEncoded := sha256.Sum256([]byte(tokenStr))
		reqCsrf := base64.StdEncoding.EncodeToString(reqCsrfEncoded[:])

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", int(claims.UserID))
		ctx = context.WithValue(ctx, "role", string(claims.Role))
		ctx = context.WithValue(ctx, "CSRF-TOKEN", reqCsrf)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h Handler) StartRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bankID, err := strconv.Atoi(r.FormValue("bank_id"))
	if err != nil {
		http.Error(w, "Некорректный bank_id", http.StatusBadRequest)
		return
	}

	user := &model.User{
		Name:           r.FormValue("name"),
		MiddleName:     r.FormValue("middlename"),
		Surname:        r.FormValue("surname"),
		Password:       r.FormValue("password"),
		PassportSeries: r.FormValue("passport_series"),
		PassportNumber: r.FormValue("passport_number"),
		Phone:          r.FormValue("phone"),
		Email:          r.FormValue("email"),
		Role:           model.Role(r.FormValue("role")),
	}

	err = h.authService.StartRegistration(user, bankID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/verify-email?email="+user.Email, http.StatusSeeOther)
}

func (h Handler) VerifyEmailPage(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	err := h.templates.ExecuteTemplate(w, "verify_email.html", map[string]any{
		"Email": email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")

	code := r.FormValue("code")

	err := h.authService.ConfirmRegistration(email, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h Handler) TransactionPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	val := r.Context().Value("user_id")
	userID, ok := val.(int)
	if !ok {
		http.Error(w, "invalid user id in context", http.StatusInternalServerError)
		return
	}
	userAccount, err := h.bankingService.GetAnyUserAccount(userID)

	if err != nil {
		http.Error(w, "Ошибка загрузки аккаунта пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	csrfVal := r.Context().Value("CSRF-TOKEN")
	csrf, _ := csrfVal.(string)

	err = h.templates.ExecuteTemplate(w, "transaction_page.html", map[string]any{
		"Account": userAccount,
		"Csrf":    csrf,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h Handler) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("qweqweqwe")

	val := r.Context().Value("user_id")
	userID, ok := val.(int)
	if !ok {
		http.Error(w, "invalid user id in context", http.StatusInternalServerError)
		return
	}
	userSourceAccount, err := h.bankingService.GetAnyUserAccount(userID)
	if err != nil {
		http.Error(w, "Ошибка загрузки аккаунта пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	amountStr := r.FormValue("amount")

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(w, "Некорректная сумма", http.StatusBadRequest)
		return
	}

	csrf := r.FormValue("CSRF-TOKEN")

	if csrf != r.Context().Value("CSRF-TOKEN") {
		http.Error(w, "Некорректный CSRF", http.StatusBadRequest)
		return
	}

	number := r.FormValue("to_account")

	userDestinationAccount, err := h.bankingService.GetAccountByNumber(number)

	if err != nil {
		http.Error(w, "Некорректный номер счета", http.StatusBadRequest)
		return
	}

	sourceUser, err := h.bankingService.GetUserById(userID)
	if err != nil {
		http.Error(w, "Ошибка загрузки отправителя", http.StatusBadRequest)
		return
	}

	destinationUser, err := h.bankingService.GetUserById(userDestinationAccount.UserId)
	if err != nil {
		http.Error(w, "Ошибка загрузки получателя", http.StatusBadRequest)
		return
	}

	transaction := &model.Transaction{
		Amount:                   amount,
		Сurrency:                 userSourceAccount.Currency,
		Status:                   "pending",
		SourceAccountId:          userSourceAccount.ID,
		DestinationAccountId:     userDestinationAccount.ID,
		SourceAccountType:        "user",
		DestinationAccountType:   "user",
		Type:                     "transfer",
		SourceBankId:             userSourceAccount.BankId,
		DestinationBankId:        userDestinationAccount.BankId,
		InitiatedByUserId:        userSourceAccount.UserId,
		SourseAccountNumber:      userSourceAccount.Number,
		DestinationAccountNumber: userDestinationAccount.Number,
		SourceAccountUser:        sourceUser,
		DestinationAccountUser:   destinationUser,
	}

	err = h.bankingService.CreationTransaction(transaction)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h Handler) MakeTransactionConfirmation(w http.ResponseWriter, r *http.Request) {

	roleVal := r.Context().Value("role")
	role, ok := roleVal.(string)
	if !ok || role != "admin" {
		http.Error(w, "У вас не хватает прав", http.StatusForbidden)
		return
	}

	transactionIdStr := r.URL.Query().Get("id")
	transactionId, err := strconv.Atoi(transactionIdStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	h.bankingService.TransactionConfirmation(transactionId)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h Handler) MakeTransactionDeclanation(w http.ResponseWriter, r *http.Request) {
	roleVal := r.Context().Value("role")
	role, ok := roleVal.(string)
	if !ok || role != "admin" {
		http.Error(w, "У вас не хватает прав", http.StatusForbidden)
		return
	}

	transactionIdStr := r.URL.Query().Get("id")
	transactionId, err := strconv.Atoi(transactionIdStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	h.bankingService.TransactionDeclination(transactionId)

	http.Redirect(w, r, "/transaction", http.StatusSeeOther)
}

func (h Handler) AdminPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	roleVal := r.Context().Value("role")
	role, ok := roleVal.(string)
	if !ok || role != "admin" {
		http.Error(w, "У вас не хватает прав", http.StatusForbidden)
		return
	}

	transactions, err := h.bankingService.GetAllTransactions()
	if err != nil {
		http.Error(w, "Ошибка загрузки транзакций: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.templates.ExecuteTemplate(w, "admin_page.html", map[string]any{
		"Transactions": transactions,
	})
}
