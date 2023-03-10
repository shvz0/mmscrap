# mmscrap
Mass media scrapper with ability to detect author of article based on 3 methods of stylometry.
## Setup
1. Install dependencies;
```bash
go mod tidy
go mod vendor
```
2. Setup Postgresql using docker-compose: 
```bash
docker-compose up
```
3. Fill up **.env** file by example (**.env** must be placed in repository directory);
4. Build app:
```
go build main.go
```
5. Run migration
```
./main --migrate
```
## Usage
1. Run server using flag **--serve**
2. Parse todays articles using flag **--parse**
3. Migrate DB using flag **--migrate**
## API
* POST /delta - calculate delta coefficient from all authors in DB. Body params: **text**.
Most accurate method 
* POST /chisquared - calculate chi squared coefficent from all authors in DB. Body params: **text**.
In the main less accurate than delta method
* POST /mendenhall - calculate Mendenhall coefficient (using word length distribution) from all authors in DB. Body params: **text**.
Simple, old and inaccurate