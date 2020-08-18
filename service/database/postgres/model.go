package postgres

import (
	"database/sql"
	"database/sql/driver"
	billscan "github.com/okcredit/billscan/api/go"
	"github.com/okcredit/go-common/encoding/json"
	"github.com/okcredit/go-common/errors"
	"time"
)

type UserDBModel struct {
	Id           string         `db:"id"`
	Mobile       sql.NullString `db:"mobile"`
	ProfileImage sql.NullString `db:"profile_image"`
	DisplayName  sql.NullString `db:"display_name"`
	Registered   bool           `db:"registered"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
	SignImage    sql.NullString `db:"sign_image"`
}

type BillDBModel struct {
	Id        string            `db:"id"`
	ImageUrl  string            `db:"image_url"`
	BillFrom  string            `db:"bill_from"`
	BillTo    string            `db:"bill_to"`
	CreatedAt time.Time         `db:"created_at"`
	UpdatedAt time.Time         `db:"updated_at"`
	Agreement *AgreementDBModel `db:"agreement"`
}

type RelationDBModel struct {
	Id             string                     `db:"id"`
	UserId         string                     `db:"user_id"`
	ContactId      string                     `db:"contact_id"`
	CreatedAt      time.Time                  `db:"created_at"`
	UpdatedAt      time.Time                  `db:"updated_at"`
	ContactProfile *billscan.Relation_Profile `db:"contact_profile"`
}

func newUserDBModel(user *billscan.User) *UserDBModel {
	return &UserDBModel{
		Id:           user.Id,
		Mobile:       nullableString(user.Mobile),
		ProfileImage: nullableString(user.ProfileImage),
		DisplayName:  nullableString(user.DisplayName),
		Registered:   user.Registered,
		SignImage:    nullableString(user.SignImage),
	}
}

func newRelationDBModel(relation *billscan.Relation) *RelationDBModel {
	return &RelationDBModel{
		Id:             relation.Id,
		UserId:         relation.UserId,
		ContactId:      relation.ContactId,
		ContactProfile: relation.ContactProfile,
	}
}

func newBillDBModel(bill *billscan.Bill) *BillDBModel {
	if bill.Agreement == nil {
		bill.Agreement = &billscan.Agreement{}
	}

	return &BillDBModel{
		Id:       bill.Id,
		ImageUrl: bill.ImageUrl,
		BillFrom: bill.BillFrom,
		BillTo:   bill.BillTo,
		Agreement: &AgreementDBModel{
			Note:   bill.Agreement.Note,
			Status: bill.Agreement.Status,
		},
	}
}

func (relation *RelationDBModel) toDomainModel() *billscan.Relation {
	return &billscan.Relation{
		Id:             relation.Id,
		UserId:         relation.UserId,
		ContactId:      relation.ContactId,
		ContactProfile: relation.ContactProfile,
		CreatedAt:      uint64(relation.CreatedAt.Unix()),
		UpdatedAt:      uint64(relation.UpdatedAt.Unix()),
	}
}

func (bill *BillDBModel) toDomainModel() *billscan.Bill {

	if bill.Agreement == nil {
		bill.Agreement = &AgreementDBModel{}
	}

	return &billscan.Bill{
		Id:        bill.Id,
		ImageUrl:  bill.ImageUrl,
		BillFrom:  bill.BillFrom,
		BillTo:    bill.BillTo,
		CreatedAt: uint64(bill.CreatedAt.Unix()),
		UpdatedAt: uint64(bill.UpdatedAt.Unix()),
		Agreement: &billscan.Agreement{
			Note:   bill.Agreement.Note,
			Status: bill.Agreement.Status,
		},
	}
}

func (user *UserDBModel) toDomainModel() *billscan.User {
	return &billscan.User{
		Id:           user.Id,
		Mobile:       user.Mobile.String,
		ProfileImage: user.ProfileImage.String,
		DisplayName:  user.DisplayName.String,
		Registered:   user.Registered,
		CreatedAt:    uint64(user.CreatedAt.Unix()),
		UpdatedAt:    uint64(user.UpdatedAt.Unix()),
		SignImage:    user.SignImage.String,
	}
}

type AgreementDBModel billscan.Agreement

func (a AgreementDBModel) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AgreementDBModel) Scan(src interface{}) error {
	jsonData, ok := src.([]byte)
	if !ok {
		return errors.New("unexpected type in database column")
	}

	json.Unmarshal(jsonData, a)
	return nil
}

type ContactProfileDBModel billscan.Relation_Profile

func (m ContactProfileDBModel) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *ContactProfileDBModel) Scan(src interface{}) error {
	jsonData, ok := src.([]byte)
	if !ok {
		return errors.New("unexpected type in database column")
	}

	json.Unmarshal(jsonData, m)
	return nil
}

func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	} else {
		return sql.NullString{String: s, Valid: true}
	}
}
