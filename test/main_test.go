package test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"ip.limit.rate/internal/api/rest"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	wg           sync.WaitGroup
	counter      = make(map[string]int)
	counterMutex sync.Mutex
)

func TestPingRouter(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	rand.Seed(time.Now().UnixNano())
	num := 3

	router := rest.NewRouter().GetRouter()
	t.Logf("number of ip addresses: %d", num)

	wg.Add(num)
	for i := 0; i < num; i++ {
		go DoPing(t, router, &wg, fmt.Sprintf("%d.%d.%d.%d", rand.Intn(257), rand.Intn(257), rand.Intn(257), rand.Intn(257)))
	}
	wg.Wait()

	keys := make([]string, 0, len(counter))
	for key := range counter {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		t.Logf("%s:%d", key, counter[key])
	}
}

func DoPing(t *testing.T, router *gin.Engine, wg *sync.WaitGroup, ip string) {
	for i := 1; i < 81; i++ {
		recorder := httptest.NewRecorder() // 取得 ResponseRecorder 物件
		req, _ := http.NewRequest("GET", "/ping", nil)
		req.RemoteAddr = ip
		router.ServeHTTP(recorder, req)

		t.Logf("ip:%s, body:%s", ip, recorder.Body.String())
		if i > 60 {
			assert.Equal(t, http.StatusTooManyRequests, recorder.Code)
			assert.Contains(t, recorder.Body.String(), "Error")
			AddCounter(ip, true)
		} else {
			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Contains(t, recorder.Body.String(), strconv.Itoa(i))
			AddCounter(ip, false)
		}
	}
	wg.Done()
}

func AddCounter(ip string, isError bool) {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	str := fmt.Sprintf("%s[%t]", ip, !isError)
	counter[str] ++
}
