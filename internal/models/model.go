package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/Sofja96/GophKeeper.git/proto"
)

// ContextKey - тип для ключа контекста
type ContextKey string

const ContextKeyUser ContextKey = "username"

// User - структура для хранения данных пользователя
// Содержит логин и пароль.
type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

// DataType - тип для представления разных типов данных
type DataType string

// String возвращает строковое представление DataType
func (d DataType) String() string {
	return string(d)
}

const (
	LoginPassword DataType = "LOGIN_PASSWORD" // Логин и пароль
	TextData      DataType = "TEXT_DATA"      // Текстовые данные
	BinaryData    DataType = "BINARY_DATA"    // Бинарные данные
	BankCard      DataType = "BANK_CARD"      // Данные банковской карты
)

// Data Общая структура данных
type Data struct {
	ID          int64     `json:"id,omitempty" db:"id"`
	UserID      int64     `json:"user_id,omitempty" db:"user_id"`
	DataType    DataType  `json:"data_type" db:"data_type"`
	FileName    string    `json:"file_name,omitempty"`
	DataContent []byte    `json:"data_content" db:"data_content"`
	Metadata    JSONB     `json:"metadata,omitempty" db:"metadata"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SetMetadata обновляет метаданные для модели.
func (d *Data) SetMetadata(key string, value interface{}) {
	if d.Metadata == nil {
		d.Metadata = JSONB{}
	}
	d.Metadata[key] = value
}

// GetMetadata возвращает значение по ключу из метаданных.
// Если метаданные не инициализированы или ключ отсутствует, возвращает nil и false.
func (d *Data) GetMetadata(key string) (interface{}, bool) {
	if d.Metadata == nil {
		return nil, false
	}
	value, ok := d.Metadata[key]
	return value, ok
}

// ProtoToModelMapping Мапа для сопоставления proto.DataType и DataType
var ProtoToModelMapping = map[string]string{
	proto.DataType_LOGIN_PASSWORD.String(): LoginPassword.String(),
	proto.DataType_TEXT_DATA.String():      TextData.String(),
	proto.DataType_BINARY_DATA.String():    BinaryData.String(),
	proto.DataType_BANK_CARD.String():      BankCard.String(),
}

// ProtoToModelMappingToProto  Мапа для сопоставления proto.DataType и DataType
var ProtoToModelMappingToProto = map[string]string{
	LoginPassword.String(): proto.DataType_LOGIN_PASSWORD.String(),
	TextData.String():      proto.DataType_TEXT_DATA.String(),
	BinaryData.String():    proto.DataType_BINARY_DATA.String(),
	BankCard.String():      proto.DataType_BANK_CARD.String(),
}

// ConvertModelDataTypeToProto преобразует DataType в соответствующий proto.DataType.
func ConvertModelDataTypeToProto(dataType DataType) (proto.DataType, error) {
	for modelType, protoType := range ProtoToModelMappingToProto {
		if dataType.String() == modelType {
			return proto.DataType(proto.DataType_value[protoType]), nil
		}
	}
	return proto.DataType(0), fmt.Errorf("unknown DataType: %v", dataType)
}

// GetModelType Получение модели по proto.DataType
func GetModelType(protoType proto.DataType) (DataType, error) {
	if modelType, ok := ProtoToModelMapping[protoType.String()]; ok {
		return DataType(modelType), nil
	}
	return "", errors.New("unknown proto.DataType")
}

// ValidateType — проверяет допустимость типа данных.
func (d *Data) ValidateType() error {
	switch d.DataType {
	case LoginPassword, TextData, BinaryData, BankCard:
		return nil
	default:
		return errors.New("unsupported data type")
	}
}

// JSONB тип для метаданных
type JSONB map[string]interface{}

// Value сериализует JSONB в значение, подходящее для записи в базу данных.
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan десериализует значение из базы данных в JSONB.
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// ConvertJSONBToStruct преобразует JSONB в *structpb.Struct для использования с protobuf.
func ConvertJSONBToStruct(metadata JSONB) (*structpb.Struct, error) {
	if metadata == nil {
		return nil, nil
	}
	structProto, err := structpb.NewStruct(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to convert JSONB to structpb.Struct: %w", err)
	}
	return structProto, nil
}
