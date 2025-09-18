package dao

import (
	"fmt"

	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

func DeferTransaction(tx *gorm.DB, log *logger.Logger, err **error, operation string) {
	if recoverErr := recover(); recoverErr != nil {
		log.Errorf("%s recover: %v", operation, recoverErr)
		**err = fmt.Errorf("%s %v", operation, recoverErr)
	}
	if **err != nil {
		if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
			**err = fmt.Errorf("%s: [%v]; tx.Rollback: [%v]", operation, **err, rollBackErr)
		}
	} else {
		if commitErr := tx.Commit().Error; commitErr != nil {
			log.Errorf("%s tx.Commit: %v", operation, commitErr)
			**err = commitErr
		}
	}
}
