package cas

import (
	r "github.com/dancannon/gorethink"
)

type RethinkDBAdapter struct {
	session  *r.Session
	dbName string
}

func NewRethinkDBAdapter(c *CAS) (*RethinkDBAdapter, error) {
	// Database setup
	dbSession, err := r.Connect(r.ConnectOpts{
		Address:  c.Config["dbHost"],
		Database: c.Config["dbName"],
	})
	if err != nil {
		return nil, err
	}

	return &RethinkDBAdapter{dbSession, c.Config["dbName"]}, nil
}

func (db *RethinkDBAdapter) GetServiceByName(serviceName string) (*CASService, *CASServerError) {
	return &CASService{serviceName, "nobody@nowhere.net"}, nil
}

func (db *RethinkDBAdapter) FindUserByEmail(username string) (*User, *CASServerError) {
	// Find the user
	cursor, err := r.
		Db(db.dbName).
		Table("users").
		Get(username).
		Run(db.session)
	if err != nil {	return nil, &InvalidEmailAddressError	}

	// Get the user from the returned cursor
	var returnedUser *User
	err = cursor.One(&returnedUser)
	if err != nil {	return nil, &InvalidEmailAddressError	}

	return returnedUser, nil
}

func (db *RethinkDBAdapter) MakeNewTicketForService(service *CASService) (*CASTicket, *CASServerError) {
	// TODO
	return &CASTicket{}, nil
}

func (db *RethinkDBAdapter) RemoveTicketsForUser(username string, service *CASService) *CASServerError {
	return nil
}

func (db *RethinkDBAdapter) FindTicketForService(ticket string, service *CASService) (*CASTicket, *CASServerError) {
	return &CASTicket{}, nil
}

func (db *RethinkDBAdapter) AddNewUser(username, password string) (*User, *CASServerError) {
	user := &User{username, password}

	// Insert user into database
	res, err := r.
		Db(db.dbName).
		Table("users").
		Insert(user, r.InsertOpts{Conflict: "error"}).
		RunWrite(db.session)
	if err != nil {
		return nil, &FailedToCreateUserError
	} else if res.Errors > 0 {
		return nil, &EmailAlreadyTakenError
	}

	return user, nil
}
