package signaturefile

import (
	"encoding/json"
	"filesigner/signaturehandler"
	"fmt"
	"os"
)

// ******** Public functions ********

func WriteSignatureFile(filePath string, signatureData *signaturehandler.SignatureData) error {
	jsonOutput, err := json.MarshalIndent(signatureData, "", "   ")
	if err != nil {
		return fmt.Errorf("Could not create json format for signature file: %w", err)
	}

	err = os.WriteFile(filePath, jsonOutput, 0600)
	if err != nil {
		return fmt.Errorf("Could not write signature file: %w", err)
	}

	return nil
}

// ReadSignatureFile reads a signature file and returns the signature data.
func ReadSignatureFile(filePath string) (*signaturehandler.SignatureData, error) {
	result, err := getSignatureData(filePath)
	if err != nil {
		return nil, err
	}

	err = checkSignatureForm(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ******** Private functions ********

// getSignatureData reads the signature data from a file and returns the data in a SignatureData structure.
func getSignatureData(filePath string) (*signaturehandler.SignatureData, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	result := &signaturehandler.SignatureData{
		Format:        signaturehandler.SignatureFormatInvalid,
		SignatureType: signaturehandler.SignatureTypeInvalid,
	}
	err = json.Unmarshal(fileContent, result)

	return result, err
}

// checkSignatureForm checks if the signature data are complete and satisfy some formal requirements.
func checkSignatureForm(signatureData *signaturehandler.SignatureData) error {
	err := checkMissingInformation(signatureData)
	if err != nil {
		return err
	}

	if signatureData.Format > signaturehandler.SignatureFormatMax {
		return fmt.Errorf("Invalid signature format id: %d", signatureData.Format)
	}

	if signatureData.SignatureType > signaturehandler.SignatureTypeMax {
		return fmt.Errorf("Invalid signature type: %d", signatureData.SignatureType)
	}

	return nil
}

// checkMissingInformation checks if any required signature result data is missing.
func checkMissingInformation(signatureData *signaturehandler.SignatureData) error {
	if len(signatureData.DataSignature) == 0 {
		return makeMissingFieldError("DataSignature")
	}
	if signatureData.FileSignatures == nil {
		return makeMissingFieldError("FileSignatures")
	}
	if signatureData.Format == signaturehandler.SignatureFormatInvalid {
		return makeMissingFieldError("Format")
	}
	if len(signatureData.Hostname) == 0 {
		return makeMissingFieldError("Hostname")
	}
	if len(signatureData.PublicKey) == 0 {
		return makeMissingFieldError("PublicKey")
	}
	if signatureData.SignatureType == signaturehandler.SignatureTypeInvalid {
		return makeMissingFieldError("SignatureType")
	}
	if len(signatureData.Timestamp) == 0 {
		return makeMissingFieldError("Timestamp")
	}
	return nil
}

// makeMissingFieldError build the error for a "missing field" error type.
func makeMissingFieldError(fieldName string) error {
	return fmt.Errorf("Field '%s' is missing from signature file", fieldName)
}
