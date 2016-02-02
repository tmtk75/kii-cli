package kiicli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/tmtk75/cli"

	"golang.org/x/net/websocket"
)

type RawLog map[string]interface{}

func (self *RawLog) Log() *Log {
	f := time.RFC3339Nano // "2006-01-02T15:04:05.999Z"
	t, err := time.Parse(f, (*self)["time"].(string))
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &Log{
		(*self)["key"].(string),
		(*self)["level"].(string),
		(*self)["description"].(string),
		t,
	}
}

type Log struct {
	Key         string
	Level       string
	Description string
	Time        time.Time
}

type AuthRequest struct {
	AppID        string `json:"appID"`
	AppKey       string `json:"appKey"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	Token        string `json:"token"`
	Command      string `json:"command"` // 'tail' or 'cat'
	UserID       string `json:"userID"`
	Level        string `json:"level"`
	Limit        int    `json:"limit"`
	DateFrom     string `json:"dateFrom"`
	DateTo       string `json:"dateTo"`
}

func (self *GlobalConfig) AuthRequest() *AuthRequest {
	req := &AuthRequest{
		AppID:        self.AppId,
		AppKey:       self.AppKey,
		ClientID:     self.ClientId,
		ClientSecret: self.ClientSecret,
		Token:        self.Token,
	}
	return req
}

func (s *AuthRequest) Parse(c *cli.Context) {
	//
	s.Command = "cat"
	if c.Bool("tail") {
		s.Command = "tail" //cat
	}
	//t, _ := time.Parse("2006-01-02 15:04:05", "2012-01-01 12:12:12")
	//return t.Format("2006-01-02 15:04:05")

	s.Limit = c.Int("num")
	s.UserID = c.String("user-id")
	s.Level = strings.ToUpper(c.String("level"))

	if c.String("date-from") != "" {
		s.DateFrom = c.String("date-from") //"2015-01-08:07:40:00",
	}
	if c.String("date-to") != "" {
		s.DateTo = c.String("date-to") //"2015-01-09:00:00:00",
	}
}

func StartLogging(c *cli.Context) {
	p := Profile()
	req := p.AuthRequest()
	req.Parse(c)

	j, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	//fmt.Printf("%s", string(j))

	url := p.EndpointUrlForApiLog()
	logger.Printf("log-url: %s", url)
	logger.Printf("%s", j)
	ws, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		log.Fatalf("%v", err)
	}
	_, err = ws.Write(j)
	if err != nil {
		log.Fatalf("%v", err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	go func() {
		_ = <-sigc
		ws.Close()
		os.Exit(0)
	}()

	rch := make(chan []RawLog)

	go func() {
		for {
			//log.Println("will read")
			msg := <-rch
			for i, m := range msg {
				m.Print(i)
			}
			if req.Command == "cat" {
				ws.Close()
				os.Exit(0)
			}
		}
	}()

	for {
		select {
		default:
			var msg []RawLog
			err = websocket.JSON.Receive(ws, &msg)
			if err == io.EOF {
				os.Exit(0)
			}
			if err != nil {
				log.Fatalf("%v", err)
			}
			rch <- msg
			//log.Printf("wrote %d", len(msg))
		}
	}

	//TODO: gracefully exit when cat mode
}

type RawFormat map[string]string
type Format map[string]*template.Template

var format Format

func LoadFormat(path string) Format {
	e, err := exists(path)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if !e {
		logger.Printf("%v is missing", path)
		return make(Format)
	}

	body, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("%v", err)
	}

	f := make(RawFormat)
	if err := json.Unmarshal(body, &f); err != nil {
		log.Fatalf("%v", err)
	}

	r := make(Format)
	for k, v := range f {
		s := convertLogFormat(v)
		t, err := template.New(k).Parse(s)
		if err != nil {
			log.Fatalf("%v", err)
		}
		r[k] = t
	}

	return r
}

func convertLogFormat(f string) string {
	re, _ := regexp.Compile("\\${[a-zA-Z-_]+}")
	k := re.ReplaceAllFunc([]byte(f), func(a []byte) []byte {
		s := norm(string(a[2 : len(a)-1]))
		return []byte(fmt.Sprintf("{{%v}}", s))
	})
	return string(k)
}

func norm(k string) string {
	if strings.Index(k, "-") > 0 {
		return fmt.Sprintf(`index . "%v"`, k)
	}
	return fmt.Sprintf(`.%v`, k)
}

func (m *RawLog) Print(idx int) {
	key := (*m)["key"].(string)
	f := format[key]
	if f == nil {
		fmt.Printf("%v\n", *m.Log())
		return
	}

	w := bytes.NewBuffer([]byte{})
	t, _ := timeFromStringInUTC((*m)["time"].(string))
	(*m)["time"] = t.Format(time.RFC3339Nano)
	f.Execute(w, *m)
	fmt.Printf("%v\n", w)
}

var LogCommands = []cli.Command{
	{
		Name:  "log",
		Usage: "Disply logs for an app",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "level", Usage: "Filtering with level. e.g) DEBUG, INFO, ERROR"},
			cli.IntFlag{Name: "num,n", Value: 100, Usage: "Show specified number of lines"},
			cli.StringFlag{Name: "user-id", Usage: "Filtering with UserID"},
			cli.BoolFlag{Name: "tail,t", Usage: "Similar to tail -f"},
			cli.StringFlag{Name: "format-file", Usage: "File path to a format file", Value: (func() string {
				d, _ := homedir.Dir()
				return fmt.Sprintf("%v/.kii/format.json", d)
			})(),
			},
			cli.StringFlag{Name: "date-from", Usage: "Filtering from specified date"},
			cli.StringFlag{Name: "date-to", Usage: "Filtering until specified date"},
		},
		Action: func(c *cli.Context) {
			format = LoadFormat(c.String("format-file"))
			StartLogging(c)
		},
	},
}
