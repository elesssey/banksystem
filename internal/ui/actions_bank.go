package ui

import (
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
