//go:build !nullauth

package jasmine

import (
	"github.com/1f349/cardcaldav"
	"github.com/charmbracelet/log"
)

func NewAuth(db string, logger *log.Logger) AuthProvider {
	return cardcaldav.NewAuth(db, logger)
}
