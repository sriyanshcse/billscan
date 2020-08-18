package postgres

import (
	"context"
	"github.com/okcredit/nap"
	uuid "github.com/satori/go.uuid"
	"time"

	//"database/sql/driver"
	_ "github.com/lib/pq"
	billscan "github.com/okcredit/billscan/api/go"
	service "github.com/okcredit/billscan/service"
)

func New(db *nap.DB) service.Database {
	return &database{db: db}
}

var _ service.Database = &database{}

type database struct {
	db *nap.DB
}

func (d *database) AddContact(ctx context.Context, userId string, contact *billscan.User) (*billscan.User, error) {
	tsNow := time.Now()

	_, err := d.db.ExecContext(ctx, `insert into relations(id, user_id, contact_id, contact_profile, created_at, updated_at)
								  values
  							      ($1, $3, $4, $5, $6, $7),
								  ($2, $4, $3, $8::jsonb, $6, $7)`,
		uuid.NewV4().String(), uuid.NewV4().String(), userId, contact.Id, &ContactProfileDBModel{DisplayName: contact.DisplayName}, tsNow, tsNow, &ContactProfileDBModel{})

	if err != nil {
		return nil, err
	}

	contact.CreatedAt = uint64(tsNow.Unix())
	contact.UpdatedAt = uint64(tsNow.Unix())

	return contact, err
}

func (d *database) UpdateContact(ctx context.Context, userId string, contact *billscan.User) (*billscan.User, error) {
	panic("implement me")
}

func (d *database) GetContact(ctx context.Context, userId, contactId string) (*billscan.User, error) {
	contact := &UserDBModel{}
	profile := &ContactProfileDBModel{}

	err := d.db.QueryRowContext(ctx, `select users.id, users.mobile, users.profile_image, relations.contact_profile, users.registered, users.created_at, users.updated_at from users join relations on users.id = relations.contact_id where relations.user_id = $1 and relations.contact_id = $2`, userId, contactId).
		Scan(&contact.Id, &contact.Mobile, &contact.ProfileImage, &profile, &contact.Registered, &contact.CreatedAt, &contact.UpdatedAt)

	contact.DisplayName = nullableString(profile.DisplayName)
	if err != nil {
		return nil, err
	}
	return contact.toDomainModel(), err
}

func (d *database) DeleteContact(ctx context.Context, userId, contactId string) error {
	panic("implement me")
}

func (d *database) ListContacts(ctx context.Context, userId string) ([]*billscan.User, error) {
	var contacts []*billscan.User

	rows, err := d.db.QueryContext(ctx, `select users.id, users.mobile, users.profile_image, relations.contact_profile, users.registered, users.created_at, users.updated_at from users join relations on users.id = relations.contact_id where relations.user_id = $1`, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		contact := &UserDBModel{}
		profile := &ContactProfileDBModel{}
		err := rows.Scan(&contact.Id, &contact.Mobile, &contact.ProfileImage, &profile, &contact.Registered, &contact.CreatedAt, &contact.UpdatedAt)
		if err != nil {
			return nil, err
		}

		contact.DisplayName = nullableString(profile.DisplayName)
		contacts = append(contacts, contact.toDomainModel())
	}

	return contacts, err
}

func (d *database) CreateUser(ctx context.Context, user *billscan.User) (*billscan.User, error) {

	userId := uuid.NewV4().String()
	tsNow := time.Now()

	userModel := newUserDBModel(user)

	_, err := d.db.ExecContext(ctx, `insert into users(id, mobile, profile_image, display_name, registered, sign_image, created_at, updated_at)
								  values($1, $2, $3, $4, $5, $6, $7, $8)`,
		userId, userModel.Mobile, userModel.ProfileImage, userModel.DisplayName, userModel.Registered, userModel.SignImage, tsNow, tsNow)

	if err != nil {
		return nil, err
	}

	user.CreatedAt = uint64(tsNow.Unix())
	user.UpdatedAt = uint64(tsNow.Unix())
	user.Id = userId
	return user, err
}

func (d *database) UpdateUser(ctx context.Context, user *billscan.User) (*billscan.User, error) {
	tsNow := time.Now()

	userModel := newUserDBModel(user)

	_, err := d.db.ExecContext(ctx, `update users
											set 
												profile_image = $2, 
												display_name = $3, 
												registered = $4, 
												sign_image = $5, 
												updated_at = $6
											where id = $1`,
		user.Id, userModel.ProfileImage, userModel.DisplayName, userModel.Registered, userModel.SignImage, tsNow)

	if err != nil {
		return nil, err
	}

	user.UpdatedAt = uint64(tsNow.Unix())
	return user, err
}

func (d *database) GetUser(ctx context.Context, userId string) (*billscan.User, error) {
	user := &UserDBModel{}
	err := d.db.QueryRowContext(ctx, `select id, mobile, profile_image, display_name, registered, sign_image, created_at, updated_at  from users where id = $1`, userId).
		Scan(&user.Id, &user.Mobile, &user.ProfileImage, &user.DisplayName, &user.Registered, &user.SignImage, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user.toDomainModel(), err
}

func (d *database) DeleteUser(ctx context.Context, userId string) error {
	panic("implement me")
}

func (d *database) CreateBill(ctx context.Context, bill *billscan.Bill) (*billscan.Bill, error) {
	billId := uuid.NewV4().String()
	tsNow := time.Now()

	billModel := newBillDBModel(bill)

	_, err := d.db.ExecContext(ctx, `insert into bills(id, image_url, bill_from, bill_to, agreement, created_at, updated_at)
								  values($1, $2, $3, $4, $5, $6, $7)`,
		billId, billModel.ImageUrl, billModel.BillFrom, billModel.BillTo, billModel.Agreement, tsNow, tsNow)

	if err != nil {
		return nil, err
	}

	bill.Id = billId
	bill.CreatedAt = uint64(tsNow.Unix())
	bill.UpdatedAt = uint64(tsNow.Unix())
	return bill, err
}

func (d *database) UpdateBill(ctx context.Context, bill *billscan.Bill) (*billscan.Bill, error) {
	tsNow := time.Now()

	billModel := newBillDBModel(bill)

	_, err := d.db.ExecContext(ctx, `update bills 
											set
												image_url = $2,  
												agreement = $3, 
												updated_at = $4
											where id = $1`,
		billModel.Id, billModel.ImageUrl, billModel.Agreement, tsNow)

	if err != nil {
		return nil, err
	}

	bill.UpdatedAt = uint64(tsNow.Unix())
	return bill, err
}

func (d *database) GetBill(ctx context.Context, billId string) (*billscan.Bill, error) {
	bill := &BillDBModel{}
	err := d.db.QueryRowContext(ctx, `select id, image_url, bill_from, bill_from, agreement, created_at, updated_at  from bills where id = $1`, billId).
		Scan(&bill.Id, &bill.ImageUrl, &bill.BillFrom, &bill.BillTo, &bill.Agreement, &bill.CreatedAt, &bill.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return bill.toDomainModel(), err
}

func (d *database) ListBills(ctx context.Context, userId, contactId string) ([]*billscan.Bill, error) {
	var bills []*billscan.Bill

	rows, err := d.db.QueryContext(ctx, `select id, image_url, bill_from, bill_to, agreement, created_at, updated_at from bills where bill_from = $1 and bill_to = $2 or bill_from = $2 and bill_to = $1`, userId, contactId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		bill := &BillDBModel{}
		err := rows.Scan(&bill.Id, &bill.ImageUrl, &bill.BillFrom, &bill.BillTo, &bill.Agreement, &bill.CreatedAt, &bill.UpdatedAt)
		if err != nil {
			return nil, err
		}

		bills = append(bills, bill.toDomainModel())
	}

	return bills, err
}

func (d *database) GetUserByMobile(ctx context.Context, mobile string) (*billscan.User, error) {
	user := &UserDBModel{}
	err := d.db.QueryRowContext(ctx, `select id, mobile, profile_image, display_name, registered, sign_image, created_at, updated_at from users where mobile = $1`, mobile).
		Scan(&user.Id, &user.Mobile, &user.ProfileImage, &user.DisplayName, &user.Registered, &user.SignImage, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user.toDomainModel(), err
}

func (d *database) DeleteBill(ctx context.Context, billId string) error {
	_, err := d.db.ExecContext(ctx, `delete from bills where id = $1`, billId)
	return err
}
