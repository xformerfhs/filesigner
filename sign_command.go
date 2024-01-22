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
	"os"
	"time"
)

// ******** Private constants ********

// timeStampFormat Format for signature file time stamp
const timeStampFormat = "2006-01-02 15:04:05 Z07:00"

// ******** Private functions ********

// doSigning signs all files with the given context id.
func doSigning(signatureType signaturehandler.SignatureType, contextId string, filePaths []string) int {
	var err error

	signatureData := &signaturehandler.SignatureData{
		Format:        signaturehandler.SignatureFormatV1,
		Timestamp:     time.Now().Format(timeStampFormat),
		SignatureType: signatureType,
	}

	signatureData.Hostname, err = os.Hostname()
	if err != nil {
		logger.PrintErrorf(31, "Could not get host name: %v", err)
		return rcProcessError
	}

	resultList := filehasher.FileHashes(filePaths, contextId)

	if existHashErrors(resultList) {
		return rcProcessError
	}

	var hashSigner hashsignature.HashSigner
	if signatureType == signaturehandler.SignatureTypeEd25519 {
		hashSigner, err = hashsignature.NewEd25519HashSigner()
	} else {
		hashSigner, err = hashsignature.NewEcDsa521HashSigner()
	}
	if err != nil {
		logger.PrintErrorf(33, "Could not create hash-signer: %v", err)
		return rcProcessError
	}
	defer hashSigner.Destroy()

	var publicKeyBytes []byte
	publicKeyBytes, err = hashSigner.GetPublicKey()
	if err != nil {
		logger.PrintErrorf(34, "Could not get public key bytes: %v", err)
		return rcProcessError
	}
	signatureData.PublicKey = base32encoding.EncodeToString(publicKeyBytes)

	var successList []string
	signatureData.FileSignatures, successList, err = filesignature.SignFileHashes(hashSigner, resultList)
	if err != nil {
		logger.PrintErrorf(35, "Could not sign file hashes: %v", err)
		return rcProcessError
	}

	err = signatureData.Sign(hashSigner, contextId)
	if err != nil {
		logger.PrintErrorf(36, "Could not sign signature file data: %v", err)
		return rcProcessError
	}

	err = signaturefile.WriteSignatureFile(signatureFileName, signatureData)
	if err != nil {
		logger.PrintError(37, err.Error())
		return rcProcessError
	}

	logger.PrintInfof(38, "Context id         : %s", contextId)
	logger.PrintInfof(39, "Public key id      : %s", keyid.KeyId(publicKeyBytes))
	logger.PrintInfof(40, "Signature timestamp: %s", signatureData.Timestamp)
	logger.PrintInfof(41, "Signature host name: %s", signatureData.Hostname)

	successCount := len(successList)
	if successCount > 0 {
		printSuccessList("Signing", successList)
	}

	successEnding := texthelper.GetCountEnding(successCount)

	logger.PrintInfof(42, "Signature%s for %d file%s successfully created", successEnding, len(successList), successEnding)
	return rcOK
}
