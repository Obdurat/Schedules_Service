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
	Where string
	WantStatus int
	WantReturn string
	Mode string
}

var tests = []testStruct{
	{
		Name: "MONGO ERROR",
		Where: `{"price": 80}`,
		WantStatus: 400,
		WantReturn: `{"error":"MOCK MONGO ERROR: Find failed"}`,
		Mode: "MONGO_FAIL",
	},
	{
		Name: "QUERY SUCCESS",
		Where: `{"price": 80}`,
		WantStatus: 200,
		WantReturn: `{"result":[{"_id":"641b551ad15248336f0aad6a","client_id":"abduhlaziz","service":"641afadff6872fffc607baef","price":80,"date":"2019","finished":false}]}`,
		Mode: "MONGO_SUCCESS",
	},
}

func getAllScheduleTest(where string, want_status int, want_return string, mode string) (error) {
	repository.Instance = new(mode)
	w := httptest.NewRecorder()
	c, _ := Prepare(w)
	req, err := http.NewRequest("DELETE", "/schedules/any/any", nil); if err != nil {
		return err
	}
	c.DefaultQuery("where", where)
	c.Request = req
	handlers.GetSchedules(c)
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
		if err := getAllScheduleTest(test.Where, test.WantStatus, test.WantReturn, test.Mode); err != nil {
			fmt.Printf("%s --- FAILED: %v\n", test.Name, err.Error())
			t.Fail()
			continue
		}
		fmt.Printf("%s --- PASSED\n", test.Name)
	}
}