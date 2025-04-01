package db

import (
	"context"
	"errors"
	"fmt"

	"sync"

	"github.com/buyandship/supply-svr/biz/common/config"
	"github.com/buyandship/supply-svr/biz/common/trace"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (h *H) UpsertAccount(ctx context.Context, account *model.Account) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "UpsertAccount")
	defer trace.EndSpan(span, err)

	err = h.cli.
		WithContext(ctx).
		Debug().
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&account).Error

	return err
}

func (h *H) InsertMessage(ctx context.Context, message *model.Message) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "InsertMessage")
	defer trace.EndSpan(span, err)

	err = h.cli.
		WithContext(ctx).
		Debug().
		Create(&message).Error

	return
}

func (h *H) InsertReview(ctx context.Context, review *model.Review) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "InsertReview")
	defer trace.EndSpan(span, err)

	err = h.cli.
		WithContext(ctx).
		Debug().
		Create(&review).Error

	return
}

func (h *H) GetTransaction(ctx context.Context, trxId string) (trx *model.Transaction, err error) {
	ctx, span := trace.StartDBOperation(ctx, "GetTransaction")
	defer trace.EndSpan(span, err)

	err = h.cli.WithContext(ctx).
		Debug().
		Where("ref_id = ?", trxId).
		First(&trx).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (h *H) InsertTransaction(ctx context.Context, transaction *model.Transaction) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "InsertTransaction")
	defer trace.EndSpan(span, err)

	err = h.cli.
		WithContext(ctx).
		Debug().
		Create(&transaction).Error

	return
}

func (h *H) InsertTokenLog(ctx context.Context, token *model.Token) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "InsertTokenLog")
	defer trace.EndSpan(span, err)

	err = h.cli.
		WithContext(ctx).
		Debug().
		Create(&token).Error

	return err
}

func (h *H) UpdateTransaction(ctx context.Context, cond *model.Transaction) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "UpdateTransaction")
	defer trace.EndSpan(span, err)

	err = h.cli.
		WithContext(ctx).
		Debug().
		Where("ref_id = ?", cond.RefID).
		Updates(cond).Error

	return
}

func (h *H) GetAccount(ctx context.Context, buyerID int32) (account *model.Account, err error) {
	ctx, span := trace.StartDBOperation(ctx, "GetAccount")
	defer trace.EndSpan(span, err)

	err = h.cli.
		WithContext(ctx).
		Debug().
		Where("buyer_id = ?", buyerID).
		First(account).Error

	if err != nil {
		return nil, err
	}
	return
}

func (h *H) GetToken(ctx context.Context) (token *model.Token, err error) {
	ctx, span := trace.StartDBOperation(ctx, "GetToken")
	defer trace.EndSpan(span, err)

	err = h.cli.
		WithContext(ctx).
		Order("created_at desc").
		First(&token).Error

	if err != nil {
		return nil, err
	}
	return
}
