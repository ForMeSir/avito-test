package integration

import (
	"avito/internal/banner"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const(
 	query="INSERT INTO banners (tag_ids,feature_id,content_title,content_text,content_url,is_active,created_at,updated_at) values($1,$2,$3,$4,$5,$6,$7,$8)"
 )

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

func TestOneUserBanner(t *testing.T, router *gin.Engine) {
	req, err := http.NewRequest("GET", "/user_banner", nil)
	if err != nil {
		t.Fatal(err)
	}
	testHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == 400

		p, err := io.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "Некорректные данные") > 0
		return statusOK && pageOK
	})
}

func TestTwoUserBanner(t *testing.T, router *gin.Engine) {
	req, err := http.NewRequest("GET", "/user_banner", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("tag_id", "16")
	q.Add("feature_id", "322")
	req.URL.RawQuery = q.Encode()


	testHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == 401

		p, err := io.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "Пользователь не авторизован") > 0

		return statusOK && pageOK
	})

}

func TestThreeUserBanner(t *testing.T, router *gin.Engine) {
	req, err := http.NewRequest("GET", "/user_banner", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("token", "user_token")

	q := req.URL.Query()
	q.Add("tag_id", "16")
	q.Add("feature_id", "322")
	req.URL.RawQuery = q.Encode()

	testHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == 404

		p, err := io.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "Баннер не найден") > 0
		return statusOK && pageOK
	})
}

func TestFourUserBanner(t *testing.T, router *gin.Engine, db *sqlx.DB) {
	
	var bann banner.Banner=banner.Banner{TagIds:[]int{17},FeatureId: 32}

	bann.Content = banner.Content{Title:"some_tile", Text: "some_tet",Url: "some_ul",IsActive: false}
	
	 _,err:=db.Exec(query, pq.Array(bann.TagIds),bann.FeatureId,bann.Content.Title,bann.Content.Text,bann.Content.Url,bann.Content.IsActive,time.Now().Format(time.RFC3339),time.Now().Format(time.RFC3339))
	 if err!=nil{
		logrus.Info(err)
		 t.Fail()
	 }
	
	 req, err := http.NewRequest("GET","/user_banner", nil)
					if err != nil {
						t.Fatal(err)
					}
					req.Header.Add("token","user_token")
	
					q:= req.URL.Query()
									q.Add("tag_id", "17")
									q.Add("feature_id", "32")
									req.URL.RawQuery = q.Encode()
	
	testHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == 403
	
		 p, err := io.ReadAll(w.Body)
		 pageOK := err == nil && strings.Index(string(p), "Пользователь не имеет доступа") > 0
		 return statusOK && pageOK
	})
	
}


func TestFiveUserBanner(t *testing.T, router *gin.Engine, db *sqlx.DB) {
	
	var bann banner.Banner=banner.Banner{TagIds:[]int{16},FeatureId: 342}

	bann.Content = banner.Content{Title:"some_tile", Text: "some_tet",Url: "some_ul",IsActive: false}
	
	 _,err:=db.Exec(query, pq.Array(bann.TagIds),bann.FeatureId,bann.Content.Title,bann.Content.Text,bann.Content.Url,bann.Content.IsActive,time.Now().Format(time.RFC3339),time.Now().Format(time.RFC3339))
	 if err!=nil{
		logrus.Info(err)
		 t.Fail()
	 }
	
	 req, err := http.NewRequest("GET","/user_banner", nil)
					if err != nil {
						t.Fatal(err)
					}
					req.Header.Add("token","admin_token")
	
					q:= req.URL.Query()
									q.Add("tag_id", "16")
									q.Add("feature_id", "342")
									req.URL.RawQuery = q.Encode()
	
	testHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == 200
	
		 p, err := io.ReadAll(w.Body)
		 pageOK := err == nil && strings.Index(string(p), "Баннер пользователя") > 0
		 return statusOK && pageOK
	})
	
}


