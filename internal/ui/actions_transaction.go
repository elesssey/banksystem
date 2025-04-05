package ui

import (
	"banksystem/internal/model"
	"fmt"
	"log"
)

func (n *NavigationManager) createTransaction() error {
	err := n.bankingService.CreateTransaction(n.state.Transaction.BuildTransaction())
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}
	n.navigateTo(ScreenBank)
	return nil
}

func (n *NavigationManager) onCreateTransactionError(err error) {
	n.showError(err.Error(), func() {
		n.state.Transaction = nil
		n.navigateTo(ScreenBank)
	})
}

func (n *NavigationManager) adminConfirmationTransaction(transaction *model.Transaction) error {
	log.Printf("sqweqwe")
	n.bankingService.TransactionConfirmation(transaction.Id)
	n.navigateTo(ScreenAdminMain)
	return nil
}
