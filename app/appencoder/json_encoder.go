package appencoder

import jsoniter "github.com/json-iterator/go"

var JsonEncoder AppJsonEncoder

type AppJsonEncoder interface {
	Marshal(obj interface{}) ([]byte, error)
	Unmarshal(data []byte, obj interface{}) error
}

func SetAppJsonEncoder(jsonEncoder AppJsonEncoder) {
	JsonEncoder = jsonEncoder
}

type AppJsonEncoderImpl struct {
	encoder jsoniter.API
}

func NewAppJsonEncoderImpl() *AppJsonEncoderImpl {
	return &AppJsonEncoderImpl{
		encoder: jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

func (a *AppJsonEncoderImpl) Marshal(obj interface{}) ([]byte, error) {
	return a.encoder.Marshal(obj)
}

func (a *AppJsonEncoderImpl) Unmarshal(data []byte, obj interface{}) error {
	return a.encoder.Unmarshal(data, obj)
}
