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

	n.state.Banks.SetBanks(banks[0], banks[1], banks[2])
	return nil
}

func (n *NavigationManager) openTransactionPage() {
	n.navigateTo(ScreenTransaction)
}

func (n *NavigationManager) openBankPage(index int) {
	n.state.Banks.SelectedBankIndex = index
	n.navigateTo(ScreenBank)
}

func (n *NavigationManager) onCreateTransactionClick(tx model.Transaction) error {
	err := n.bankingService.CreateTransaction(tx)
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}
	return nil
}

func (n *NavigationManager) onCreateTransactionError(err error) {
	n.showError(err.Error(), func() { n.navigateTo(ScreenBank) })
}
