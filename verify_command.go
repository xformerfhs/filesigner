package main

import (
	"encoding/json"
	"filesigner/base32encoding"
	"filesigner/filehashing"
	"filesigner/hashsigner"
	"filesigner/logger"
	"filesigner/signaturehandler"
	"filesigner/texthelper"
	"fmt"
	"golang.org/x/exp/maps"
	"os"
	"path/filepath"
)

// ******** Private constants ********

const invalidType byte = 255

// ******** Private functions ********

// doVerification verifies a signature file.
func doVerification(contextId string) int {
	signatureData, err := getSignatureData(signatureFileName)
	if err != nil {
		logger.PrintErrorf(51, "Error reading signature file: %v", err)
		return rcProcessError
	}

	err = checkSignatureForm(signatureData)
	if err != nil {
		logger.PrintErrorf(52, "Malformed signature file: %v", err)
		return rcProcessError
	}

	var hashVerifier hashsigner.HashVerifier
	var publicKeyId string
	hashVerifier, publicKeyId, err = getHashVerifier(signatureData)
	if err != nil {
		logger.PrintError(53, err.Error())
		return rcProcessError
	}

	var ok bool
	ok, err = signatureData.Verify(hashVerifier, contextId)
	if err == nil {
		if !ok {
			logger.PrintError(54, "Signature file has been tampered with or wrong context id")
			return rcProcessError
		}
	} else {
		logger.PrintErrorf(55, "Error verifying signature file data signature: %v", err)
		return rcProcessError
	}

	logger.PrintInfof(56, "Context id         : %s", contextId)
	logger.PrintInfof(57, "Public key id      : %s", publicKeyId)
	logger.PrintInfof(58, "Signature timestamp: %s", signatureData.Timestamp)
	logger.PrintInfof(59, "Signature host name: %s", signatureData.Hostname)

	successCount, errorCount, rc := verifyFiles(contextId, signatureData, hashVerifier)

	successEnding := texthelper.GetCountEnding(successCount)
	errorEnding := texthelper.GetCountEnding(errorCount)

	switch rc {
	case rcOK:
		logger.PrintInfof(60, "Verification of %d file%s successful", successCount, successEnding)

	case rcProcessWarning:
		logger.PrintWarningf(61, "Verification of %d file%s successful and warnings present", successCount, successEnding)

	case rcProcessError:
		logger.PrintErrorf(62, "Verification of %d file%s successful and %d file%s unsuccessful", successCount, successEnding, errorCount, errorEnding)
	}

	return rc
}

func verifyFiles(contextId string,
	signatureData *signaturehandler.SignatureResult,
	hashVerifier hashsigner.HashVerifier) (int, int, int) {
	filePaths, rc := getExistingFiles(maps.Keys(signatureData.FileSignatures))

	if len(filePaths) == 0 {
		logger.PrintWarning(63, "No files from signature file present")
		return 0, 0, rcProcessWarning
	}

	hashList := filehashing.GetFileHashes(filePaths, contextId)
	if existHashErrors(hashList) {
		return 0, 0, rcProcessError
	}

	successList, errorList := signaturehandler.VerifyFileHashes(hashVerifier, signatureData.FileSignatures, hashList)

	successCount := len(successList)
	if successCount > 0 {
		printSuccessList("Verification", successList)
	}

	errorCount := len(errorList)
	if errorCount > 0 {
		printErrorList(errorList)
		rc = rcProcessError
	}

	return successCount, errorCount, rc
}

func getHashVerifier(signatureData *signaturehandler.SignatureResult) (hashsigner.HashVerifier, string, error) {
	publicKeyBytes, err := base32encoding.DecodeFromString(signatureData.PublicKey)
	if err != nil {
		return nil, "", fmt.Errorf("Could not convert public key to bytes: %w", err)
	}

	var hashVerifier hashsigner.HashVerifier
	if signatureData.SignatureType == signaturehandler.SignatureTypeEd25519 {
		hashVerifier, err = hashsigner.NewEd25519HashVerifier(publicKeyBytes)
	} else {
		hashVerifier, err = hashsigner.NewEc521HashVerifier(publicKeyBytes)
	}
	if err != nil {
		return nil, "", fmt.Errorf("Could not create hash verifier: %w", err)
	}

	return hashVerifier, getPublicKeyId(publicKeyBytes), nil
}

func getPublicKeyId(publicKeyBytes []byte) string {
	var publicKeyHash []byte
	publicKeyHash = getKeyHash(publicKeyBytes)

	return base32encoding.EncodeKey(publicKeyHash)
}

func checkSignatureForm(signatureData *signaturehandler.SignatureResult) error {
	err := checkMissingInformation(signatureData)
	if err != nil {
		return err
	}

	if signatureData.Format != signaturehandler.MaxFormatId {
		return fmt.Errorf("Invalid signature format id: %d", signatureData.Format)
	}

	if signatureData.SignatureType != signaturehandler.SignatureTypeEd25519 &&
		signatureData.SignatureType != signaturehandler.SignatureTypeEcDsaP521 {
		return fmt.Errorf("Invalid signature type: %d", signatureData.SignatureType)
	}

	return nil
}

// getSignatureData reads the signature data from a file and returns the data in a SignatureResult structure.
func getSignatureData(filePath string) (*signaturehandler.SignatureResult, error) {
	result := &signaturehandler.SignatureResult{Format: invalidType, SignatureType: invalidType}
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fileContent, result)

	return result, err
}

// getExistingFiles gets the files from a signature list that exist in the directory that is to be verified.
func getExistingFiles(filePaths []string) ([]string, int) {
	rc := rcOK

	result := make([]string, 0, len(filePaths))
	for _, fp := range filePaths {
		fi, err := os.Stat(fp)
		if err != nil {
			logger.PrintWarningf(64, "'%s' from signature file does not exist", fp)
			rc = rcProcessWarning
		} else {
			if fi.IsDir() {
				logger.PrintWarningf(65, "'%s' from signature file is a directory", fp)
				rc = rcProcessWarning
			} else {
				result = append(result, filepath.FromSlash(fp))
			}
		}
	}

	return result, rc
}

// checkMissingInformation checks if any required signature result data is missing.
func checkMissingInformation(signatureData *signaturehandler.SignatureResult) error {
	if len(signatureData.DataSignature) == 0 {
		return makeMissingFieldError("DataSignature")
	}
	if signatureData.FileSignatures == nil {
		return makeMissingFieldError("FileSignatures")
	}
	if signatureData.Format == invalidType {
		return makeMissingFieldError("Format")
	}
	if len(signatureData.Hostname) == 0 {
		return makeMissingFieldError("Hostname")
	}
	if len(signatureData.PublicKey) == 0 {
		return makeMissingFieldError("PublicKey")
	}
	if signatureData.SignatureType == invalidType {
		return makeMissingFieldError("SignatureType")
	}
	if len(signatureData.Timestamp) == 0 {
		return makeMissingFieldError("Timestamp")
	}
	return nil
}

func makeMissingFieldError(fieldName string) error {
	return fmt.Errorf("Field '%s' is missing from signature file", fieldName)
}
