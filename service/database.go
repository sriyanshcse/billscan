package billscan

import (
	"context"
	billscan "github.com/okcredit/billscan/api/go"
)

type Database interface {
	AddContact(ctx context.Context, userId string, contact *billscan.User) (*billscan.User, error)

	UpdateContact(ctx context.Context, userId string, contact *billscan.User) (*billscan.User, error)

	GetContact(ctx context.Context, userId, contactId string) (*billscan.User, error)

	DeleteContact(ctx context.Context, userId, contactId string) error

	ListContacts(ctx context.Context, userId string) ([]*billscan.User, error)

	CreateUser(ctx context.Context, user *billscan.User) (*billscan.User, error)

	UpdateUser(ctx context.Context, user *billscan.User) (*billscan.User, error)

	GetUser(ctx context.Context, userId string) (*billscan.User, error)

	GetUserByMobile(ctx context.Context, mobile string) (*billscan.User, error)

	DeleteUser(ctx context.Context, userId string) error

	CreateBill(ctx context.Context, bill *billscan.Bill) (*billscan.Bill, error)

	UpdateBill(ctx context.Context, bill *billscan.Bill) (*billscan.Bill, error)

	GetBill(ctx context.Context, billId string) (*billscan.Bill, error)

	ListBills(ctx context.Context, userId, contactId string) ([]*billscan.Bill, error)

	DeleteBill(ctx context.Context, billId string) error
}
