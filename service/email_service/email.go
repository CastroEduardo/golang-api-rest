package email_service

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/CastroEduardo/golang-api-rest/pkg/setting"
)

type Person struct {
	Name string
}

type entry struct {
	Name string
	Done bool
}

type ToDo struct {
	User string
	List []entry
}

func SendEmailToList(email string, subject string, temp string) bool {

	var doc bytes.Buffer
	content, err := ioutil.ReadFile("html_templates/list.html")
	if err != nil {
		log.Fatal(err)
	}

	//t := template.Must(template.New("html-tmpl").ParseFiles(string(content)))
	t, err := template.New("templateName").Parse(string(content))
	if err != nil {
		panic(err)
	}

	listTodo := []entry{}
	for i := 0; i < 100; i++ {
		newItem := entry{
			Name: "Item Todo: " + strconv.Itoa(i),
			Done: false,
		}
		listTodo = append(listTodo, newItem)
	}

	todoStruc := ToDo{User: "Eduardo Castro Test List", List: listTodo}
	t.Execute(&doc, todoStruc)
	if err != nil {
		panic(err)
	}
	result := sendGmailSMTP(email, subject, doc.String())

	return result
}

func SendEmailToTest(email string, subject string, temp string) bool {

	dataStruc := Person{"Eduardo Castro TESt#1"}
	//td := Person{"Eduardo Castro"}
	var doc bytes.Buffer

	content, err := ioutil.ReadFile("html_templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("templateName").Parse(string(content))
	if err != nil {
		panic(err)
	}

	t.Execute(&doc, dataStruc)
	if err != nil {
		panic(err)
	}

	result := sendGmailSMTP(email, subject, doc.String())

	return result

}

func sendGmailSMTP(email string, sub string, template string) bool {

	emailHost := setting.EmailSetting.Host
	emailFrom := setting.EmailSetting.User
	emailPassword := setting.EmailSetting.Password
	emailTo := email //[]string{"to@example.com"} //email
	subj := sub
	body := template
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = "<" + emailFrom + ">" //from.String()
	headers["To"] = emailTo
	headers["Subject"] = subj
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html"
	headers["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	servername := "smtp.gmail.com:465"

	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", emailFrom, emailPassword, emailHost)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
		return false
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
		return false
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
		return false
	}

	// To && From
	if err = c.Mail(emailFrom); err != nil {
		log.Panic(err)
		return false
	}

	if err = c.Rcpt(emailTo); err != nil {
		log.Panic(err)
		return false
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
		return false
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
		return false
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
		return false
	}

	c.Quit()
	return true

}

func sendGmailSMTPFile() bool {
	// content, err := ioutil.ReadFile("html_templates/list.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	path, _ := os.Getwd()
	// handle err
	fmt.Println(path)
	var (
		serverAddr = setting.EmailSetting.Host
		password   = setting.EmailSetting.Password
		emailAddr  = setting.EmailSetting.User
		portNumber = 465
		tos        = []string{
			"data2354@gmail.com",
			"eduardocastro_39@hotmail.com",
		}
		cc = []string{
			setting.EmailSetting.User,
		}

		attachmentFilePath = path + "/html_templates/list.html"
		filename           = "logFileUI.txt"
		delimeter          = "**=myohmy689407924327"
	)

	log.Println("======= Test Gmail client (with attachment) =========")
	log.Println("NOTE: user need to turn on 'less secure apps' options")
	log.Println("URL:  https://myaccount.google.com/lesssecureapps\n\r")

	tlsConfig := tls.Config{
		ServerName:         serverAddr,
		InsecureSkipVerify: true,
	}

	log.Println("Establish TLS connection")
	conn, connErr := tls.Dial("tcp", fmt.Sprintf("%s:%d", serverAddr, portNumber), &tlsConfig)
	if connErr != nil {
		log.Panic(connErr)
	}
	defer conn.Close()

	log.Println("create new email client")
	client, clientErr := smtp.NewClient(conn, serverAddr)
	if clientErr != nil {
		log.Panic(clientErr)
	}
	defer client.Close()

	log.Println("setup authenticate credential")
	auth := smtp.PlainAuth("", emailAddr, password, serverAddr)

	if err := client.Auth(auth); err != nil {
		log.Panic(err)
	}

	log.Println("Start write mail content")
	log.Println("Set 'FROM'")
	if err := client.Mail(emailAddr); err != nil {
		log.Panic(err)
	}
	log.Println("Set 'TO(s)'")
	for _, to := range tos {
		if err := client.Rcpt(to); err != nil {
			log.Panic(err)
		}
	}

	writer, writerErr := client.Data()
	if writerErr != nil {
		log.Panic(writerErr)
	}

	//basic email headers
	sampleMsg := fmt.Sprintf("From: %s\r\n", emailAddr)
	sampleMsg += fmt.Sprintf("To: %s\r\n", strings.Join(tos, ";"))
	if len(cc) > 0 {
		sampleMsg += fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ";"))
	}
	sampleMsg += "Subject: Send FILE FROM GOLANG TEST\r\n"

	log.Println("Mark content to accept multiple contents")
	sampleMsg += "MIME-Version: 1.0\r\n"
	sampleMsg += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", delimeter)

	//place HTML message
	log.Println("Put HTML message")
	sampleMsg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	sampleMsg += "Content-Type: text/html; charset=\"utf-8\"\r\n"
	sampleMsg += "Content-Transfer-Encoding: 7bit\r\n"
	sampleMsg += fmt.Sprintf("\r\n%s", "<html><body><h1>Hi FROM GOLANG</h1>"+
		"<p>This FILE IS SEnd TO GOLANG TO TEST</p></body></html>\r\n")

	//place file
	log.Println("Put file attachment")
	sampleMsg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	sampleMsg += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
	sampleMsg += "Content-Transfer-Encoding: base64\r\n"
	sampleMsg += "Content-Disposition: attachment;filename=\"" + filename + "\"\r\n"
	//read file
	rawFile, fileErr := ioutil.ReadFile(attachmentFilePath)
	if fileErr != nil {
		log.Panic(fileErr)
	}
	sampleMsg += "\r\n" + base64.StdEncoding.EncodeToString(rawFile)

	//write into email client stream writter
	log.Println("Write content into client writter I/O")
	if _, err := writer.Write([]byte(sampleMsg)); err != nil {
		log.Panic(err)
	}

	if closeErr := writer.Close(); closeErr != nil {
		log.Panic(closeErr)
	}

	client.Quit()

	log.Print("done.")
	return true

}
