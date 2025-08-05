package unit_test

import (
	"github.com/stretchr/testify/mock"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Config() *schemaregistry.Config {
	args := m.Called()
	return args.Get(0).(*schemaregistry.Config)
}

func (m *MockClient) GetAllContexts() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockClient) Register(subject string, schema schemaregistry.SchemaInfo, normalize bool) (int, error) {
	args := m.Called(subject, schema, normalize)
	return args.Int(0), args.Error(1)
}

func (m *MockClient) RegisterFullResponse(subject string, schema schemaregistry.SchemaInfo, normalize bool) (schemaregistry.SchemaMetadata, error) {
	args := m.Called(subject, schema, normalize)
	return args.Get(0).(schemaregistry.SchemaMetadata), args.Error(1)
}

func (m *MockClient) GetBySubjectAndID(subject string, id int) (schemaregistry.SchemaInfo, error) {
	args := m.Called(subject, id)
	return args.Get(0).(schemaregistry.SchemaInfo), args.Error(1)
}

func (m *MockClient) GetByGUID(guid string) (schemaregistry.SchemaInfo, error) {
	args := m.Called(guid)
	return args.Get(0).(schemaregistry.SchemaInfo), args.Error(1)
}

func (m *MockClient) GetSubjectsAndVersionsByID(id int) ([]schemaregistry.SubjectAndVersion, error) {
	args := m.Called(id)
	return args.Get(0).([]schemaregistry.SubjectAndVersion), args.Error(1)
}

func (m *MockClient) GetID(subject string, schema schemaregistry.SchemaInfo, normalize bool) (int, error) {
	args := m.Called(subject, schema, normalize)
	return args.Int(0), args.Error(1)
}

func (m *MockClient) GetIDFullResponse(subject string, schema schemaregistry.SchemaInfo, normalize bool) (schemaregistry.SchemaMetadata, error) {
	args := m.Called(subject, schema, normalize)
	return args.Get(0).(schemaregistry.SchemaMetadata), args.Error(1)
}

func (m *MockClient) GetLatestSchemaMetadata(subject string) (schemaregistry.SchemaMetadata, error) {
	args := m.Called(subject)
	return args.Get(0).(schemaregistry.SchemaMetadata), args.Error(1)
}

func (m *MockClient) GetSchemaMetadata(subject string, version int) (schemaregistry.SchemaMetadata, error) {
	args := m.Called(subject, version)
	return args.Get(0).(schemaregistry.SchemaMetadata), args.Error(1)
}

func (m *MockClient) GetSchemaMetadataIncludeDeleted(subject string, version int, deleted bool) (schemaregistry.SchemaMetadata, error) {
	args := m.Called(subject, version, deleted)
	return args.Get(0).(schemaregistry.SchemaMetadata), args.Error(1)
}

func (m *MockClient) GetLatestWithMetadata(subject string, metadata map[string]string, deleted bool) (schemaregistry.SchemaMetadata, error) {
	args := m.Called(subject, metadata, deleted)
	return args.Get(0).(schemaregistry.SchemaMetadata), args.Error(1)
}

func (m *MockClient) GetAllVersions(subject string) ([]int, error) {
	args := m.Called(subject)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockClient) GetVersion(subject string, schema schemaregistry.SchemaInfo, normalize bool) (int, error) {
	args := m.Called(subject, schema, normalize)
	return args.Int(0), args.Error(1)
}

func (m *MockClient) GetVersionIncludeDeleted(subject string, schema schemaregistry.SchemaInfo, normalize bool, deleted bool) (int, error) {
	args := m.Called(subject, schema, normalize, deleted)
	return args.Int(0), args.Error(1)
}

func (m *MockClient) GetAllSubjects() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockClient) DeleteSubject(subject string, permanent bool) ([]int, error) {
	args := m.Called(subject, permanent)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockClient) DeleteSubjectVersion(subject string, version int, permanent bool) (int, error) {
	args := m.Called(subject, version, permanent)
	return args.Int(0), args.Error(1)
}

func (m *MockClient) TestSubjectCompatibility(subject string, schema schemaregistry.SchemaInfo) (bool, error) {
	args := m.Called(subject, schema)
	return args.Bool(0), args.Error(1)
}

func (m *MockClient) TestCompatibility(subject string, version int, schema schemaregistry.SchemaInfo) (bool, error) {
	args := m.Called(subject, version, schema)
	return args.Bool(0), args.Error(1)
}

func (m *MockClient) GetCompatibility(subject string) (schemaregistry.Compatibility, error) {
	args := m.Called(subject)
	return args.Get(0).(schemaregistry.Compatibility), args.Error(1)
}

func (m *MockClient) UpdateCompatibility(subject string, update schemaregistry.Compatibility) (schemaregistry.Compatibility, error) {
	args := m.Called(subject, update)
	return args.Get(0).(schemaregistry.Compatibility), args.Error(1)
}

func (m *MockClient) GetDefaultCompatibility() (schemaregistry.Compatibility, error) {
	args := m.Called()
	return args.Get(0).(schemaregistry.Compatibility), args.Error(1)
}

func (m *MockClient) UpdateDefaultCompatibility(update schemaregistry.Compatibility) (schemaregistry.Compatibility, error) {
	args := m.Called(update)
	return args.Get(0).(schemaregistry.Compatibility), args.Error(1)
}

func (m *MockClient) GetConfig(subject string, defaultToGlobal bool) (schemaregistry.ServerConfig, error) {
	args := m.Called(subject, defaultToGlobal)
	return args.Get(0).(schemaregistry.ServerConfig), args.Error(1)
}

func (m *MockClient) UpdateConfig(subject string, update schemaregistry.ServerConfig) (schemaregistry.ServerConfig, error) {
	args := m.Called(subject, update)
	return args.Get(0).(schemaregistry.ServerConfig), args.Error(1)
}

func (m *MockClient) GetDefaultConfig() (schemaregistry.ServerConfig, error) {
	args := m.Called()
	return args.Get(0).(schemaregistry.ServerConfig), args.Error(1)
}

func (m *MockClient) UpdateDefaultConfig(update schemaregistry.ServerConfig) (schemaregistry.ServerConfig, error) {
	args := m.Called(update)
	return args.Get(0).(schemaregistry.ServerConfig), args.Error(1)
}

func (m *MockClient) ClearLatestCaches() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockClient) ClearCaches() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockClient) Close() error {
	args := m.Called()
	return args.Error(0)
}
