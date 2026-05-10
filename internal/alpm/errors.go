package alpm

// #cgo pkg-config: libalpm
// #include <alpm.h>
import "C"

const (
	ErrOK                          Errno = C.ALPM_ERR_OK
	ErrMemory                      Errno = C.ALPM_ERR_MEMORY
	ErrSystem                      Errno = C.ALPM_ERR_SYSTEM
	ErrBadPerms                    Errno = C.ALPM_ERR_BADPERMS
	ErrNotAFile                    Errno = C.ALPM_ERR_NOT_A_FILE
	ErrNotADir                     Errno = C.ALPM_ERR_NOT_A_DIR
	ErrWrongArgs                   Errno = C.ALPM_ERR_WRONG_ARGS
	ErrDiskSpace                   Errno = C.ALPM_ERR_DISK_SPACE
	ErrHandleNull                  Errno = C.ALPM_ERR_HANDLE_NULL
	ErrHandleNotNull               Errno = C.ALPM_ERR_HANDLE_NOT_NULL
	ErrHandleLock                  Errno = C.ALPM_ERR_HANDLE_LOCK
	ErrDBOpen                      Errno = C.ALPM_ERR_DB_OPEN
	ErrDBCreate                    Errno = C.ALPM_ERR_DB_CREATE
	ErrDBNull                      Errno = C.ALPM_ERR_DB_NULL
	ErrDBNotNull                   Errno = C.ALPM_ERR_DB_NOT_NULL
	ErrDBNotFound                  Errno = C.ALPM_ERR_DB_NOT_FOUND
	ErrDBInvalid                   Errno = C.ALPM_ERR_DB_INVALID
	ErrDBInvalidSig                Errno = C.ALPM_ERR_DB_INVALID_SIG
	ErrDBVersion                   Errno = C.ALPM_ERR_DB_VERSION
	ErrDBWrite                     Errno = C.ALPM_ERR_DB_WRITE
	ErrDBRemove                    Errno = C.ALPM_ERR_DB_REMOVE
	ErrServerBadUrl                Errno = C.ALPM_ERR_SERVER_BAD_URL
	ErrServerNone                  Errno = C.ALPM_ERR_SERVER_NONE
	ErrTransNotNull                Errno = C.ALPM_ERR_TRANS_NOT_NULL
	ErrTransNull                   Errno = C.ALPM_ERR_TRANS_NULL
	ErrTransDupTarget              Errno = C.ALPM_ERR_TRANS_DUP_TARGET
	ErrTransDupFilename            Errno = C.ALPM_ERR_TRANS_DUP_FILENAME
	ErrTransNotInitialized         Errno = C.ALPM_ERR_TRANS_NOT_INITIALIZED
	ErrTransNotPrepared            Errno = C.ALPM_ERR_TRANS_NOT_PREPARED
	ErrTransAbort                  Errno = C.ALPM_ERR_TRANS_ABORT
	ErrTransType                   Errno = C.ALPM_ERR_TRANS_TYPE
	ErrTransNotLocked              Errno = C.ALPM_ERR_TRANS_NOT_LOCKED
	ErrTransHookFailed             Errno = C.ALPM_ERR_TRANS_HOOK_FAILED
	ErrPkgNotFound                 Errno = C.ALPM_ERR_PKG_NOT_FOUND
	ErrPkgIgnored                  Errno = C.ALPM_ERR_PKG_IGNORED
	ErrPkgInvalid                  Errno = C.ALPM_ERR_PKG_INVALID
	ErrPkgInvalidChecksum          Errno = C.ALPM_ERR_PKG_INVALID_CHECKSUM
	ErrPkgInvalidSig               Errno = C.ALPM_ERR_PKG_INVALID_SIG
	ErrPkgMissingSig               Errno = C.ALPM_ERR_PKG_MISSING_SIG
	ErrPkgOpen                     Errno = C.ALPM_ERR_PKG_OPEN
	ErrPkgCantRemove               Errno = C.ALPM_ERR_PKG_CANT_REMOVE
	ErrPkgInvalidName              Errno = C.ALPM_ERR_PKG_INVALID_NAME
	ErrPkgInvalidArch              Errno = C.ALPM_ERR_PKG_INVALID_ARCH
	ErrSigMissing                  Errno = C.ALPM_ERR_SIG_MISSING
	ErrSigInvalid                  Errno = C.ALPM_ERR_SIG_INVALID
	ErrUnsatisfiedDeps             Errno = C.ALPM_ERR_UNSATISFIED_DEPS
	ErrConflictingDeps             Errno = C.ALPM_ERR_CONFLICTING_DEPS
	ErrFileConflicts               Errno = C.ALPM_ERR_FILE_CONFLICTS
	ErrRetrievePrepare             Errno = C.ALPM_ERR_RETRIEVE_PREPARE
	ErrRetrieve                    Errno = C.ALPM_ERR_RETRIEVE
	ErrInvalidRegex                Errno = C.ALPM_ERR_INVALID_REGEX
	ErrLibArchive                  Errno = C.ALPM_ERR_LIBARCHIVE
	ErrLibCURL                     Errno = C.ALPM_ERR_LIBCURL
	ErrExternalDownload            Errno = C.ALPM_ERR_EXTERNAL_DOWNLOAD
	ErrGPGME                       Errno = C.ALPM_ERR_GPGME
	ErrMissingCapabilitySignatures Errno = C.ALPM_ERR_MISSING_CAPABILITY_SIGNATURES
)

type Errno int

func (err Errno) Error() string {
	return C.GoString(C.alpm_strerror(C.alpm_errno_t(err)))
}
