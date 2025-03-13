package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/Sofja96/GophKeeper.git/internal/client/encryption"
	"github.com/Sofja96/GophKeeper.git/pkg"
	"github.com/Sofja96/GophKeeper.git/proto"
)

// DataType интерфейс для всех моделей данных
type DataType interface {
	Validate() error
	ToJSON() ([]byte, error)
}

// TextDataType - структура для текстовых данных.
type TextDataType struct {
	Text string `json:"text"`
}

// Validate проверяет, что текст не пуст.
func (t *TextDataType) Validate() error {
	if len(t.Text) == 0 {
		return errors.New("поле не может быть пустым")
	}
	return nil
}

// ToJSON текст в JSON
func (t *TextDataType) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

// BinaryDataType - структура бинарных данных
type BinaryDataType struct {
	FilePath string `json:"file_path"`
	Content  []byte `json:"content"`
	Filename string `json:"filename"`
}

// Validate проверяет, что путь к файлу указан и объект данных инициализирован.
func (b *BinaryDataType) Validate() error {
	if b == nil {
		return errors.New("объект данных не инициализирован, проверьте путь к файлу")
	}

	fmt.Println("file_path", b.FilePath)
	if len(b.FilePath) == 0 {
		return errors.New("поле не может быть пустым")
	}
	return nil
}

// ToJSON кодирует бинарные данные в base64 и возвращает их в формате JSON.
func (b *BinaryDataType) ToJSON() []byte {
	return []byte(encryption.EncodeData(b.Content))
}

// BankCardType - структура банковской карты
type BankCardType struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
	HolderName string `json:"holder_name"`
}

// IsValid проверяет карту по алгоритму Луна
func (bc *BankCardType) IsValid() bool {
	return pkg.ValidateLuhn(bc.CardNumber)
}

// Validate проверяет, что карта корректная
func (bc *BankCardType) Validate() error {
	bc.CardNumber = strings.ReplaceAll(bc.CardNumber, " ", "")
	if !bc.IsValid() {
		return errors.New("номер карты не проходит проверку по Луну")
	}
	if len(bc.CVV) != 3 {
		return fmt.Errorf("cvv must be 3 digits")
	}
	if bc.CardNumber == "" || bc.ExpiryDate == "" || bc.CVV == "" || bc.HolderName == "" {
		return errors.New("не все поля заполнены")
	}
	return nil
}

// ToJSON сериализует карту в JSON
func (bc *BankCardType) ToJSON() ([]byte, error) {
	return json.Marshal(bc)
}

// LoginPasswordType - структура логина и пароля
type LoginPasswordType struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate проверяет логин/пароль
func (lp *LoginPasswordType) Validate() error {
	if lp.Username == "" || lp.Password == "" {
		return errors.New("логин и пароль не могут быть пустыми")
	}
	return nil
}

// ToJSON сериализует логин/пароль в JSON
func (lp *LoginPasswordType) ToJSON() ([]byte, error) {
	return json.Marshal(lp)
}

// CreateData - структура для создания данных.
type CreateData struct {
	Data          DataType
	DataType      proto.DataType
	EncryptionKey []byte
	Metadata      *structpb.Struct
	Filename      string
}
