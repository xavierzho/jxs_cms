package util

type TelConfig struct {
	AreaCode           string
	FullPhoneNumLength int
}

func (t TelConfig) GetFullPhoneNum(phoneNum string) (bool, string) {
	if len(t.AreaCode+phoneNum) == t.FullPhoneNumLength {
		return true, t.AreaCode + phoneNum
	}
	if len(phoneNum) == t.FullPhoneNumLength && phoneNum[:2] == t.AreaCode {
		return true, phoneNum
	}
	return false, ""
}
