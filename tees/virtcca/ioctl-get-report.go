package virtcca

import (
	"fmt"
	"github.com/vtolstov/go-ioctl"
	"golang.org/x/sys/unix"
	"os"
	"unsafe"
)

// Size in bytes of the SHA512 measurement
const SHA512_SIZE = 64

// Size in bytes of the SHA256 measurement
const SHA256_SIZE = 32

/*
 * Size in bytes of the largest measurement type that can be supported.
 * This constant needs to be updated accordingly if new algorithms are supported.
 */
const MAX_MEASUREMENT_SIZE = SHA512_SIZE
const MAX_DEV_CERT_SIZE = 4096

const GRANULE_SIZE = 4096
const MAX_TOKEN_GRANULE_COUNT = 2
const CHALLENGE_SIZE = 64

type CVMMeasurement struct {
	Index int32
	Value [MAX_MEASUREMENT_SIZE]byte
}

type CVMTSIVersion struct {
	Major int32
	Minor int32
}

type CVMAttestationCmd struct {
	Challenge [CHALLENGE_SIZE]byte // input: challenge value
	Token     [GRANULE_SIZE * MAX_TOKEN_GRANULE_COUNT]byte
	TokenSize uint64 // return: token size
}

type CCADevCert struct {
	Size  uint64
	Value [MAX_DEV_CERT_SIZE]byte
}

const TSI_DEV = "/dev/tsi"
const TSI_MAGIC = 'T'

// Ioctl commands
var TMM_GET_TSI_VERSION = ioctl.IOR(TSI_MAGIC, 0, 8)
var TMM_GET_ATTESTATION_TOKEN = ioctl.IOWR(TSI_MAGIC, 1, 8264)
var TMM_GET_DEV_CERT = ioctl.IOR(TSI_MAGIC, 2, 4104)

func GetVersion() (*CVMTSIVersion, error) {
	fd, err := os.OpenFile(TSI_DEV, os.O_RDWR|os.O_EXCL, 0)
	if err != nil {
		return nil, fmt.Errorf("open /dev/tsi failed: %s", err)
	}
	defer fd.Close()

	var version CVMTSIVersion
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, fd.Fd(), TMM_GET_TSI_VERSION, uintptr(unsafe.Pointer(&version)))
	if errno != 0 {
		return nil, fmt.Errorf("failed to get TSI version: %v", errno)
	}
	return &version, nil
}

func GetAttestationToken(challenge []byte) ([]byte, error) {
	if len(challenge) > CHALLENGE_SIZE {
		return nil, fmt.Errorf("challenge too long")
	}

	var cmd CVMAttestationCmd
	copy(cmd.Challenge[:], challenge)

	fd, err := os.OpenFile(TSI_DEV, os.O_RDWR|os.O_EXCL, 0)
	if err != nil {
		return nil, fmt.Errorf("open /dev/tsi failed: %s", err)
	}
	defer fd.Close()

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, fd.Fd(), TMM_GET_ATTESTATION_TOKEN, uintptr(unsafe.Pointer(&cmd)))
	if errno != 0 {
		return nil, fmt.Errorf("failed to get attestation token: %v", errno)
	}

	return cmd.Token[:cmd.TokenSize], nil
}

func GetDevCert() ([]byte, error) {
	fd, err := os.OpenFile(TSI_DEV, os.O_RDWR|os.O_EXCL, 0)
	if err != nil {
		return nil, fmt.Errorf("open /dev/tsi failed: %s", err)
	}
	defer fd.Close()

	var cert CCADevCert
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, fd.Fd(), TMM_GET_DEV_CERT, uintptr(unsafe.Pointer(&cert)))
	if errno != 0 {
		return nil, fmt.Errorf("failed to get dev cert: %v", errno)
	}

	return cert.Value[:cert.Size], nil
}
