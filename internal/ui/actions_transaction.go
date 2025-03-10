package ui

import (
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

func (n *NavigationManager) onCreateTransactionError(err error) {
	n.showError(err.Error(), func() {
		n.state.Transaction = nil
		n.navigateTo(ScreenBank)
	})
}
