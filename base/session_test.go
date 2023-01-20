package base

import (
	"testing"
)

func TestSessionEqual(t *testing.T) {
	tests := []struct {
		name   string
		first  *Session
		second *Session
		want   bool
	}{
		{
			"Empty sessions",
			&Session{},
			&Session{},
			true,
		},
		{
			"Identical process id",
			&Session{Pid: 1},
			&Session{Pid: 1},
			true,
		},
		{
			"Different process id",
			&Session{Pid: 1},
			&Session{Pid: 2},
			false,
		},
		{
			"Identical users",
			&Session{User: "test"},
			&Session{User: "test"},
			true,
		},
		{
			"Different users",
			&Session{User: "test"},
			&Session{User: "random"},
			false,
		},
		{
			"Identical databases",
			&Session{Db: "test"},
			&Session{Db: "test"},
			true,
		},
		{
			"Different databases",
			&Session{Db: "test"},
			&Session{Db: "random"},
			false,
		},
		{
			"Identical users and databases",
			&Session{User: "test", Db: "test"},
			&Session{User: "test", Db: "test"},
			true,
		},
		{
			"Different users and same databases",
			&Session{User: "test_1", Db: "test"},
			&Session{User: "test_2", Db: "test"},
			false,
		},
		{
			"Different databases and same user",
			&Session{User: "test", Db: "test_1"},
			&Session{User: "test", Db: "test_2"},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.first.Equal(tc.second)
			if got != tc.want {
				t.Errorf("got %t; want %t", got, tc.want)
			} else {
				t.Logf("got %t; want %t", got, tc.want)
			}
		})
	}
}

func TestSessionInSlice(t *testing.T) {
	sessions := []*Session{
		{User: "test"},
		{User: "test_1"},
		{User: "test_2"},
		{User: "postgres"},
		{Db: "test"},
	}

	tests := []struct {
		name  string
		input *Session
		want  bool
	}{
		{
			"Empty session",
			&Session{},
			false,
		},
		{
			"Session with user in slice",
			&Session{User: "test"},
			true,
		},
		{
			"Session with user not in slice",
			&Session{User: "random"},
			false,
		},
		{
			"Session with db in slice",
			&Session{Db: "test"},
			true,
		},
		{
			"Session with db not in slice",
			&Session{Db: "random"},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.InSlice(sessions)
			if got != tc.want {
				t.Errorf("got %t; want %t", got, tc.want)
			} else {
				t.Logf("got %t; want %t", got, tc.want)
			}
		})
	}
}
