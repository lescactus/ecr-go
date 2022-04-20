package summary

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepositoryFailedUpdate(t *testing.T) {
	var r1 RepositoryFailedUpdate
	var r2 RepositoryFailedUpdate

	r1 = NewRepositoryFailedUpdate()

	m := make(map[string]error)
	r2.ErrorRepositoryName = m

	assert.Equal(t, r1, r2)
}

func TestNewRepositorySuccededUpdate(t *testing.T) {
	var r1 RepositorySuccededUpdate
	var r2 RepositorySuccededUpdate

	r1 = NewRepositorySuccededUpdate()

	assert.Equal(t, r1, r2)
}

func TestRFAdd(t *testing.T) {
	tests := []struct {
		desc string
		want RepositoryFailedUpdate
	}{
		{
			desc: "Add entries to ErrorRepositoryName hashmap",
			want: RepositoryFailedUpdate{
				ErrorRepositoryName: map[string]error{
					"repoName1": errors.New("repoName1"),
					"repoName2": errors.New("repoName2"),
					"repoName3": errors.New("repoName3"),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			r := NewRepositoryFailedUpdate()
			r.Add("repoName1", errors.New("repoName1"))
			r.Add("repoName2", errors.New("repoName2"))
			r.Add("repoName3", errors.New("repoName3"))
			assert.Equal(t, test.want, r)
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		desc string
		want map[string]error
	}{
		{
			desc: "Read entries from errorRepositoryName hashmap",
			want: map[string]error{
				"repoName1": errors.New("repoName1"),
				"repoName2": errors.New("repoName2"),
				"repoName3": errors.New("repoName3"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			r := NewRepositoryFailedUpdate()
			r.Add("repoName1", errors.New("repoName1"))
			r.Add("repoName2", errors.New("repoName2"))
			r.Add("repoName3", errors.New("repoName3"))
			m := r.GetAll()
			assert.Equal(t, test.want, m)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		desc string
		key  string
		want error
	}{
		{
			desc: "Read first entry from errorRepositoryName hashmap",
			key:  "repoName1",
			want: errors.New("errorName1"),
		},
		{
			desc: "Read second entry from errorRepositoryName hashmap",
			key:  "repoName2",
			want: errors.New("errorName2"),
		},
		{
			desc: "Read third entry from errorRepositoryName hashmap",
			key:  "repoName3",
			want: errors.New("errorName3"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			r := NewRepositoryFailedUpdate()
			r.Add("repoName1", errors.New("errorName1"))
			r.Add("repoName2", errors.New("errorName2"))
			r.Add("repoName3", errors.New("errorName3"))

			e := r.Get(test.key)
			assert.Error(t, test.want, e)
			assert.EqualError(t, e, test.want.Error())
		})
	}
}

func TestRSAdd(t *testing.T) {
	tests := []struct {
		desc string
		want RepositorySuccededUpdate
	}{
		{
			desc: "Add entries to RepositoryNames slice",
			want: RepositorySuccededUpdate{
				RepositoryNames: []string{
					"repoName1",
					"repoName2",
					"repoName3",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			r := NewRepositorySuccededUpdate()
			r.Add("repoName1")
			r.Add("repoName2")
			r.Add("repoName3")
			assert.Equal(t, test.want, r)
		})
	}
}
