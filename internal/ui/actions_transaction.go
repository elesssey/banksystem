package ui

import (
	"banksystem/internal/model"
	"fmt"
)

func (n *NavigationManager) createTransaction() error {
	err := n.bankingService.CreateTransaction(n.state.Transaction.BuildTransaction())
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}
	n.navigateTo(ScreenBank)
	return nil
}
func (n *NavigationManager) freezeAccount() error {
	err := n.bankingService.FreezeAccount(n.state.Banks.WorkingAccount.ID)
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}
	n.navigateTo(ScreenBankSelector)
	return nil
}

func (n *NavigationManager) unFreezeAccount() error {
	err := n.bankingService.UnFreezeAccount(n.state.Banks.WorkingAccount.ID)
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}
	n.navigateTo(ScreenBankSelector)
	return nil
}

func (n *NavigationManager) onCreateTransactionError(err error) {
	n.showError(err.Error(), func() {
		n.state.Transaction = nil
		n.navigateTo(ScreenBank)
	})
}

func (n *NavigationManager) onFreezingAccountError(err error) {
	n.showError(err.Error(), func() {
		n.navigateTo(ScreenBankSelector)
	})
}

func (n *NavigationManager) adminConfirmationTransaction(transaction *model.Transaction) error {
	n.bankingService.TransactionConfirmation(transaction.Id)
	n.navigateTo(ScreenBankSelector)
	return nil
}

func (n *NavigationManager) adminDeclineTransaction(transaction *model.Transaction) error {
	n.bankingService.TransactionDeclination(transaction.Id)
	n.navigateTo(ScreenBankSelector)
	return nil
}

func (n *NavigationManager) adminConfirmationCredit(credit *model.Credit) error {
	n.bankingService.CreditConfirmation(credit.Id)
	n.navigateTo(ScreenBankSelector)
	return nil
}

func (n *NavigationManager) adminDeclineCredit(credit *model.Credit) error {
	n.bankingService.CreditDeclination(credit.Id)
	n.navigateTo(ScreenBankSelector)
	return nil
}
