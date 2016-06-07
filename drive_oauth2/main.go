package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const htmlIndex = `<html><body>
Logged in with <a href="/login">Index page code</a>
</body></html>
`

var (
	_cacheFilePath string
	_ctx = context.Background()
	_oauthConfig *oauth2.Config
	_oauthToken *oauth2.Token
	_oauthClient *http.Client
	_googleDriveService *drive.Service
// random string for oauth2 API calls to protect against CSRF
	oauthStateString = "thisshouldberandom"
)


// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func createClient(w http.ResponseWriter, r *http.Request) {
	if (_oauthConfig != nil) {
		var err error = nil
		_cacheFilePath, err = tokenCacheFilePath()
		if err != nil {
			log.Fatalf("Unable to get path to cached credential file. %v", err)
		}

		_oauthToken, err = tokenFromFile(_cacheFilePath)
		if err != nil {
			log.Print("Generate new token")
			authURL := getAuthTokenURL()
			http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
			//body := r.Body
		}else {
			log.Print("Get cached token")

			_oauthClient = _oauthConfig.Client(_ctx, _oauthToken)

			token_acces := _oauthToken.AccessToken
			//resp, err := _oauthClient.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token="+token_acces)
			//resp1, err1 := _oauthClient.Head("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token="+token_acces)

			var target Target
			err := getJson("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=" + token_acces, target)

			print(err)
			//print(resp, err, resp1, err1)


			redirectToMainPage(w, r)
		}
	}else {
		log.Fatal("Unable to get oauth configuartion!")
	}
}

type Target struct {
	audience string
	user_id string
	scope string
	expires_in string
	error string
	error_description string
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func getAuthTokenURL() string {
	authURL := _oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("drive-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func auth() {
	absPath, _ := filepath.Abs("../key/client_secret.json")
	b, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	_oauthConfig, err = google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	//fmt.Println("get value by pointer = ", &_oauthConfig)
	//fmt.Println("get refferense on pointer = ", _oauthConfig)
	//var _ = _oauthConfig
}

// ?
func redirectToMainPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/grow_up", http.StatusTemporaryRedirect)
}

// http://127.0.0.1:7000/
func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

// /login
func handleLogin(w http.ResponseWriter, r *http.Request) {
	createClient(w, r)
}

// /aouthCallnack
func handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	fmt.Println("state = ", state)
	code := r.FormValue("code")
	fmt.Println("code = ", code)

	//TODO state error check
	// converts an authorization code into a token
	//_oauthToken, err := _oauthConfig.Exchange(oauth2.NoContext, code)
	_oauthToken, err := _oauthConfig.Exchange(_ctx, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	saveToken(_cacheFilePath, _oauthToken)

	redirectToMainPage(w, r)
}

// /grow_up
func handleGrowUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	//var htmlTemlate string = `<html><body>`
	//Grow UP!
	//</body></html>`

	var err error = nil
	_googleDriveService, err = drive.New(_oauthClient)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	reader, err := _googleDriveService.Files.List().PageSize(10).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files.", err)
	}

	htmlBody := "Files: \n"
	fmt.Println("Files:")
	if len(reader.Files) > 0 {
		for _, i := range reader.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
			htmlBody = htmlBody + i.Name + "(" + i.Id + ")\n"
		}
	} else {
		htmlBody = htmlBody + "No files found"
		fmt.Print("No files found.")
	}

	htmlStart := "<html><body>"
	htmlEnd := "</body></html>"
	htmlText := htmlStart + htmlBody + htmlEnd
	w.Write([]byte(htmlText))
	//w.Write([]byte(htmlTemlate))
}

func main() {
	auth()

	// server
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/oauth_callback", handleOAuthCallback)
	http.HandleFunc("/grow_up", handleGrowUp)

	fmt.Print("Started running on http://127.0.0.1:7000\n")
	fmt.Println(http.ListenAndServe(":7000", nil))

}
