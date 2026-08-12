package helper

import (
	"github.com/allape/golang/model"
	"github.com/allape/gophorward"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserCanAccessTo(user *gophorward.SimpleUser, gallery *model.Gallery, context *gin.Context, db *gorm.DB) (bool, error) {
	if gallery.IsPublic {
		return true, nil
	}

	if gallery.CreatedBy == user.ID {
		return true, nil
	}

	var count int64
	if err := db.Model(&model.UserGallery{}).Where("`user_id` = ? AND `gallery_id` = ?", user.ID, gallery.ID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
