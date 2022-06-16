package database

import (
	"database/sql/driver"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/sioncojp/famili-api/domain"
	"github.com/sioncojp/famili-api/domain/model"
	"github.com/sioncojp/famili-api/utils"
)

// テストスイートの構造体
type TodoRepositoryTestSuite struct {
	suite.Suite
	mock           sqlmock.Sqlmock
	todoRepository todoRepository
	dummy          *model.Todo
	dummys         []model.Todo
}

// テストのセットアップ
// (sqlmockをNew、Gormで発行されるクエリがモックに送られるように)
func (s *TodoRepositoryTestSuite) SetupTest() {
	//db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	todoRepository := todoRepository{}
	todoRepository.db, _ = gorm.Open(
		mysql.Dialector{Config: &mysql.Config{DriverName: "mysql", Conn: db, SkipInitializeWithVersion: true}},
		&gorm.Config{},
	)
	s.mock = mock
	s.todoRepository = todoRepository
}

// テスト終了時の処理（データベース接続のクローズ）
func (s *TodoRepositoryTestSuite) TearDownTest() {
	db, _ := s.todoRepository.db.DB()
	db.Close()
}

// テストスイートの実行
func TestTodoRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TodoRepositoryTestSuite))
}

// BeforeTest...テストスイート前に実行するメソッド
func (s *TodoRepositoryTestSuite) BeforeTest(suiteName string, testName string) {
	s.dummy = &model.Todo{
		Model:       model.Model{ID: utils.MakeRandomUintExcludeZero(100)},
		Title:       faker.Word(),
		Description: faker.Word(),
		Completed:   false,
	}

	s.dummys = []model.Todo{
		{
			Model:       model.Model{ID: utils.MakeRandomUintExcludeZero(100)},
			Title:       faker.Word(),
			Description: faker.Sentence(),
			Completed:   false,
		},
		{
			Model:       model.Model{ID: utils.MakeRandomUintExcludeZero(100)},
			Title:       faker.Word(),
			Description: faker.Sentence(),
			Completed:   true,
		},
	}
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

var anyTime = AnyTime{}

func (s *TodoRepositoryTestSuite) TestTodoGetById() {
	s.Run("GetById", func() {
		rows := s.mock.NewRows([]string{"id", "title", "description", "completed"}).
			AddRow(s.dummy.ID, s.dummy.Title, s.dummy.Description, s.dummy.Completed)

		s.mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT * FROM `todos` WHERE id = ? ORDER BY `todos`.`id` LIMIT 1")).
			WithArgs(strconv.FormatUint(uint64(s.dummy.ID), 10)).
			WillReturnRows(rows)

		data, err := s.todoRepository.GetById(domain.Id(strconv.Itoa(int(s.dummy.ID))))
		require.NoError(s.T(), err)

		assert.Equal(s.T(), data.ID, s.dummy.ID, "unexpected id")
		assert.Equal(s.T(), data.Title, s.dummy.Title, "unexpected title")
		assert.Equal(s.T(), data.Description, s.dummy.Description, "unexpected description")
		assert.Equal(s.T(), data.Completed, s.dummy.Completed, "unexpected completed")
	})
}

func (s *TodoRepositoryTestSuite) TestTodoList() {
	s.Run("List", func() {
		rows := sqlmock.NewRows([]string{"id", "title", "description", "completed"})
		for _, v := range s.dummys {
			rows.AddRow(v.ID, v.Title, v.Description, v.Completed)
		}
		s.mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT * FROM `todos`")).
			WillReturnRows(rows)

		data, err := s.todoRepository.List()
		require.NoError(s.T(), err)

		num := 0
		for _, v := range data {
			assert.Equal(s.T(), v.ID, s.dummys[num].ID, "unexpected id")
			assert.Equal(s.T(), v.Title, s.dummys[num].Title, "unexpected title")
			assert.Equal(s.T(), v.Description, s.dummys[num].Description, "unexpected description")
			assert.Equal(s.T(), v.Completed, s.dummys[num].Completed, "unexpected completed")
			num++
		}
	})
}

func (s *TodoRepositoryTestSuite) TestTodoCreate() {
	s.Run("Create", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec("INSERT").
			WithArgs(anyTime, anyTime, s.dummy.Title, s.dummy.Description, s.dummy.Completed).
			WillReturnResult(sqlmock.NewResult(int64(s.dummy.ID), 1))
		s.mock.ExpectCommit()

		data := &model.Todo{
			Title:       s.dummy.Title,
			Description: s.dummy.Description,
			Completed:   s.dummy.Completed,
		}
		err := s.todoRepository.Create(data)
		require.NoError(s.T(), err)
		if err == nil {
			assert.Equal(s.T(), data.ID, s.dummy.ID, "unexpected ID")
			assert.Equal(s.T(), data.Title, s.dummy.Title, "unexpected title")
			assert.Equal(s.T(), data.Description, s.dummy.Description, "unexpected description")
			assert.Equal(s.T(), data.Completed, s.dummy.Completed, "unexpected completed")
		}
	})
}

func (s *TodoRepositoryTestSuite) TestTodoUpdate() {
	s.Run("Update", func() {
		data := &model.Todo{
			Model:       model.Model{ID: s.dummy.ID},
			Title:       faker.Word(),
			Description: faker.Sentence(),
			Completed:   true,
		}

		s.mock.ExpectBegin()
		s.mock.ExpectExec("UPDATE").
			WithArgs(anyTime, anyTime, data.Title, data.Description, data.Completed, data.ID).
			WillReturnResult(sqlmock.NewResult(int64(s.dummy.ID), 1))
		s.mock.ExpectCommit()

		err := s.todoRepository.Update(data)
		require.NoError(s.T(), err)
		if err == nil {
			assert.Equal(s.T(), data.ID, s.dummy.ID, "unexpected ID")
			assert.NotEqual(s.T(), data.Title, s.dummy.Title, "unexpected title")
			assert.NotEqual(s.T(), data.Description, s.dummy.Description, "unexpected description")
			assert.NotEqual(s.T(), data.Completed, s.dummy.Completed, "unexpected completed")
		}
	})
}

func (s *TodoRepositoryTestSuite) TestTodoDelete() {
	s.Run("Delete", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec("DELETE").
			WithArgs(s.dummy.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.todoRepository.Delete(s.dummy)
		require.NoError(s.T(), err)
	})
}
