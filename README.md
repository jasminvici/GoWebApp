Web App using Go Programming Language and MySQL

Simple web app that allows you to perform CRUD operations over a user database.

Prerequisites: 
1) Go programming language
2) MySQL

Installation:
1) Clone this repository to your local machine
2) Move to the root directory of a project
3) Open terminal and run following command: go get -u github.com/go-sql-driver/mysql

Configuring Go with MySQL:
1) Create a database in MySQL
2) In "main.go", update the following line with your corresponding data (change 'root','password' and 'webappdb' with your MySQL username, password and
database name)

Running the App:
1) Open terminal and move to the root directory of a project
2) Run following command to start the app: go run main.go
3) The app should now be running on "http://localhost:8080
