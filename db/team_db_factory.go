package db

import (
	"database/sql"
	"errors"
)

//go:generate counterfeiter . TeamDBFactory

type TeamDBFactory interface {
	GetTeamDBById(int) TeamDB
	GetTeamDBByName(string) (TeamDB, error)
}

type teamDBFactory struct {
	conn        Conn
	bus         *notificationsBus
	lockFactory LockFactory
}

func NewTeamDBFactory(conn Conn, bus *notificationsBus, lockFactory LockFactory) TeamDBFactory {
	return &teamDBFactory{
		conn:        conn,
		bus:         bus,
		lockFactory: lockFactory,
	}
}

func (f *teamDBFactory) GetTeamDBById(teamId int) TeamDB {
	return &teamDB{
		teamId:       teamId,
		conn:         f.conn,
		buildFactory: newBuildFactory(f.conn, f.bus, f.lockFactory),
	}
}

func (f *teamDBFactory) GetTeamDBByName(teamName string) (TeamDB, error) {
	var id int
	row := f.conn.QueryRow(`
	SELECT id FROM teams WHERE LOWER(name) = LOWER($1)
`, teamName)

	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("This team does not exist in db")
		}

		return nil, err
	}

	team := f.GetTeamDBById(id)
	return team, nil
}
