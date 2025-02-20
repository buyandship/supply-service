package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/buyandship/supply-svr/biz/common/config"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
)

var (
	once    sync.Once
	Handler *H
)

type H struct {
	cli *gorm.DB
}

func GetHandler() *H {
	once.Do(func() {
		c, err := gorm.Open(mysql.New(
			mysql.Config{
				DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
					config.GlobalServerConfig.Mysql.Username,
					config.GlobalServerConfig.Mysql.Password,
					config.GlobalServerConfig.Mysql.Address,
					config.GlobalServerConfig.Mysql.DBName,
				),
			},
		), &gorm.Config{})
		if err != nil {
			hlog.Fatal("mysql init err:", err)
		}
		Handler = &H{cli: c}
	})
	return Handler
}

func (h *H) HealthCheck() error {
	sqlDB, err := h.cli.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Ping(); err != nil {
		return err
	}
	return nil
}

func (h *H) UpsertAccount(ctx context.Context, account *model.Account) error {
	if err := h.cli.
		WithContext(ctx).
		Debug().
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&account).Error; err != nil {
		return err
	}
	return nil
}

func (h *H) InsertMessage(ctx context.Context, message *model.Message) error {
	if err := h.cli.
		WithContext(ctx).
		Debug().
		Create(&message).Error; err != nil {
		return err
	}
	return nil
}

func (h *H) GetTransaction(ctx context.Context, trxId string) (*model.Transaction, error) {
	var trx *model.Transaction
	if err := h.cli.WithContext(ctx).
		Debug().
		Where("ref_id = ?", trxId).
		First(&trx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return trx, nil
}

func (h *H) InsertTransaction(ctx context.Context, transaction *model.Transaction) error {
	if err := h.cli.
		WithContext(ctx).
		Debug().
		Create(&transaction).Error; err != nil {
		return err
	}
	return nil
}

func (h *H) InsertTokenLog(ctx context.Context, token *model.Token) error {
	if err := h.cli.
		WithContext(ctx).
		Debug().
		Create(&token).Error; err != nil {
		return err
	}
	return nil
}

func (h *H) UpdateTransaction(ctx context.Context, cond *model.Transaction) error {
	if err := h.cli.
		WithContext(ctx).
		Debug().
		Where("ref_id = ?", cond.RefID).
		Updates(cond).Error; err != nil {
		return err
	}
	return nil
}

func (h *H) GetAccount(ctx context.Context, buyerID int32) (*model.Account, error) {
	account := &model.Account{}
	if err := h.cli.
		WithContext(ctx).
		Debug().
		Where("buyer_id = ?", buyerID).
		First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (h *H) GetToken() (*model.Token, error) {
	var token model.Token
	if err := h.cli.
		Order("created_at desc").
		First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}
