package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Sofja96/GophKeeper.git/internal/client/encryption"
	"github.com/Sofja96/GophKeeper.git/internal/client/grpcclient"
	"github.com/Sofja96/GophKeeper.git/internal/client/localstorage"
	"github.com/Sofja96/GophKeeper.git/internal/models"
	mlogging "github.com/Sofja96/GophKeeper.git/internal/server/logger/mocks"
	"github.com/Sofja96/GophKeeper.git/proto"
	mproto "github.com/Sofja96/GophKeeper.git/proto/mocks"
)

func TestStartCLI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	client := &grpcclient.Client{
		Client: mockClient,
	}

	input := "8\n"

	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		defer w.Close()
		_, err := w.WriteString(input)
		assert.NoError(t, err)
	}()

	oldStdout := os.Stdout
	rout, wout, _ := os.Pipe()
	os.Stdout = wout

	go func() {
		err := StartCLI(client)
		assert.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	wout.Close()
	os.Stdin = oldStdin
	os.Stdout = oldStdout

	var output bytes.Buffer
	_, err := io.Copy(&output, rout)
	assert.NoError(t, err)

	assert.Contains(t, output.String(), "Выход из программы.")
}

func TestVersionCmd(t *testing.T) {
	var buf bytes.Buffer
	cmd := VersionCmd()
	cmd.SetOut(&buf)

	cmd.Run(cmd, []string{})

	assert.Contains(t, buf.String(), "Version:")
	assert.Contains(t, buf.String(), "Build Date:")
}

func TestRegisterCmd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	mockClient.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(&proto.RegisterResponse{}, nil).
		Times(1)

	var buf bytes.Buffer
	cmd := RegisterCmd(&grpcclient.Client{
		Client: mockClient,
	})
	cmd.SetOut(&buf)

	input := "testuser\npassword123\n"
	cmd.SetIn(bytes.NewBufferString(input))

	cmd.Run(cmd, []string{})

	assert.Contains(t, buf.String(), "Registration successful!")
}

func TestLoginCmd_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	mockLogger := mlogging.NewMockILogger(ctrl)

	mockClient.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("invalid credentials")).
		Times(1)

	mockLogger.EXPECT().
		Error("Login failed: %v", gomock.Any()).
		Times(1)

	mockClient.EXPECT().GetAllData(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("ошибка получения данных с сервера: %v", gomock.Any())).
		Times(1)

	var buf bytes.Buffer
	cmd := LoginCmd(&grpcclient.Client{
		Client: mockClient,
		Logger: mockLogger,
	})
	cmd.SetOut(&buf)

	input := "wronguser\nwrongpassword\n"

	cmd.SetIn(bytes.NewBufferString(input))

	cmd.Run(cmd, []string{})
}

func TestRegisterCmd_RegistrationFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	mockLogger := mlogging.NewMockILogger(ctrl)

	mockClient.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("user exists")).
		Times(1)

	mockLogger.EXPECT().
		Error("Registration failed: %v", gomock.Any()).
		Times(1)

	var buf bytes.Buffer
	cmd := RegisterCmd(&grpcclient.Client{
		Client: mockClient,
		Logger: mockLogger,
	})
	cmd.SetOut(&buf)

	input := "wronguser\nwrongpassword\n"
	cmd.SetIn(bytes.NewBufferString(input))

	cmd.Run(cmd, []string{})
}

func TestCreateDataCmd(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testData := []byte("test data")
	if _, err := tmpFile.Write(testData); err != nil {
		t.Fatalf("Не удалось записать данные во временный файл: %v", err)
	}
	tmpFile.Close()

	userID := int64(2)
	userDir := filepath.Join("user_data", fmt.Sprintf("%d", userID))
	defer os.RemoveAll(userDir)

	testCases := []struct {
		name           string
		input          string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "Успешное_создание_данных_логин_пароль",
			input:          "1\nunutest\ntest\n\n", // Выбор типа данных, логин, пароль, пустая строка для метаданных
			expectedOutput: "Данные успешно сохранены с ID:",
			expectedError:  false,
		},
		{
			name:           "Успешное создание данных (текстовые данные)",
			input:          "2\nSample text data\n\n", // Выбор типа данных, текстовые данные, пустая строка для метаданных
			expectedOutput: "Данные успешно сохранены с ID:",
			expectedError:  false,
		},
		{
			name:           "Успешное создание данных (банковская карта с метаданными)",
			input:          "4\n4111111111111111\n12/25\n123\nJohn Doe\nkey1=value1\n\n", // Выбор типа данных, данные карты, метаданные
			expectedOutput: "Данные успешно сохранены с ID:",
			expectedError:  false,
		},
		{
			name:           "Невалидные данные (банковская карта)",
			input:          "4\n1234567890123456\n12/25\n123\nJohn Doe\n\n", // Выбор типа данных, данные карты, пустая строка для метаданных
			expectedOutput: "Ошибка: ошибка валидации данных: номер карты не проходит проверку по Луну",
			expectedError:  true,
		},
		{
			name:           "Неверный выбор типа данных",
			input:          "5\n", // Неверный выбор типа данных
			expectedOutput: "Неверный выбор",
			expectedError:  true,
		},
		{
			name:           "Ошибка ввода",
			input:          "\n", // Неверный выбор типа данных
			expectedOutput: "Ошибка ввода. Введите число от 1 до 4.",
			expectedError:  true,
		},
		{
			name:           "Некорректный ввод (текст вместо числа)",
			input:          "abc\n", // Текст вместо числа
			expectedOutput: "Ошибка ввода. Введите число от 1 до 4.",
			expectedError:  true,
		},
		{
			name:           "Ошибка ввода логин_пароль",
			input:          "1\n\n\n\n", // Пустой логин и пароль
			expectedOutput: "Ошибка: ошибка валидации данных: логин и пароль не могут быть пустыми",
			expectedError:  true,
		},
		{
			name:           "Ошибка ввода (пустые текстовые данные)",
			input:          "2\n\n\n", // Пустые текстовые данные
			expectedOutput: "поле не может быть пустым",
			expectedError:  true,
		},
		{
			name:           "Успешное создание данных (бинарные данные)",
			input:          fmt.Sprintf("3\n%s\nTestFileName\n\n", tmpFile.Name()), // Выбор типа данных, путь к файлу, имя файла, пустая строка для метаданных
			expectedOutput: "Данные успешно сохранены с ID:",
			expectedError:  false,
		},
		{
			name:           "Ошибка ввода (несуществующий файл)",
			input:          "3\n/nonexistent/file.txt\nTestFileName\n\n", // Выбор типа данных, несуществующий файл, имя файла, пустая строка для метаданных
			expectedOutput: "Ошибка: ошибка валидации данных: объект данных не инициализирован, проверьте путь к файлу",
			expectedError:  true,
		},
		{
			name:           "Ошибка синхронизации данных",
			input:          "1\nunutest\ntest\n\n",
			expectedOutput: "Ошибка синхронизации данных:",
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mproto.NewMockGophKeeperClient(ctrl)
			if !tc.expectedError {
				mockClient.EXPECT().
					CreateData(gomock.Any(), gomock.Any()).
					Return(&proto.CreateDataResponse{DataId: 123}, nil).
					AnyTimes()
			}

			if tc.expectedError {
				mockClient.EXPECT().
					CreateData(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf(tc.expectedOutput)).
					AnyTimes()
			}

			mockClient.EXPECT().
				GetAllData(gomock.Any(), gomock.Any()).
				Return(&proto.GetAllDataResponse{}, nil).
				AnyTimes()

			masterKey := make([]byte, 32)
			copy(masterKey, "16-byte-master-key")
			client := &grpcclient.Client{
				Client:        mockClient,
				EncryptionKey: masterKey,
			}

			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			r, w, err := os.Pipe()
			assert.NoError(t, err)
			os.Stdin = r

			_, err = w.WriteString(tc.input)
			assert.NoError(t, err)
			w.Close()

			cmd := CreateDataCmd(client)
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err = cmd.Execute()

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Contains(t, buf.String(), tc.expectedOutput)
		})
	}
}

func TestUpdateDataCmd(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testData := []byte("test data")
	if _, err := tmpFile.Write(testData); err != nil {
		t.Fatalf("Не удалось записать данные во временный файл: %v", err)
	}
	tmpFile.Close()

	userID := int64(1)
	filepath.Join("user_data", fmt.Sprintf("%d", userID))

	masterKey := make([]byte, 32)
	copy(masterKey, "16-byte-master-key")

	encryptedData, err := encryption.EncryptData([]byte(`{"username":"unutest","password":"test"}`), masterKey)
	if err != nil {
		t.Fatalf("Не удалось зашифровать тестовые данные: %v", err)
	}

	testDatas := models.Data{
		ID:          1,
		DataType:    "LOGIN_PASSWORD",
		DataContent: []byte(encryptedData),
		Metadata:    nil,
		UpdatedAt:   time.Date(2025, 3, 2, 15, 22, 0, 0, time.FixedZone("MSK", 3*60*60)),
	}

	err = localstorage.SaveData(userID, testDatas)
	if err != nil {
		t.Fatalf("Не удалось сохранить тестовые данные: %v", err)
	}

	testCases := []struct {
		name           string
		input          string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "Успешное_обновлени_данных_логин_пароль",
			input:          "1\nunutest\ntest\n\n", // Выбор типа данных, логин, пароль, пустая строка для метаданных
			expectedOutput: "Данные успешно обновлены!",
			expectedError:  false,
		},
		{
			name:           "Некорректный ID",
			input:          "5\n",
			expectedOutput: "Ошибка: данные с таким ID не найдены",
			expectedError:  true,
		},
		{
			name:           "Ошибка ввода",
			input:          "\n",
			expectedOutput: "Ошибка ввода: unexpected newline",
			expectedError:  true,
		},
		{
			name:           "Некорректный ввод (текст вместо числа)",
			input:          "abc\n",
			expectedOutput: "Ошибка: неверный формат ID",
			expectedError:  true,
		},
		{
			name:           "Ошибка ввода логин_пароль",
			input:          "1\n\n\n\n",
			expectedOutput: "ошибка валидации данных: логин и пароль не могут быть пустыми",
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mproto.NewMockGophKeeperClient(ctrl)

			if !tc.expectedError {
				mockClient.EXPECT().
					CreateData(gomock.Any(), gomock.Any()).
					Return(&proto.CreateDataResponse{DataId: 123}, nil).
					AnyTimes()
			}

			if tc.expectedError {
				mockClient.EXPECT().
					CreateData(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf(tc.expectedOutput)).
					AnyTimes()
			}

			masterKey := make([]byte, 32)
			copy(masterKey, "16-byte-master-key")
			client := &grpcclient.Client{
				Client:        mockClient,
				EncryptionKey: masterKey,
				UserID:        int64(1),
			}

			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			r, w, err := os.Pipe()
			assert.NoError(t, err)
			os.Stdin = r

			_, err = w.WriteString(tc.input)
			assert.NoError(t, err)
			w.Close()

			cmd := UpdateDataCmd(client)
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err = cmd.Execute()

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Contains(t, buf.String(), tc.expectedOutput)
		})
	}
}

func TestDeleteDataCmd(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
		expectedError  bool
		mockBehavior   func(*mproto.MockGophKeeperClient)
	}{
		{
			name:           "Успешное удаление данных",
			input:          "123\n",
			expectedOutput: "Данные успешно удалены",
			expectedError:  false,
			mockBehavior: func(mockClient *mproto.MockGophKeeperClient) {
				mockClient.EXPECT().
					DeleteData(gomock.Any(), gomock.Any()).
					Return(&proto.DeleteDataResponse{}, nil).
					AnyTimes()
			},
		},
		{
			name:           "Ошибка удаления данных",
			input:          "1\n",
			expectedOutput: "Ошибка удаления данных:",
			expectedError:  true,
			mockBehavior: func(mockClient *mproto.MockGophKeeperClient) {
				mockClient.EXPECT().
					DeleteData(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("ошибка удаления данных")).
					AnyTimes()
			},
		},
		{
			name:           "Ошибка ввода ID (некорректный формат)",
			input:          "abc\n",
			expectedOutput: "Ошибка: неверный формат ID",
			expectedError:  true,
			mockBehavior: func(mockClient *mproto.MockGophKeeperClient) {
				mockClient.EXPECT().
					DeleteData(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("неверный формат ID")).
					AnyTimes()
			},
		},
		{
			name:           "Ошибка ввода",
			input:          "\n",
			expectedOutput: "Ошибка ввода: unexpected newline",
			expectedError:  true,
			mockBehavior: func(mockClient *mproto.MockGophKeeperClient) {
				mockClient.EXPECT().
					DeleteData(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("ошибка ввода")).
					AnyTimes()
			},
		},
		{
			name:           "нет данных для удаления",
			input:          "\n",
			expectedOutput: "У вас нет данных для удаления",
			expectedError:  true,
			mockBehavior: func(mockClient *mproto.MockGophKeeperClient) {
				os.RemoveAll("user_data")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mproto.NewMockGophKeeperClient(ctrl)

			tc.mockBehavior(mockClient)

			masterKey := make([]byte, 32)
			copy(masterKey, "16-byte-master-key")
			client := &grpcclient.Client{
				Client:        mockClient,
				EncryptionKey: masterKey,
				UserID:        int64(0),
			}

			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			r, w, err := os.Pipe()
			assert.NoError(t, err)
			os.Stdin = r

			_, err = w.WriteString(tc.input)
			assert.NoError(t, err)
			w.Close()

			cmd := DeleteDataCmd(client)
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err = cmd.Execute()

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Contains(t, buf.String(), tc.expectedOutput)
		})
	}
}

func TestDeleteDataCmd_ErrorGettingData(t *testing.T) {
	t.Cleanup(func() {
		os.RemoveAll("user_data")
	})

	userID := int64(1)
	userDir := filepath.Join("user_data", fmt.Sprintf("%d", userID))
	err := os.MkdirAll(userDir, 0700)
	if err != nil {
		t.Fatalf("Не удалось создать директорию пользователя: %v", err)
	}

	dataFilePath := filepath.Join(userDir, "data.json")
	err = os.WriteFile(dataFilePath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Не удалось создать пустой файл data.json: %v", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)

	mockClient.EXPECT().
		DeleteData(gomock.Any(), gomock.Any()).
		Return(nil,
			fmt.Errorf("ошибка получения данных из локального хранилища: "+
				"ошибка чтения данных: ошибка десериализации данных: unexpected end of JSON input")).
		AnyTimes()

	masterKey := make([]byte, 32)
	copy(masterKey, "16-byte-master-key")
	client := &grpcclient.Client{
		Client:        mockClient,
		EncryptionKey: masterKey,
		UserID:        userID,
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdin = r

	_, err = w.WriteString("\n")
	assert.NoError(t, err)
	w.Close()

	cmd := DeleteDataCmd(client)
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err = cmd.Execute()

	assert.Error(t, err)
	output := buf.String()
	assert.Contains(t, output, "Ошибка получения данных: ошибка получения данных из локального хранилища: ошибка чтения данных пользователя: ошибка десериализации данных: unexpected end of JSON input")
}

func TestGetDataCmd_Integration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)
	t.Cleanup(func() {
		os.RemoveAll("user_data")
	})

	userID := int64(1)
	userDir := filepath.Join("user_data", fmt.Sprintf("%d", userID))
	err := os.MkdirAll(userDir, 0700)
	if err != nil {
		t.Fatalf("Не удалось создать директорию пользователя: %v", err)
	}

	masterKey := make([]byte, 32)
	copy(masterKey, "16-byte-master-key")

	encryptedData, err := encryption.EncryptData([]byte(`{"username":"unutest","password":"test"}`), masterKey)
	if err != nil {
		t.Fatalf("Не удалось зашифровать тестовые данные: %v", err)
	}

	testData := models.Data{
		ID:          1,
		DataType:    "LOGIN_PASSWORD",
		DataContent: []byte(encryptedData),
		Metadata:    nil,
		UpdatedAt:   time.Date(2025, 3, 2, 15, 22, 0, 0, time.FixedZone("MSK", 3*60*60)),
	}

	err = localstorage.SaveData(userID, testData)
	if err != nil {
		t.Fatalf("Не удалось сохранить тестовые данные: %v", err)
	}

	client := &grpcclient.Client{
		Client:        mockClient,
		EncryptionKey: masterKey,
		UserID:        int64(1),
	}

	cmd := GetDataCmd(client)

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	cmd.Run(cmd, []string{})

	output := buf.String()
	assert.Contains(t, output, "ID: 1")
	assert.Contains(t, output, "Тип данных: LOGIN_PASSWORD")
	assert.Contains(t, output, `Содержимое: {"username":"unutest","password":"test"}`)
	assert.Contains(t, output, "Метаданные: map[]")
	assert.Regexp(t, `Обновлено: 2025-03-02 15:22:00 \+0300( MSK)?\n---\n`, output)
}

func TestGetDataCmd_Integration_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mproto.NewMockGophKeeperClient(ctrl)
	t.Cleanup(func() {
		os.RemoveAll("user_data")
	})

	userID := int64(1)
	userDir := filepath.Join("user_data", fmt.Sprintf("%d", userID))
	err := os.MkdirAll(userDir, 0700)
	if err != nil {
		t.Fatalf("Не удалось создать директорию пользователя: %v", err)
	}

	masterKey := make([]byte, 32)
	copy(masterKey, "16-byte-master-key")

	testData := models.Data{
		ID:          1,
		DataType:    "LOGIN_PASSWORD",
		DataContent: []byte(`{"username":"unutest","password":"test"}`), // Не зашифрованные данные
		Metadata:    nil,
		UpdatedAt:   time.Date(2025, 3, 2, 15, 22, 0, 0, time.FixedZone("MSK", 3*60*60)),
	}

	err = localstorage.SaveData(userID, testData)
	if err != nil {
		t.Fatalf("Не удалось сохранить тестовые данные: %v", err)
	}

	client := &grpcclient.Client{
		Client:        mockClient,
		EncryptionKey: masterKey,
		UserID:        int64(1),
	}

	cmd := GetDataCmd(client)

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	cmd.Run(cmd, []string{})

	output := buf.String()
	assert.Contains(t, output, "Ошибка получения данных:")
}
