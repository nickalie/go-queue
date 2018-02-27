package queue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
	"errors"
)

type concurrentStringSlice struct {
	*sync.Mutex
	items []string
}

func (cs *concurrentStringSlice) Append(item string) {
	cs.Lock()
	defer cs.Unlock()
	cs.items = append(cs.items, item)
}

func (cs *concurrentStringSlice) Contains(item string) bool {
	cs.Lock()
	defer cs.Unlock()
	for _, v := range cs.items {
		if v == item {
			return true
		}
	}

	return false
}

func (cs *concurrentStringSlice) Len() int {
	cs.Lock()
	defer cs.Unlock()
	return len(cs.items)
}

type baseSuite struct {
	suite.Suite
}

func (suite *baseSuite) TestMulti() {
	t := suite.T()
	cs := &concurrentStringSlice{Mutex: &sync.Mutex{}}
	consumers := rand.Intn(10) + 10
	wg := &sync.WaitGroup{}
	wg.Add(consumers)

	for i := 0; i < consumers; i++ {
		go consumer(t, wg, cs)
	}

	messages := consumers * (rand.Intn(100) + 100)

	for i := 0; i < messages; i++ {
		err := Put("messages", fmt.Sprintf("message %d", i))
		assert.Nil(t, err)
	}

	for i := 0; i < consumers; i++ {
		Put("messages", "done")
	}

	wg.Wait()
	assert.Equal(t, messages, cs.Len())
}

func consumer(t *testing.T, wg *sync.WaitGroup, cs *concurrentStringSlice) {
	defer wg.Done()
	for {
		var message string
		err := Get("messages", &message)
		assert.Nil(t, err)
		assert.NotZero(t, message)

		if message == "done" {
			return
		}

		assert.False(t, cs.Contains(message))
		cs.Append(message)
	}
}

func (suite *baseSuite) TestString() {
	t := suite.T()
	key := randString(10)
	value := randString(20)
	err := Put(key, value)
	assert.Nil(t, err)
	var result string
	err = Get(key, &result)
	assert.Nil(t, err)
	assert.Equal(t, value, result)
}

func (suite *baseSuite) TestStringError() {
	t := suite.T()
	key := randString(10)
	value := randString(20)
	err := Put(key, value)
	assert.Nil(t, err)
	value = randString(20)
	var result string
	err = Get(key, &result)
	assert.Nil(t, err)
	assert.NotEqual(t, value, result)
}

func (suite *baseSuite) TestMap() {
	t := suite.T()
	key := randString(10)
	value := randMap(100)
	err := Put(key, value)
	assert.Nil(t, err)
	var result map[string]interface{}
	err = Get(key, &result)
	assert.Nil(t, err)
	assert.Equal(t, value, result)
}

func (suite *baseSuite) TestMapError() {
	t := suite.T()
	key := randString(10)
	value := randMap(100)
	err := Put(key, value)
	assert.Nil(t, err)
	value = randMap(100)
	var result map[string]interface{}
	err = Get(key, &result)
	assert.Nil(t, err)
	assert.NotEqual(t, value, result)
}

func (suite *baseSuite) TestFloat() {
	t := suite.T()
	key := randString(10)
	value := rand.Float64()
	err := Put(key, value)
	assert.Nil(t, err)
	var result float64
	err = Get(key, &result)
	assert.Nil(t, err)
	assert.Equal(t, value, result)
}

func (suite *baseSuite) TestFloatError() {
	t := suite.T()
	key := randString(10)
	value := rand.Float64()
	err := Put(key, value)
	assert.Nil(t, err)
	value = rand.Float64()
	var result float64
	err = Get(key, &result)
	assert.Nil(t, err)
	assert.NotEqual(t, value, result)
}

func (suite *baseSuite) TestInt() {
	t := suite.T()
	key := randString(10)
	value := rand.Int63()
	err := Put(key, value)
	assert.Nil(t, err)
	var result int64
	err = Get(key, &result)
	assert.Nil(t, err)
	assert.Equal(t, value, result)
}

func (suite *baseSuite) TestIntError() {
	t := suite.T()
	key := randString(10)
	value := rand.Int63()
	err := Put(key, value)
	assert.Nil(t, err)
	value = rand.Int63()
	var result int64
	err = Get(key, &result)
	assert.Nil(t, err)
	assert.NotEqual(t, value, result)
}

func (suite *baseSuite) TestObject() {
	t := suite.T()
	key := randString(10)
	value := randUser()
	err := Put(key, value)
	assert.Nil(t, err)
	var result testUser
	err = Get(key, &result)
	assert.Nil(t, err)
	assert.Nil(t, equalUsers(value, result))
}

func (suite *baseSuite) TestObjectAsync() {
	t := suite.T()
	key := randString(10)
	value := randUser()

	go func() {
		timer := time.NewTimer(time.Second)
		<-timer.C
		err := Put(key, value)
		assert.Nil(t, err)
	}()

	var result testUser
	err := Get(key, &result)
	assert.Nil(t, err)
	assert.Nil(t, equalUsers(value, result))
}

func (suite *baseSuite) TestObjectError() {
	t := suite.T()
	key := randString(10)
	value := randUser()
	err := Put(key, value)
	assert.Nil(t, err)
	value = randUser()
	var result testUser
	err = Get(key, &result)
	assert.Nil(t, err)
	assert.NotNil(t, equalUsers(value, result))
}


func randMap(keys int) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < keys; i++ {
		d := i % 3
		if d == 0 {
			m[randString(20)] = randString(100)
		} else {
			m[randString(20)] = rand.Float64()
		}
	}

	return m
}

type testName struct {
	FirstName string
	LastName  string
}

type testCompany struct {
	Name    string
	Domains []string
}

type testUser struct {
	Name       testName
	Companies  []testCompany
	Address    string
	Birthday   time.Time
	Duration   time.Duration
	Count      int64
	CountFloat float64
}

func (t *testUser) String() string {
	return t.Name.FirstName + " " + t.Name.LastName + " " + strconv.Itoa(len(t.Companies))
}

func randUser() testUser {
	name := testName{FirstName: randString(20), LastName: randString(10)}
	companies := make([]testCompany, rand.Intn(10)+10)

	for index := range companies {
		companies[index] = randCompany()
	}

	return testUser{
		Name:       name,
		Companies:  companies,
		Address:    randString(30),
		Birthday:   time.Now(),
		Duration:   time.Duration(rand.Intn(20)) * time.Minute,
		Count:      rand.Int63(),
		CountFloat: rand.Float64(),
	}
}

func randCompany() testCompany {
	domains := make([]string, rand.Intn(10)+10)

	for index := range domains {
		domains[index] = randString(20)
	}

	return testCompany{Name: randString(20), Domains: domains}
}

func equalUsers(user1, user2 testUser) error {
	if user1.Name != user2.Name {
		return errors.New("user names should be equal")
	}

	if user1.Address != user2.Address {
		return errors.New("user addresses should be equal")
	}

	if !user1.Birthday.Equal(user2.Birthday) {
		return errors.New("user birthdays should be equal")
	}

	if user1.Duration != user2.Duration {
		return errors.New("user durations should be equal")
	}

	if user1.Count != user2.Count {
		return errors.New("user counts should be equal")
	}

	if user1.CountFloat != user2.CountFloat {
		return errors.New("user countFloats should be equal")
	}

	if len(user1.Companies) != len(user2.Companies) {
		return errors.New("user companies len should be equal")
	}

	for i, v := range user1.Companies {
		err := equalCompanies(v, user2.Companies[i])

		if err != nil {
			return fmt.Errorf("user companies at index %d should be equal: %v", i, err)
		}
	}

	return nil
}

func equalCompanies(company1, company2 testCompany) error {
	if company1.Name != company2.Name {
		return errors.New("company names should be equal")
	}

	if len(company1.Domains) != len(company2.Domains) {
		return errors.New("company domains len should be equal")
	}

	for i, v := range company1.Domains {
		if v != company2.Domains[i] {
			return fmt.Errorf("company domains at index %d should be equal", i)
		}
	}

	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const letterBytesLen = len(letterBytes)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(letterBytesLen)]
	}
	return string(b)
}
