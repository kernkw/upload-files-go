package main
 
import (
    "html/template"
	"io"
	"net/http"
	"os"
	//"./sendemail"
	"./consuela"
	"log"
	//"./api"
	"fmt"
  //"mime"
  "github.com/dchest/uniuri"
  "github.com/gorilla/sessions"
  "github.com/gorilla/context"

)

var token = uniuri.New()
var store = sessions.NewCookieStore([]byte(token))
//Compile templates on start
var templates = template.Must(template.ParseFiles("tmpl/auth.html","tmpl/email.html", "tmpl/upload.html"))
 
//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

func credentialHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    //GET displays the upload form.
    case "GET":
      display(w, "auth", nil)
   
    //POST takes the uploaded file(s) and saves it to disk.
    case "POST":
        fmt.Println(r)
        
        // pass list to consuela to remove bad addresses
        username := r.FormValue("username")
        password := r.FormValue("password")

        uri := "https://api.sendgrid.com/api/bounces.count.json?api_user=" + username + "&api_key=" + password
        response, err  := http.StatusText(uri)
        if err != nil {
        log.Fatal(err)
        }
        fmt.Println(resp)

        // existing session: Get() always returns a session, even if empty.
        session, _ := store.Get(r, "session")

        // Set some session values.
        session.Values["username"] = username
        session.Values["password"] = password

        // Save it.
        session.Save(r, w)
      
        //Check if correct creds

        // set unique token
        w.Header().Set("StatusCode", "202")
        w.Header().Set("Content Type", "application/json")
        w.Header().Set("Body", "“message”: “success”")
        w.Header().Set("Body", token)
        fmt.Println(w)
        //display success message.
        display(w, "auth", "Authenticated!")

    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
  }
}
func emailHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    //GET displays the upload form.
  case "GET":
    display(w, "email", nil)
 
  //POST takes the uploaded file(s) and saves it to disk.
  case "POST":
      // existing session: Get() always returns a session, even if empty.
      session, _ := store.Get(r, "session")
      // pass list to consuela to remove bad addresses
      email := r.FormValue("email")
      emailToken := r.FormValue("token")

      if emailToken == token {
        session.Values["email"] = email
        w.Header().Set("StatusCode", "202")
        w.Header().Set("Content Type", "application/json")
        w.Header().Set("Body", "“message”: “success”")
      } else {
        w.Header().Set("StatusCode", "401")
        w.Header().Set("Content Type", "application/json")
      }
     
      fmt.Println("email: " + email)
      //display success message.
      display(w, "email", "Email added successfully.")

  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
  }
}
 
//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
  	//GET displays the upload form.
  	case "GET":
  		display(w, "upload", nil)
   
  	//POST takes the uploaded file(s) and saves it to disk.
  	case "POST":
      // existing session: Get() always returns a session, even if empty.
      session, _ := store.Get(r, "session")

  		//get the multipart reader for the request.
  		reader, err := r.MultipartReader()

  		if err != nil {
  			http.Error(w, err.Error(), http.StatusInternalServerError)
  			return
  		}

   
  		//copy each part to destination.
  		for {
    			part, err := reader.NextPart()
    			if err == io.EOF {
    				break
    			}

         // fmt.Println(form.Value("username"))
     
    			//if part.FileName() is empty, skip this iteration.
    			if part.FileName() == "" {
    				continue
    			}
    			dst, err := os.Create("/Users/sendgrid1/upload-files-go/" + part.FileName())
    			defer dst.Close()
     
    			if err != nil {
    				http.Error(w, err.Error(), http.StatusInternalServerError)
    				return
    			}
    			
    			if _, err := io.Copy(dst, part); err != nil {
    				http.Error(w, err.Error(), http.StatusInternalServerError)
    				return
    			}
    			// pass list to consuela to remove bad addresses
    			inputfile, err := os.Open("/Users/sendgrid1/upload-files-go/" + part.FileName())
    				if err != nil {
    				log.Fatal(err)
    					}
    			defer inputfile.Close()

    			update.Consuela(inputfile, session.Values["username"].(string), session.Values["password"].(string))
    		  //unsubscribe.Add(dst, username, password)

    		  //send output link in email using SendGrid
    			// emailRecipient := r.FormValue("email")
    			// fmt.Println(emailRecipient)
    			// email.SendEmail("kyle.w.kern@gmail.com")
    		}
    		//display success message.
    		display(w, "upload", "Upload successful.")

    default:
    		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

 
func main() {
  http.HandleFunc("/auth", credentialHandler)
  http.HandleFunc("/email", emailHandler)
  http.HandleFunc("/upload", uploadHandler)
	//static file handler.
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
 
	//Listen on port 8080
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}