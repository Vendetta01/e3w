package auth_test

/*import (
	"testing"

	auth "github.com/VendettA01/e3w/auth"
	mock_auth "github.com/VendettA01/e3w/auth/mocks"
	"github.com/golang/mock/gomock"
)*/

/*func setup() {
	test1 := mocks.userAuthentication{}
	test1.On("init").Return(true)
	test1.On("login").Return(false)
	test1.On("getName").Return("test1")
	test2 := mocks.userAuthentication{}
	test2.On("init").Return(true)
	test2.On("login", &userCredentials{Username: "testuser", Password: "testpw"}).Return(true)
	//test2.On("login").Return(false)
	test2.On("getName").Return("test2")
	authImpls = make(map[string]userAuthentication)
	activeAuths = make([]string, 2)
	authImpls["test1"] = test1
	authImpls["test2"] = test2
	activeAuths = append(activeAuths, "test1")
	activeAuths = append(activeAuths, "test2")
}

func shutdown() {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}*/

/*func TestCanLogInValidUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserAuth1 := mock_auth.NewMockUserAuthentication(mockCtrl)
	mockUserAuth1.EXPECT().Init().Return(true).Times(1)
	mockUserAuth1.EXPECT().GetName().Return("mockUserAuth1").Times(2)

	mockUserAuth2 := mock_auth.NewMockUserAuthentication(mockCtrl)
	mockUserAuth2.EXPECT().Init().Return(true).Times(1)
	mockUserAuth1.EXPECT().GetName().Return("mockUserAuth2").Times(2)

	auth.

	userCreds := auth.UserCredentials{
		Username: "testuser",
		Password: "testpw",
	}
	expected := true
	res, err := auth.CanLogIn(userCreds)

	if err != nil || res != expected {
		t.Fail()
	}
}

/*func TestCanLogInInvalidUser(t *testing.T) {
	userCreds := &UserCredentials{
		Username: "testuserXXX",
		Password: "testpwXXX",
	}
	assert.False(t, canLogin(userCreds))
	//t.Error("Test not implemented yet")
}

func TestInitAuthFromConf(t *testing.T) {
	t.Error("Test not implemented yet")
}*/
