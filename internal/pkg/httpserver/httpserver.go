package httpserver

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Server struct {
	*http.Server
	window         window
	ServerHitCount int
}

type HandlerFunc = http.HandlerFunc

func NewServer(addr string) *Server {
	windowSize, err := strconv.Atoi(os.Getenv("SERVER_WINDOW_SIZE_IN_SECONDS"))
	if err != nil {
		windowSize = 10
	}
	srv := &Server{
		Server: &http.Server{Addr: addr},
		window: window{
			SLIDING_WINDOW_SIZE: windowSize,
			start:               -1,
			end:                 -1,
			running:             false,
		},
	}
	return srv
}

// Public Members
func (s *Server) Handle(path string, countServerHits bool, handler HandlerFunc) {
	handler = http.HandlerFunc(handler)
	if countServerHits {
		http.Handle(path, s.countHitsMiddleware(handler))
	} else {
		http.Handle(path, handler)
	}

}
func (s *Server) Cleanup() {
	rootPath := os.Getenv("PROJECT_ROOT")
	errH := s.writeHitsFile(rootPath)
	errW := s.writeWindowmetadataFile(rootPath)
	if errH != nil || errW != nil {
		//Delete both files
		os.Remove(filepath.Join(rootPath, ".hits"))
		os.Remove(filepath.Join(rootPath, ".windowmetadata"))
	}

}

// Private Members
func (s *Server) countHitsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.countHits() //Sets s.ServerHitsCount which is the no. of hits in the last n seconds; n is s.window.SLIDING_WINDOW_SIZE
		w.Write([]byte("No. of Hits to the server in the last " + strconv.Itoa(s.window.SLIDING_WINDOW_SIZE) + " seconds: " + strconv.Itoa(s.ServerHitCount) + "\n"))
		next.ServeHTTP(w, r)
	})
}

func (s *Server) countHits() {
	rootPath := os.Getenv("PROJECT_ROOT")
	hitTime := int64(0)
	if rootPath == "" {
		log.Fatalf("bad argument for env variable PROJECT_ROOT")
	}
	hitsFilePath := filepath.Join(rootPath, ".hits")                     //persists the timestamp in time.Unix() format of every hit in comma-separated format
	windowMetadataFilePath := filepath.Join(rootPath, ".windowmetadata") //persists the start of the sliding window, end of the sliding window and hit-counts when server kill/terminates
	if !s.window.running {                                               //checking if it is the first hit to the server instance
		// if it is first hit, checking for persisted files to load the structs for continuity
		_, err := os.Stat(hitsFilePath)
		if os.IsNotExist(err) {
			//initialise when the files don't exist
			s.initialiseWindow()
			return
		} else {
			//load from file
			s.loadWindowFromFile(hitsFilePath, windowMetadataFilePath)
		}

	}
	s.window.running = true
	hitTime = s.window.appendNewHit()

	//Extend window at the back by 1 such that end pointer always points to the currhit-1 (skipping the current hit because `previous n seconds` and not including current second)
	timeAtCurrEnd := s.window.moveEnd()
	diff := int(hitTime - timeAtCurrEnd)
	if s.window.SLIDING_WINDOW_SIZE < diff { //When there are no hits that are within the last n seconds - HitCount resets to 0
		s.window.hitCountInWindow = 0
		s.window.start = s.window.end + 1
		s.ServerHitCount = s.window.hitCountInWindow

		return
	}
	timeAtCurrStart := s.window.hitTimeRecords[s.window.start]
	//Shrink window from front until the start pointer points to a hit that happened within the last n seconds
	for {
		diff := int(hitTime - timeAtCurrStart)
		if s.window.start == s.window.end || s.window.SLIDING_WINDOW_SIZE >= diff {
			s.ServerHitCount = s.window.hitCountInWindow

			return
		}
		timeAtCurrStart = s.window.moveStart()
	}

}

func (s *Server) initialiseWindow() {
	s.window.end = -1
	s.window.start = s.window.end + 1
	s.window.hitCountInWindow = 0
	s.window.hitTimeRecords = make([]int64, 0)
	s.window.appendNewHit()
	s.ServerHitCount = s.window.hitCountInWindow
	s.window.running = true
}

func (s *Server) loadWindowFromFile(hitsFilePath, windowMetadataFilePath string) {

	hitRecordsByteArray, err := os.ReadFile(hitsFilePath)
	if err != nil {
		log.Fatalf("failed starting server: cannot load hits")
	}

	hitTimesSerialisedString := string(hitRecordsByteArray)
	s.window.hitTimeRecordsSerialised = hitTimesSerialisedString
	tmpHitTimes := strings.Split(hitTimesSerialisedString, ",")
	for _, v := range tmpHitTimes {
		hitTime, _ := strconv.Atoi(v)
		s.window.hitTimeRecords = append(s.window.hitTimeRecords, int64(hitTime))

	}
	windowmetdataByteArray, err := os.ReadFile(windowMetadataFilePath)
	if err != nil {
		log.Fatalf("failed starting server: cannot load metadata")
	}

	splitMetadataString := strings.Split(string(windowmetdataByteArray), ",")
	s.window.start, _ = strconv.Atoi(splitMetadataString[0])
	s.window.end, _ = strconv.Atoi(splitMetadataString[1])
	s.window.hitCountInWindow, _ = strconv.Atoi(splitMetadataString[2])

}

func (s *Server) writeHitsFile(rootPath string) error {
	filename := filepath.Join(rootPath, ".hits")
	b := []byte(s.window.hitTimeRecordsSerialised)
	os.WriteFile(filename, b, 0666)
	return nil

}

func (s *Server) writeWindowmetadataFile(rootPath string) error {
	filename := filepath.Join(rootPath, ".windowmetadata")
	b := []byte(strconv.Itoa(s.window.start) + "," + strconv.Itoa(s.window.end) + "," + strconv.Itoa(s.window.hitCountInWindow))
	os.WriteFile(filename, b, 0666)
	return nil
}
