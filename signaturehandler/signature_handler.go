package signaturehandler

import (
	"filesigner/base32encoding"
	"filesigner/contexthasher"
	"filesigner/hashsigner"
	"filesigner/maphelper"
	"filesigner/stringhelper"
	"golang.org/x/crypto/sha3"
	"hash"
)

// ******** Public types ********

type SignatureFormat byte
type SignatureType byte

type SignatureData struct {
	Format         SignatureFormat
	PublicKey      string
	Timestamp      string
	Hostname       string
	SignatureType  SignatureType
	FileSignatures map[string]string
	DataSignature  string
}

// ******** Public constants ********

// SignatureFormatMax is the maximum (i.e. newest) format id of the signature result format.
const (
	SignatureFormatInvalid SignatureFormat = iota
	SignatureFormatV1
	SignatureFormatMax = iota - 1
)

const (
	SignatureTypeInvalid SignatureType = iota
	SignatureTypeEd25519
	SignatureTypeEcDsaP521
	SignatureTypeMax = iota - 1
)

// ******** Public type functions ********

// Sign adds the data signature to a SignatureData
func (sd *SignatureData) Sign(hashSigner hashsigner.HashSigner, contextId string) error {
	hashValue := getHashValueOfSignatureData(sd, contextId)
	signatureValue, err := hashSigner.SignHash(hashValue)
	if err != nil {
		return err
	}

	sd.DataSignature = base32encoding.EncodeToString(signatureValue)

	return nil
}

// Verify verifies the data signature of a SignatureData
func (sd *SignatureData) Verify(hashVerifier hashsigner.HashVerifier, contextId string) (bool, error) {
	dataSignature, err := base32encoding.DecodeFromString(sd.DataSignature)
	if err != nil {
		return false, err
	}

	hashValue := getHashValueOfSignatureData(sd, contextId)
	return hashVerifier.VerifyHash(hashValue, dataSignature)
}

// ******** Private functions ********

// getHashValueOfSignatureData calculates the hash value of a SignatureData
func getHashValueOfSignatureData(signatureData *SignatureData, contextId string) []byte {
	hasher := contexthasher.NewContextHasher(sha3.New512(), contextId)

	position := make([]byte, 1)
	tempSlice := make([]byte, 1)

	hashPosition(hasher, position)
	tempSlice[0] = byte(signatureData.Format)
	hasher.Write(tempSlice)

	hashPosition(hasher, position)
	hasher.Write([]byte(signatureData.PublicKey))

	hashPosition(hasher, position)
	hasher.Write([]byte(signatureData.Timestamp))

	hashPosition(hasher, position)
	hasher.Write([]byte(signatureData.Hostname))

	hashPosition(hasher, position)
	tempSlice[0] = byte(signatureData.SignatureType)
	hasher.Write(tempSlice)

	sortedFileNames := maphelper.SortedKeys(signatureData.FileSignatures)
	for _, fileName := range sortedFileNames {
		hashPosition(hasher, position)
		hasher.Write(stringhelper.UnsafeStringBytes(fileName))

		hashPosition(hasher, position)
		hasher.Write(stringhelper.UnsafeStringBytes(signatureData.FileSignatures[fileName]))
	}

	return hasher.Sum(nil)
}

// hashPosition writes the position into the hasher
func hashPosition(hasher hash.Hash, position []byte) {
	position[0]++
	hasher.Write(position)
}