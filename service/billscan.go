package billscan

import (
	"context"
	"database/sql"
	"github.com/golang/protobuf/ptypes/empty"
	billscan "github.com/okcredit/billscan/api/go"
	"github.com/okcredit/go-common/errors"
)

var _ billscan.APIServer = &Service{}

type Service struct {
	Database Database
}

func (s Service) ListContacts(ctx context.Context, request *billscan.ListContactsRequest) (*billscan.ListContactsResponse, error) {
	customers, err := s.Database.ListContacts(ctx, request.UserId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &billscan.ListContactsResponse{Contacts: customers}, nil
}

func (s Service) AddContact(ctx context.Context, request *billscan.AddContactRequest) (*billscan.AddContactResponse, error) {
	if request.Contact.Mobile == "" {
		return nil, errors.New("invalid mobile")
	}

	customer, err := s.Database.GetUserByMobile(ctx, request.Contact.Mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			newUser, err := s.Database.CreateUser(ctx, &billscan.User{
				Mobile: request.Contact.Mobile,
			})
			if err != nil {
				return nil, err
			}

			request.Contact.Id = newUser.Id
		} else {
			return nil, err
		}
	} else {
		request.Contact.Id = customer.Id
	}

	if request.Contact.Id == request.UserId {
		return nil, errors.From(409, "contact same as user")
	}

	customer, err = s.Database.AddContact(ctx, request.UserId, request.Contact)
	if err != nil {
		return nil, err
	}

	return &billscan.AddContactResponse{Contact: customer}, nil
}

func (s Service) DeleteContact(ctx context.Context, request *billscan.DeleteContactRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (s Service) ListBills(ctx context.Context, request *billscan.ListBillsRequest) (*billscan.ListBillsResponse, error) {

	bills, err := s.Database.ListBills(ctx, request.UserId, request.ContactId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &billscan.ListBillsResponse{Bills: bills}, nil
}

func (s Service) CreateBill(ctx context.Context, request *billscan.CreateBillRequest) (*billscan.CreateBillResponse, error) {
	bill, err := s.Database.CreateBill(ctx, request.Bill)
	if err != nil {
		return nil, err
	}

	return &billscan.CreateBillResponse{Bill: bill}, nil
}

func (s Service) UpdateBill(ctx context.Context, request *billscan.UpdateBillRequest) (*billscan.UpdateBillResponse, error) {
	bill, err := s.Database.UpdateBill(ctx, request.Bill)
	if err != nil {
		return nil, err
	}

	return &billscan.UpdateBillResponse{Bill: bill}, nil
}

func (s Service) DeleteBill(ctx context.Context, request *billscan.DeleteBillRequest) (*empty.Empty, error) {
	err := s.Database.DeleteBill(ctx, request.BillId)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s Service) GetContact(ctx context.Context, request *billscan.GetContactRequest) (*billscan.GetContactResponse, error) {
	customer, err := s.Database.GetContact(ctx, request.UserId, request.ContactId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.From(404, "contact not found")
		}
		return nil, err
	}

	return &billscan.GetContactResponse{Contact: customer}, nil
}

func (s Service) UpdateContact(ctx context.Context, request *billscan.UpdateContactRequest) (*billscan.UpdateContactResponse, error) {
	customer, err := s.Database.UpdateContact(ctx, request.UserId, request.Contact)
	if err != nil {
		return nil, err
	}

	return &billscan.UpdateContactResponse{Contact: customer}, nil
}

func (s Service) GetUser(ctx context.Context, request *billscan.GetUserRequest) (*billscan.GetUserResponse, error) {
	user, err := s.Database.GetUser(ctx, request.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.From(404, "user not found")
		}
		return nil, err
	}

	return &billscan.GetUserResponse{User: user}, nil
}

func (s Service) CreateUser(ctx context.Context, request *billscan.CreateUserRequest) (*billscan.CreateUserResponse, error) {
	_, err := s.Database.GetUserByMobile(ctx, request.User.Mobile)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	} else {
		return nil, errors.From(409, "User with mobile already exists")
	}

	request.User.Registered = true

	user, err := s.Database.CreateUser(ctx, request.User)
	if err != nil {
		return nil, err
	}

	return &billscan.CreateUserResponse{User: user}, nil
}

func (s Service) UpdateUser(ctx context.Context, request *billscan.UpdateUserRequest) (*billscan.UpdateUserResponse, error) {
	request.User.Id = request.UserId
	user, err := s.Database.UpdateUser(ctx, request.User)
	if err != nil {
		return nil, err
	}

	return &billscan.UpdateUserResponse{User: user}, nil
}

func (s Service) DeleteUser(ctx context.Context, request *billscan.DeleteUserRequest) (*empty.Empty, error) {
	err := s.Database.DeleteUser(ctx, request.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.From(404, "bill not found")
		}
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s Service) GetBill(ctx context.Context, request *billscan.GetBillRequest) (*billscan.GetBillResponse, error) {
	bill, err := s.Database.GetBill(ctx, request.BillId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.From(404, "bill not found")
		}
		return nil, err
	}

	return &billscan.GetBillResponse{Bill: bill}, nil
}

func (s Service) GetUserByMobile(ctx context.Context, request *billscan.GetUserByMobileRequest) (*billscan.GetUserByMobileResponse, error) {
	user, err := s.Database.GetUserByMobile(ctx, request.Mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.From(404, "user not found")
		}
		return nil, err
	}

	return &billscan.GetUserByMobileResponse{User: user}, nil
}

func (s Service) Login(ctx context.Context, request *billscan.LoginRequest) (*billscan.LoginResponse, error) {
	user, err := s.Database.GetUserByMobile(ctx, request.Mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			user, err = s.Database.CreateUser(ctx, &billscan.User{
				Mobile:      request.Mobile,
				DisplayName: request.Name,
				Registered:  true,
			})
			if err != nil {
				return nil, err
			}
			return &billscan.LoginResponse{User: user}, nil
		}

		return nil, err
	}

	if !user.Registered {
		user.Registered = true
		user, err = s.Database.UpdateUser(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	return &billscan.LoginResponse{User: user}, nil
}
