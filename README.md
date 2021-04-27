[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
# company-api

**company-api** is a simple Golang API that is still in development. The purpose of the application is to populate a database with CSV data and allow for importing (merging) and querying the data in the database. 

## Installation

Just clone the repository. and to a `go build`

## Usage

We have provided a Makefile with scripts to make common tasks easier. By default, if the `make`command is called without arguments, a `build` and `test`will be executed. Below each functionality is explained:

### clean

The `make clean` command will run a `go clean` in the directory and remove any sqlite database files that are leftover.

### build

The `make build` will build the application outputting the binary to the specified file in the `BINARY_NAME` variable on the Makefile.

### run

The `make run` will run the build command and the execute the application using the compiled executable

### test

This command will run the unit test suite for the project.

### package

This command will make a "production build" using the Dockerfile contained in this repo, meant for a production environment. The resulting image has only the compiled binary inside it.

## Interacting with the API

As previsously mentioned, the project, when starting will populate the database with the example data from the CSV located at `input_data/initial_data.csv`. After that, it is possible to:
1. Merge new data by using the endpoint `/import`
2. Query the existing data using the endpoint `/companies/search`

By default the API will run on port 8000, but setting it to a different port will be possible using `viper` to manage config files.

As of now, the insertion of new companies is not available using the API. 

### Merging new data

To add the website data into the database, a POST request with the desired CSV is needed. An example of siad request using curl can be found below:
```terminal
curl localhost:8000/import -F file=@input_data/additional_data.csv
```
The endpoint `/import`expects:
- A POST request
- With a field named `file` that contains the file to be imported. 

If the request is accepted (no missing file, no error opening it) then the application will proceed to:
1. "Download" the file locally;
2. Read its contents and format the data as desired;
3. Handle the data merging
4. Remove the file from local storage

In any case of error, either a status 400 (Bad Request) response will be provided (in errors with the request itself) or a status 500 (Server Error) will be used (in case there was data corruption or error when opening the file). If the request method is not a POST, an error 405 (Not Allowed) will be sent to the client. 

### Querying existing data

To query for existing data, a GET request to the endpoint `/companies/search` can be made. The request should contain a JSON body as follows:
```json
{
    "name": "tola",
    "zip_code": "00000"
}
```

Where name can be a fragment of the name, while the zip_code must be matching length 5 and be correct with the desired place. To check if there's a match in our database we use a `LIKE` operand when checking the name and an equals when checking the zip_code. Both fields must be present and filled for the request to succeed.

If no matches where found a 404 error will be sent, and if there was no body in the request, an error 400 will be send. If the request method is not a GET, an error 405 (Not Allowed) will be sent to the client.

### Developer notes

This project uses SQlite as "self contained" database, but interacts with the database using `gorm`, meaning that in a production scenario, given that a new "adapter" is implemented respecting the current interfaces, everything should work out of the box. Right now the project in in the process of migrating all database poperations to the same interface, but interaction using the current methods is still possible.

>A note about validation
>>For this project there was an attept of using gorm validations + govalidator to do schema validation on the requests but to do that the usage of a "fork" of gorm is needed. To not give up on stability, the schema validation because a "next-iteration" feature.

One point of interest is migrating the database querying endpoint to operate using query arguments instead of a json body. 

### Running tests
This project has unit tests, and PRs have to pass in all tests to be merged.

```bash
make test
```
Or, buy using go directly
```bash
go test -v
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Built with

- Golang [Docs](https://golang.google.cn/doc/)
- Mux Router [Github](https://github.com/gorilla/mux)
- Some ideas for automation from [golang-cookiecutter](https://github.com/lacion/cookiecutter-golang)

## License
[MIT](https://choosealicense.com/licenses/mit/)
