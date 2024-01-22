package main

import (
	"filesigner/base32encoding"
	"filesigner/filehasher"
	"filesigner/filesignature"
	"filesigner/hashsignature"
	"filesigner/keyid"
	"filesigner/logger"
	"filesigner/signaturefile"
	"filesigner/signaturehandler"
	"filesigner/texthelper"
	"fmt"
	"golang.org/x/exp/maps"
	"os"
	"path/filepath"
)

// ******** Private functions ********

// doVerification verifies a signature file.
func doVerification(contextId string) int {
	signatureData, err := signaturefile.ReadSignatureFile(signatureFileName)
	if err != nil {
		logger.PrintError(51, err.Error())
		return rcProcessError
	}

	var hashVerifier hashsignature.HashVerifier
	var publicKeyId string
	hashVerifier, publicKeyId, err = getHashVerifier(signatureData)
	if err != nil {
		logger.PrintError(52, err.Error())
		return rcProcessError
	}

	var ok bool
	ok, err = signatureData.Verify(hashVerifier, contextId)
	if err == nil {
		if !ok {
			logger.PrintError(53, "Signature file has been tampered with or wrong context id")
			return rcProcessError
		}
	} else {
		logger.PrintErrorf(54, "Error verifying signature file data signature: %v", err)
		return rcProcessError
	}

	logger.PrintInfof(55, "Context id         : %s", contextId)
	logger.PrintInfof(56, "Public key id      : %s", publicKeyId)
	logger.PrintInfof(57, "Signature timestamp: %s", signatureData.Timestamp)
	logger.PrintInfof(58, "Signature host name: %s", signatureData.Hostname)

	successCount, errorCount, rc := verifyFiles(contextId, signatureData, hashVerifier)

	successEnding := texthelper.GetCountEnding(successCount)
	errorEnding := texthelper.GetCountEnding(errorCount)

	switch rc {
	case rcOK:
		logger.PrintInfof(59, "Verification of %d file%s successful", successCount, successEnding)

	case rcProcessWarning:
		logger.PrintWarningf(60, "Verification of %d file%s successful and warnings present", successCount, successEnding)

	case rcProcessError:
		logger.PrintErrorf(61, "Verification of %d file%s successful and %d file%s unsuccessful", successCount, successEnding, errorCount, errorEnding)
	}

	return rc
}

// verifyFiles verifies the signatures of the files in the signature data.
func verifyFiles(contextId string,
	signatureData *signaturehandler.SignatureData,
	hashVerifier hashsignature.HashVerifier) (int, int, int) {
	filePaths, rc := getExistingFiles(maps.Keys(signatureData.FileSignatures))

	if len(filePaths) == 0 {
		logger.PrintWarning(62, "No files from signature file present")
		return 0, 0, rcProcessWarning
	}

	hashList := filehasher.FileHashes(filePaths, contextId)
	if existHashErrors(hashList) {
		return 0, 0, rcProcessError
	}

	successList, errorList := filesignature.VerifyFileHashes(hashVerifier, signatureData.FileSignatures, hashList)

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

// getHashVerifier constructs the hash verifier and the key id from the signature data.
func getHashVerifier(signatureData *signaturehandler.SignatureData) (hashsignature.HashVerifier, string, error) {
	publicKeyBytes, err := base32encoding.DecodeFromString(signatureData.PublicKey)
	if err != nil {
		return nil, "", fmt.Errorf("Could not convert public key to bytes: %w", err)
	}

	var hashVerifier hashsignature.HashVerifier
	if signatureData.SignatureType == signaturehandler.SignatureTypeEd25519 {
		hashVerifier, err = hashsignature.NewEd25519HashVerifier(publicKeyBytes)
	} else {
		hashVerifier, err = hashsignature.NewEc521HashVerifier(publicKeyBytes)
	}
	if err != nil {
		return nil, "", fmt.Errorf("Could not create hash verifier: %w", err)
	}

	return hashVerifier, keyid.KeyId(publicKeyBytes), nil
}

// getExistingFiles gets the files from a signature list that exist in the directory that is to be verified.
func getExistingFiles(filePaths []string) ([]string, int) {
	rc := rcOK

	result := make([]string, 0, len(filePaths))
	for _, fp := range filePaths {
		nfp := filepath.FromSlash(fp)
		fi, err := os.Stat(nfp)
		if err != nil {
			logger.PrintWarningf(63, "'%s' from signature file does not exist", nfp)
			rc = rcProcessWarning
		} else {
			if fi.IsDir() {
				logger.PrintWarningf(64, "'%s' from signature file is a directory", nfp)
				rc = rcProcessWarning
			} else {
				result = append(result, nfp)
			}
		}
	}

	return result, rc
}
