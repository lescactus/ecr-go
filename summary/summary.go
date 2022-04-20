package summary

import (
	"sync"
)

type RepositoryFailedUpdate struct {
	ErrorRepositoryName map[string]error // Hashmap to store the repositories names that failed to be updated and their error for the summary
	sync.RWMutex                         // Mutex to protect the hashmap from concurrent accesses
}

type RepositorySuccededUpdate struct {
	RepositoryNames []string // Slice to store the repositories names that were successfully updated for the summary
}

// NewRepositoryFailedUpdate instanciate a RepositoryFailedUpdate struct
// It returns a RepositoryFailedUpdate with an initialized errorRepositoryName map
func NewRepositoryFailedUpdate() RepositoryFailedUpdate {
	return RepositoryFailedUpdate{
		ErrorRepositoryName: make(map[string]error),
	}
}

// NewRepositorySuccededUpdate instanciate a RepositorySuccededUpdate
// It returns a RepositorySuccededUpdate
func NewRepositorySuccededUpdate() RepositorySuccededUpdate {
	return RepositorySuccededUpdate{}
}

// Add will add the faulted repository name and its associated error in the hashmap
func (r *RepositoryFailedUpdate) Add(repo string, e error) {
	r.Lock()
	defer r.Unlock()

	r.ErrorRepositoryName[repo] = e
}

// GetAll will return the errorRepositoryName
// It returns a map[string]error containing the faulted repositories names and their associated errors
func (r *RepositoryFailedUpdate) GetAll() map[string]error {
	r.RLock()
	defer r.RUnlock()

	return r.ErrorRepositoryName
}

// Get will retrieve the error of the given repository name
// It returns the error associated to the repository name
func (r *RepositoryFailedUpdate) Get(repository string) error {
	r.RLock()
	defer r.RUnlock()

	return r.ErrorRepositoryName[repository]
}

// Add will add the succeded repository in the slice
func (r *RepositorySuccededUpdate) Add(repository string) {
	r.RepositoryNames = append(r.RepositoryNames, repository)
}
