package terminator

import (
	"reflect"
	"testing"

	"github.com/jouir/pgterminate/base"
)

func TestFilterUsers(t *testing.T) {

	sessions := []*base.Session{
		{User: "test"},
		{User: "test_1"},
		{User: "test_2"},
		{User: "postgres"},
	}

	tests := []struct {
		name   string
		config *base.Config
		want   []*base.Session
	}{
		{
			"No filter",
			&base.Config{},
			sessions,
		},
		{
			"Include a single user",
			&base.Config{IncludeUsers: []string{"test"}},
			[]*base.Session{{User: "test"}},
		},
		{
			"Include multiple users",
			&base.Config{IncludeUsers: []string{"test_1", "test_2"}},
			[]*base.Session{{User: "test_1"}, {User: "test_2"}},
		},
		{
			"Exclude a single user",
			&base.Config{ExcludeUsers: []string{"test"}},
			[]*base.Session{{User: "test_1"}, {User: "test_2"}, {User: "postgres"}},
		},
		{
			"Exclude multiple users",
			&base.Config{ExcludeUsers: []string{"test_1", "test_2"}},
			[]*base.Session{{User: "test"}, {User: "postgres"}},
		},
		{
			"Include multiple users and exclude one",
			&base.Config{IncludeUsers: []string{"test", "test_1", "test_2"}, ExcludeUsers: []string{"test"}},
			[]*base.Session{{User: "test_1"}, {User: "test_2"}},
		},
		{
			"Include users from list and regex",
			&base.Config{
				IncludeUsers:      []string{"test"},
				IncludeUsersRegex: "^test_[0-9]$",
			},
			[]*base.Session{{User: "test"}, {User: "test_1"}, {User: "test_2"}},
		},
		{
			"Exclude users from list and regex",
			&base.Config{
				ExcludeUsers:      []string{"test"},
				ExcludeUsersRegex: "^test_[0-9]$",
			},
			[]*base.Session{{User: "postgres"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.CompileRegexes()
			if err != nil {
				t.Errorf("Failed to compile regex: %v", err)
			}
			tc.config.CompileFilters()
			terminator := &Terminator{config: tc.config}
			got := terminator.filterUsers(sessions)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v; want %+v", ListUsers(got), ListUsers(tc.want))
			} else {
				t.Logf("got %+v; want %+v", ListUsers(got), ListUsers(tc.want))
			}
		})
	}
}

// ListUsers extract usernames from a list of sessions
func ListUsers(sessions []*base.Session) (users []string) {
	for _, session := range sessions {
		users = append(users, session.User)
	}
	return users
}

func TestFilterDatabases(t *testing.T) {

	sessions := []*base.Session{
		{Db: "test"},
		{Db: "test_1"},
		{Db: "test_2"},
		{Db: "postgres"},
	}

	tests := []struct {
		name   string
		config *base.Config
		want   []*base.Session
	}{
		{
			"No filter",
			&base.Config{},
			sessions,
		},
		{
			"Include a single database",
			&base.Config{IncludeDatabases: []string{"test"}},
			[]*base.Session{{Db: "test"}},
		},
		{
			"Include multiple databases",
			&base.Config{IncludeDatabases: []string{"test_1", "test_2"}},
			[]*base.Session{{Db: "test_1"}, {Db: "test_2"}},
		},
		{
			"Exclude a single database",
			&base.Config{ExcludeDatabases: []string{"test"}},
			[]*base.Session{{Db: "test_1"}, {Db: "test_2"}, {Db: "postgres"}},
		},
		{
			"Exclude multiple databases",
			&base.Config{ExcludeDatabases: []string{"test_1", "test_2"}},
			[]*base.Session{{Db: "test"}, {Db: "postgres"}},
		},
		{
			"Include multiple databases and exclude one",
			&base.Config{IncludeDatabases: []string{"test", "test_1", "test_2"}, ExcludeDatabases: []string{"test"}},
			[]*base.Session{{Db: "test_1"}, {Db: "test_2"}},
		},
		{
			"Include databases from list and regex",
			&base.Config{
				IncludeDatabases:      []string{"test"},
				IncludeDatabasesRegex: "^test_[0-9]$",
			},
			[]*base.Session{{Db: "test"}, {Db: "test_1"}, {Db: "test_2"}},
		},
		{
			"Exclude databases from list and regex",
			&base.Config{
				ExcludeDatabases:      []string{"test"},
				ExcludeDatabasesRegex: "^test_[0-9]$",
			},
			[]*base.Session{{Db: "postgres"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.CompileRegexes()
			if err != nil {
				t.Errorf("Failed to compile regex: %v", err)
			}
			tc.config.CompileFilters()
			terminator := &Terminator{config: tc.config}
			got := terminator.filterDatabases(sessions)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v; want %+v", ListDatabases(got), ListDatabases(tc.want))
			} else {
				t.Logf("got %+v; want %+v", ListDatabases(got), ListDatabases(tc.want))
			}
		})
	}
}

// ListDatabases extract usernames from a list of sessions
func ListDatabases(sessions []*base.Session) (databases []string) {
	for _, session := range sessions {
		databases = append(databases, session.Db)
	}
	return databases
}
