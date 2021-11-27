package main



import (

    "encoding/csv"

    "os"

    "fmt"

    "math/rand"

    "time"

    "log"

    "net/http"

    "encoding/json"

    "strconv"

    "github.com/gorilla/mux"

    "io/ioutil"

    "image"

    _ "image/gif"

    _ "image/jpeg"

    _ "image/png"

)



type Store struct {

    

    AreaCode     int    `json:"areaCode"`

    StoreName    string `json:"storeName"`

    StoreID      string `json:"storeId"`

}



type JobResults struct {

    JobResults  [] Result `json:"jobResults"`

}



type Result struct {


    JobId     int    `json:"jobid"`

    ImageResolution   [] int `json:"imageResolution"`

}



type JobData struct {

    StoreID       string    `json:"store_id"`

    ImageUrl     []string `json:"image_url"`

    VisitTime     string     `json:"visit_time"`

}



type Job struct {

    Count        int  `json:"count"`

    Visits     [] JobData `json:"visits"`

}



func createStoreList(data [][]string) []Store {



    var storeList []Store

    for i, line := range data {

        if i > 0 {

            var rec Store

            for j, field := range line {

                if j == 0 {

                    var err error

                    rec.AreaCode, err = strconv.Atoi(field)

                    if err != nil {

                        continue

                    }

                } else if j == 1 {

                    rec.StoreName = field

                } else if j == 2 {

                    rec.StoreID = field

                }

            }

            storeList = append(storeList, rec)

        }

    }

    return storeList

}



func readCsvFile(filePath string) [][]string {

    f, err := os.Open(filePath)

    if err != nil {

        log.Fatal("Unable to read input file " + filePath, err)

    }

    defer f.Close()



    csvReader := csv.NewReader(f)

    records, err := csvReader.ReadAll()

    if err != nil {

        log.Fatal("Unable to parse file as CSV for " + filePath, err)

    }



    return records

}



func homePage(w http.ResponseWriter, r *http.Request){

    fmt.Fprintf(w, "Welcome to the HomePage!")

    fmt.Println("Endpoint Hit: homePage")

}



func processImage(uri string) int {

    resp, err := http.Get(uri)

    if err != nil {

        log.Fatal(err)

    }

    defer resp.Body.Close()



    m, _, err := image.Decode(resp.Body)

    if err != nil {

        log.Fatal(err)

    }

    g := m.Bounds()



    height := g.Dy()

    width := g.Dx()



    resolution := 2*(height + width)



    return resolution;

}



func returnAllStores(w http.ResponseWriter, r *http.Request){

    fmt.Println("Endpoint Hit: returnAllArticles")

    records := readCsvFile("./StoreMasterAssignment.csv")

    storeList := createStoreList(records)

    jsonData, err := json.MarshalIndent(storeList, "", " ")

    if err != nil {

        log.Fatal(err)

    }

    fmt.Fprintf(w, "%+v", string(jsonData))

}



func createNewJob(w http.ResponseWriter, r *http.Request) {

    // get the body of our POST request

    // return the string response containing the request body    

    reqBody, _ := ioutil.ReadAll(r.Body)

    var job Job

    var imgRes []int

    json.Unmarshal(reqBody, &job)

    for i := 0; i < job.Count; i++ {

        imageCount := len(job.Visits[i].ImageUrl)

        for j := 0; j < imageCount; j++ {

            rand.Seed(time.Now().UnixNano())

            x := rand.Intn(400)

            time.Sleep(time.Duration(x) * time.Millisecond)

            var resolution int = processImage(job.Visits[i].ImageUrl[j])

	    /*var temp int = job.Visits[i].StoreID[j]
	    var f int = 0
            for k := 0; k < len(storeList); k++ {
            	if temp == storeList[k] {
                    f = 0
                    break
                }
                else {
                    f = 1
                }
		
            }*/
	    
	    //if f == 0 {

                imgRes  = append(imgRes , resolution)
	    //}
            //else {
                //-1 for error :- ID not present
            //    imgRes  = append(imgRes , -1)
            //}


        }

    }



    filename := "result.json"



    file, err := ioutil.ReadFile(filename)

    if err != nil {

        fmt.Println(err)

    }



    data := []Result{}



    json.Unmarshal(file, &data)



    newStruct := &Result {

        JobId : rand.Intn(100000),

        ImageResolution : imgRes, 

    }



    data = append(data, *newStruct)



    // Preparing the data to be marshalled and written.

    dataBytes, err := json.Marshal(data)

    if err != nil {

        fmt.Println(err)

    }



    err = ioutil.WriteFile(filename, dataBytes, 0644)

    if err != nil {

        fmt.Println(err)

    }

    json.NewEncoder(w).Encode("Job Created Successfully")

}



func returnSingleJob(w http.ResponseWriter, r *http.Request){

    query := r.URL.Query()

    jobId, err := strconv.Atoi(query.Get("jobid"))

    if err != nil {

        fmt.Println(err)

    }



    filename := "result.json"



    file, err := ioutil.ReadFile(filename)

    if err != nil {

        fmt.Println(err)

    }



    data := []Result{}



    json.Unmarshal(file, &data)

    for i := 0; i < len(data); i++ {

       if(data[i].JobId == jobId) {

            json.NewEncoder(w).Encode(data[i])

            return;

       }

    }

    json.NewEncoder(w).Encode("NO Found")

}



func handleRequests() {

    myRouter := mux.NewRouter().StrictSlash(true)



    myRouter.HandleFunc("/", homePage)

    myRouter.HandleFunc("/stores", returnAllStores)

    myRouter.HandleFunc("/job", createNewJob).Methods("POST")

    myRouter.HandleFunc("/status", returnSingleJob)

    log.Fatal(http.ListenAndServe(":8000", myRouter))

}



func main() {

    handleRequests()

}