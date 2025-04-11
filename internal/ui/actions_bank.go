package ui

import (
	"banksystem/internal/model"
	"fmt"
)

func (n *NavigationManager) initializeBankPageData() error {
	if n.state.Banks.IsInitialized {
		return nil
	}

	banks, err := n.bankingService.GetBanks()
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}

	if len(banks) == 0 || len(banks) > 3 {
		return fmt.Errorf("ошибка: Меньше трех банков в хранилище")
	}

	n.state.Banks.SetBanks(banks)
	return nil
}

func (n *NavigationManager) initializeAdminPageData(bankId int) error {
	transactions, err := n.bankingService.GetTransactions(bankId)
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}
	n.state.Banks.AdminTransactionsList = make([]*model.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		if transaction.Status == "pending" {
			n.state.Banks.AdminTransactionsList = append(n.state.Banks.AdminTransactionsList, transaction)
		}
	}

	credits, err := n.bankingService.GetCredits(bankId)
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}
	n.state.Banks.AdminCreditList = make([]*model.Credit, 0, len(credits))
	for _, credit := range credits {
		if credit.Status == "pending" {
			n.state.Banks.AdminCreditList = append(n.state.Banks.AdminCreditList, credit)
		}
	}
	return nil
}

func (n *NavigationManager) initializeUserPageData(bankId int, userId int) error {
	transactions, err := n.bankingService.GetTransactions(bankId)
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}
	n.state.Banks.AdminTransactionsList = make([]*model.Transaction, 0, len(transactions))
	for _, tratransaction := range transactions {
		if tratransaction.InitiatedByUserId == userId {
			n.state.User.UserTransactionList = append(n.state.User.UserTransactionList, tratransaction)
		}
	}

	credits, err := n.bankingService.GetCredits(bankId)
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}
	n.state.Banks.AdminCreditList = make([]*model.Credit, 0, len(credits))
	for _, credit := range credits {
		if credit.InitiatedByUserId == userId {
			n.state.User.UserCreditList = append(n.state.User.UserCreditList, credit)
		}
	}
	return nil
}

func (n *NavigationManager) openTransactionPage() {
	n.navigateTo(ScreenTransaction)
}

func (n *NavigationManager) openCreditPage() {
	n.navigateTo(ScreenCredit)
}
func (n *NavigationManager) openLogsPage() {
	n.navigateTo(ScreenWatchLogs)
}

func (n *NavigationManager) openBankPage(index int) {
	n.state.Banks.SelectedBankIndex = index
	n.navigateTo(ScreenBank)
}

func (n *NavigationManager) openAdminMain(index int) {
	n.state.Banks.SelectedBankIndex = index
	n.navigateTo(ScreenAdminMain)
}
