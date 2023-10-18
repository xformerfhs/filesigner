package main

import (
	"encoding/json"
	"filesigner/base32encoding"
	"filesigner/filehashing"
	"filesigner/hashsigner"
	"filesigner/logger"
	"filesigner/signaturehandler"
	"filesigner/texthelper"
	"golang.org/x/exp/maps"
	"os"
	"path/filepath"
)

// ******** Private functions ********

// doVerification verifies a signature file.
func doVerification(contextId string) int {
	signatureData, err := getSignatureData(signatureFileName)
	if err != nil {
		logger.PrintErrorf(51, "Error reading signature file: %v", err)
		return rcProcessError
	}

	if signatureData.Format != signaturehandler.MaxFormatId {
		logger.PrintErrorf(52, "Invalid signature format id: %d", signatureData.Format)
		return rcProcessError
	}
	if signatureData.SignatureType != signaturehandler.SignatureTypeEd25519 &&
		signatureData.SignatureType != signaturehandler.SignatureTypeEcDsaP521 {
		logger.PrintErrorf(53, "Invalid signature type: %d", signatureData.SignatureType)
		return rcProcessError
	}

	filePaths, rc := getExistingFiles(maps.Keys(signatureData.FileSignatures))

	if len(filePaths) == 0 {
		logger.PrintWarning(54, "No files from signature file present")
		return rcProcessWarning
	}

	var publicKeyBytes []byte
	publicKeyBytes, err = base32encoding.DecodeFromString(signatureData.PublicKey)
	if err != nil {
		logger.PrintErrorf(55, "Could not convert public key to bytes: %v", err)
		return rcProcessError
	}

	hashList := filehashing.GetFileHashes(filePaths, contextId)
	if existHashErrors(hashList) {
		return rcProcessError
	}

	var hashVerifier hashsigner.HashVerifier
	if signatureData.SignatureType == signaturehandler.SignatureTypeEd25519 {
		hashVerifier, err = hashsigner.NewEd25519HashVerifier(publicKeyBytes)
	} else {
		hashVerifier, err = hashsigner.NewEc521HashVerifier(publicKeyBytes)
	}
	if err != nil {
		logger.PrintErrorf(56, "Could not create hash-verifier: %v", err)
		return rcProcessError
	}

	var publicKeyHash []byte
	publicKeyHash = getKeyHash(publicKeyBytes)
	if err != nil {
		logger.PrintErrorf(57, "Could not get public key hash: %v", err)
		return rcProcessError
	}

	logger.PrintInfof(58, "Context id         : %s", contextId)
	logger.PrintInfof(59, "Public key id      : %s", base32encoding.EncodeKey(publicKeyHash))
	logger.PrintInfof(60, "Signature timestamp: %s", signatureData.Timestamp)
	logger.PrintInfof(61, "Signature host name: %s", signatureData.Hostname)

	var ok bool
	ok, err = signatureData.Verify(hashVerifier, contextId)
	if err == nil {
		if !ok {
			logger.PrintError(62, "Signature file has been tampered with or wrong context id")
			return rcProcessError
		}
	} else {
		logger.PrintErrorf(63, "Error verifying signature file data signature: %v", err)
		return rcProcessError
	}

	successList, errorList := signaturehandler.VerifyFileHashes(hashVerifier, signatureData, hashList)

	successCount := len(successList)
	if successCount > 0 {
		printSuccessList("Verification", successList)
	}

	errorCount := len(errorList)
	if errorCount > 0 {
		printErrorList(errorList)
		rc = rcProcessError
	}

	successEnding := texthelper.GetCountEnding(successCount)
	errorEnding := texthelper.GetCountEnding(errorCount)

	switch rc {
	case rcOK:
		logger.PrintInfof(64, "Verification of %d file%s successful", successCount, successEnding)

	case rcProcessWarning:
		logger.PrintWarningf(65, "Verification of %d file%s successful and warnings present", successCount, successEnding)

	case rcProcessError:
		logger.PrintErrorf(66, "Verification of %d file%s successful and %d file%s unsuccessful", successCount, successEnding, errorCount, errorEnding)
	}

	return rc
}

// getSignatureData reads the signature data from a file and returns the data in a SignatureResult structure.
func getSignatureData(filePath string) (*signaturehandler.SignatureResult, error) {
	result := &signaturehandler.SignatureResult{}
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
			logger.PrintWarningf(67, "File '%s' from signature file does not exist", fp)
			rc = rcProcessWarning
		} else {
			if fi.IsDir() {
				logger.PrintWarningf(68, "File '%s' from signature file is a directory", fp)
				rc = rcProcessWarning
			} else {
				result = append(result, filepath.FromSlash(fp))
			}
		}
	}

	return result, rc
}
