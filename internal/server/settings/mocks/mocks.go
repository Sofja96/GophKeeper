package mocks

import (
	"reflect"

	"github.com/golang/mock/gomock"

	_ "github.com/Sofja96/GophKeeper.git/internal/server/settings"
)

// MockSettings мок для settings.Settings
type MockSettings struct {
	ctrl     *gomock.Controller
	recorder *MockSettingsMockRecorder
}

// MockSettingsMockRecorder мок-рекордер для MockSettings
type MockSettingsMockRecorder struct {
	mock *MockSettings
}

// NewMockSettings создает новый мок для settings.Settings
func NewMockSettings(ctrl *gomock.Controller) *MockSettings {
	mock := &MockSettings{ctrl: ctrl}
	mock.recorder = &MockSettingsMockRecorder{mock}
	return mock
}

// EXPECT возвращает объект, который позволяет указывать ожидаемые вызовы
func (m *MockSettings) EXPECT() *MockSettingsMockRecorder {
	return m.recorder
}

// DbDsn мок для метода DbDsn
func (m *MockSettings) DbDsn() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DbDsn")
	ret0, _ := ret[0].(string)
	return ret0
}

// DbDsn указывает ожидаемый вызов DbDsn
func (mr *MockSettingsMockRecorder) DbDsn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DbDsn", reflect.TypeOf((*MockSettings)(nil).DbDsn))
}

// DbAutoMigration мок для метода DbAutoMigration
func (m *MockSettings) DbAutoMigration() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DbAutoMigration")
	ret0, _ := ret[0].(bool)
	return ret0
}

// DbAutoMigration указывает ожидаемый вызов DbAutoMigration
func (mr *MockSettingsMockRecorder) DbAutoMigration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DbAutoMigration", reflect.TypeOf((*MockSettings)(nil).DbAutoMigration))
}
