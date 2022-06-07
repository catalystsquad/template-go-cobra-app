package test

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ExampleSuite struct {
	suite.Suite
}

// called before the entire suite is run
func (s *ExampleSuite) SetupSuite() {}

// called after the entire suite is run
func (s *ExampleSuite) TearDownSuite() {}

// called before each test
func (s *ExampleSuite) SetupTest() {}

// called after each test
func (s *ExampleSuite) TearDownTest() {}

// runs the entire suite
func TestExampleSuite(t *testing.T) {
	suite.Run(t, new(ExampleSuite))
}

// an individual test example
func (s *ExampleSuite) TestExample() {
	require.True(s.T(), true)
}
