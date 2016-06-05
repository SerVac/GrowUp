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
// random string for oauth2 API calls to protect against CSRF
	oauthStateString = "thisshouldberandom"
)


// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(w http.ResponseWriter, r *http.Request) {
	_cacheFilePath, err := tokenCacheFilePath()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}

	_oauthToken, err = tokenFromFile(_cacheFilePath)
	if err != nil {
		authURL := getAuthTokenURL()
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
		//body := r.Body
	}else {
		_oauthClient = _oauthConfig.Client(_ctx, _oauthToken)
		redirectToMainPage(w,r)

	}
}

func getClient1(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFilePath, err := tokenCacheFilePath()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFilePath)
	if err != nil {
		//tok = getTokenFromWeb(config)
		//saveToken(cacheFilePath, tok)
		//generateTokenFromWeb(config)
	}
	return config.Client(ctx, tok)
}

func getAuthTokenURL() string {
	authURL := _oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	//_oauthURL := &config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	//fmt.Printf("Go to the following link in your browser then type the " +
	//"authorization code: \n%v\n", _oauthURL)

	//_oauthURL = authURL

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
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
	fmt.Println("get value by pointer = ", &_oauthConfig)
	fmt.Println("get refferense on pointer = ", _oauthConfig)
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
	if (_oauthConfig != nil) {
		getClient(w, r)
		//client := getClient(_ctx, _oauthConfig)
	}
}
// /aouthCallnack
func handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	fmt.Println("state = ", state)
	code := r.FormValue("code")
	fmt.Println("code = ", code)

	_oauthToken, err := _oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	saveToken(_cacheFilePath, _oauthToken)

}
// /grow_up
func handleGrowUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	var htmlTemlate string = `<html><body>
	Grow UP!
	</body></html>`

	w.Write([]byte(htmlTemlate))
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

func main1() {
	/*
		ctx := context.Background()

		absPath, _ := filepath.Abs("../key/client_secret.json")
		b, err := ioutil.ReadFile(absPath)
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}

		// If modifying these scopes, delete your previously saved credentials
		// at ~/.credentials/drive-go-quickstart.json
		config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		client := getClient(ctx, config)

		srv, err := drive.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve drive Client %v", err)
		}

		r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve files.", err)
		}

		fmt.Println("Files:")
		if len(r.Files) > 0 {
			for _, i := range r.Files {
				fmt.Printf("%s (%s)\n", i.Name, i.Id)
			}
		} else {
			fmt.Print("No files found.")
		}
	*/
}