package utils

import (
	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

type QRUtils struct{}

func (*QRUtils) GenerateQrCodeData(email string) ([]byte, string, error) {
	secret := gotp.RandomSecret(16)
	totp := gotp.NewDefaultTOTP(secret)
	uri := totp.ProvisioningUri(email, "CTY AE TNH")
	data, err := qrcode.Encode(uri, qrcode.Medium, 256)

	if err != nil {
		return nil, "", err
	}
	return data, secret, nil
}
