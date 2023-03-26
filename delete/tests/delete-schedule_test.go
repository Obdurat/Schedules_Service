package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/Obdurat/Schedules/handlers"
	repository "github.com/Obdurat/Schedules/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Prepare(w *httptest.ResponseRecorder) (*gin.Context, *gin.Engine) {
	gin.SetMode("release")
	return gin.CreateTestContext(w)
}

type testStruct struct {
	Name string
	ID string
	WantStatus int
	WantReturn string
	Mode string
}

var tests = []testStruct{
	{
		Name: "MONGO ERROR",
		ID: "641f44e939bddf7b7f8b756f",
		WantStatus: 400,
		WantReturn: `{"error":"MOCK ERROR: FindOneAndDelete"}`,
		Mode: "MONGO_FAIL",
	},
	{
		Name: "MALFORMED ID",
		ID: "641f44e939bddf7b7f8b756",
		WantStatus: 400,
		WantReturn: `{"error":"the provided hex string is not a valid ObjectID"}`,
		Mode: "MONGO_SUCCESS",
	},
	{
		Name: "DELETE SUCCESS",
		ID: "641f44e939bddf7b7f8b756f",
		WantStatus: 200,
		WantReturn: `{"message":"Deleted sucessfully"}`,
		Mode: "MONGO_SUCCESS",
	},
}

func deleteScheduleTest(id string, want_status int, want_return string, mode string) (error) {
	repository.Instance = new(mode)
	w := httptest.NewRecorder()
	c, _ := Prepare(w)
	req, err := http.NewRequest("DELETE", "/schedules/any/any" + id, nil); if err != nil {
		return err
	}
	c.AddParam("id", id)
	c.Request = req
	handlers.DeleteSchedule(c)
	if want_status != w.Code {
		return fmt.Errorf("Status Code unexpected: want: %d, received: %d", want_status, w.Code)
	}
	if want_return != w.Body.String() {
		return fmt.Errorf("Return unexpected: want: %s, received: %s", want_return, w.Body.String())
	}
	return nil
}

func Test_Delete_Schedule(t *testing.T) {
	logrus.SetOutput(io.Discard)
	for _, test := range tests {
		if err := deleteScheduleTest(test.ID, test.WantStatus, test.WantReturn, test.Mode); err != nil {
			fmt.Printf("%s --- FAILED: %v\n", test.Name, err.Error())
			t.Fail()
			continue
		}
		fmt.Printf("%s --- PASSED\n", test.Name)
	}
}