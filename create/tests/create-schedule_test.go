package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Obdurat/Schedules/handlers"
	"github.com/Obdurat/Schedules/mongo"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Prepare(w *httptest.ResponseRecorder) (*gin.Context, *gin.Engine) {
	gin.SetMode("release")
	return gin.CreateTestContext(w)
}

type testStruct struct {
	Name string
	Body string
	WantStatus int
	WantReturn string
	Mode string
}

var tests = []testStruct{
	{
		Name: "MONGO ERROR",
		Body: `{
			"client_id": "Babayeteu",
			"service": "641afc1ff6872fffc607baf6",
			"price": 100,
			"date": "2024",
			"finished": false
		}`,
		WantStatus: 400,
		WantReturn: `{"error":"inducedError"}`,
		Mode: "MONGO_FAIL",
	},
	{
		Name: "JSON VALIDATION",
		Body: `{
			"client_id": "Babayeteu",
			"service": "641afc1ff6872fffc607baf6",
			"price": 100,
			"finished": false
		}`,
		WantStatus: 400,
		WantReturn: `{"error":[{"param":"Date","message":"Date is required"}]}`,
		Mode: "MONGO_SUCCESS",
	},
	{
		Name: "MALFORMED JSON",
		Body: `{
			"client_id": "Babayeteu",
			"service": "641afc1ff6872fffc607baf6",
			"price": 100,
			"date": "2024"
			"finished": false
		}`,
		WantStatus: 400,
		WantReturn: `{"error":"invalid character '\"' after object key:value pair"}`,
		Mode: "MONGO_SUCCESS",
	},
	{
		Name: "INSERT SUCCESS",
		Body: `{
			"client_id": "Babayeteu",
			"service": "641afc1ff6872fffc607baf6",
			"price": 100,
			"date": "2024",
			"finished": false
		}`,
		WantStatus: 201,
		WantReturn: `{"message":"Created sucessfully"}`,
		Mode: "MONGO_SUCCESS",
	},
}

func createScheduleTest(body string, want_status int, want_return string, mode string) (error) {
	mongo.Repo = new(mode)
	w := httptest.NewRecorder()
	c, _ := Prepare(w)
	req, err := http.NewRequest("POST", "/schedules/sulaTopHair", strings.NewReader(string(body))); if err != nil {
		return err
	}
	c.Request = req
	handlers.CreateSchedule(c)
	if want_status != w.Code {
		return fmt.Errorf("Status Code unexpected: want: %d, received: %d", want_status, w.Code)
	}
	if want_return != w.Body.String() {
		return fmt.Errorf("Return unexpected: want: %s, received: %s", want_return, w.Body.String())
	}
	return nil
}

func Test_Create_Schedule(t *testing.T) {
	logrus.SetOutput(io.Discard)
	for _, test := range tests {
		if err := createScheduleTest(test.Body, test.WantStatus, test.WantReturn, test.Mode); err != nil {
			fmt.Printf("%s --- FAILED: %v\n", test.Name, err.Error())
			t.Fail()
			continue
		}
		fmt.Printf("%s --- PASSED\n", test.Name)
	}
}