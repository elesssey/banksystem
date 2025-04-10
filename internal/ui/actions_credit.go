package ui

import (
	"fmt"
)

func (n *NavigationManager) createCredit() error {
	err := n.bankingService.CreateCredit(n.state.Credit.BuildCredit())
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}
	n.navigateTo(ScreenBank)
	return nil
}

func (n *NavigationManager) onCreateCreditError(err error) {
	n.showError(err.Error(), func() {
		n.state.Credit = nil
		n.navigateTo(ScreenBank)
	})
}

// func (n *NavigationManager) adminConfirmationTransaction(transaction *model.Transaction) error {
// 	n.bankingService.TransactionConfirmation(transaction.Id)
// 	n.navigateTo(ScreenBankSelector)
// 	return nil
// }

// func (n *NavigationManager) adminDeclineTransaction(transaction *model.Transaction) error {
// 	n.bankingService.TransactionDeclination(transaction.Id)
// 	n.navigateTo(ScreenBankSelector)
// 	return nil
// }
