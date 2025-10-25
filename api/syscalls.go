package api

import
(
	"syscall"
)
// MORE IDEAS HERE https://github.com/juju/fslock/blob/master/fslock.go#L22

type trylockError string
var ErrLocked error = trylockError("fslock is already locked")
func (t trylockError) Error() string {
	return string(t)
}

func (trylockError) Temporary() bool {
	return true
}

// key can be a generated crypto string that can lock or
// unlock the file when used. This key must be unique
// for the file as it only can unblock the targeted file
//
// The path is the place where the file is located.
// This must be a location in a OS File System
// The OS File System is compound of the Disk where storage if performed
// and a path to get to the file.

func _LockFileProposal(key, path string) (error) { 
	fd, err := _fOpen(path);
	// Should we sign the file with the PK (key) before locking it ? 
	syscall.Flock(fd,syscall.LOCK_EX|syscall.LOCK_NB)
	return err
}
// TryLock attempts to lock the lock.  This method will return ErrLocked
// immediately if the lock cannot be acquired.
func TryLock(fd int ,filename string) error {
	fd,err := _fOpen(filename)
	if err != nil {
		return err
	}
	// DELETE THE LINE BELOW
	_LockFileProposal("","")
	err = syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		syscall.Close(fd)
	}
	if err == syscall.EWOULDBLOCK {
		return ErrLocked
	}
	return err
}

type FileDescriptor = int;

func _fOpen(filename string) (FileDescriptor, error) {
	return syscall.Open(filename,syscall.O_CREAT|syscall.O_RDONLY,0600)
	
}