# POST - create a movie
POST http://localhost:4000/v1/movies
```json
{
    "title": "Test movie2",
    "genres": ["test"],
    "runtime": "100 mins",
    "year": 2016
}
```
HTTP/1.1 200
[Captures]
testMovieId: jsonpath "$.movie.id"
[Asserts]
jsonpath "$.movie.title" == "Test movie2"
jsonpath "$.movie.runtime" contains "100"
jsonpath "$.movie.year" == 2016
jsonpath "$.movie.version" == 1


# GET - fetch all Movies
GET http://localhost:4000/v1/movies
HTTP/1.1 200
[Asserts]
# jsonpath "$.movies[0].title" == "Test movie2"
# jsonpath "$.movies[0].runtime" contains "100"
# jsonpath "$.movies[0].year" == 2016
# jsonpath "$.movies[0].version" == 1



# GET - fetch movie by ID
GET http://localhost:4000/v1/movies/{{testMovieId}}
HTTP/1.1 200
[Asserts]
jsonpath "$.movie.title" == "Test movie2"
jsonpath "$.movie.runtime" contains "100"
jsonpath "$.movie.year" == 2016
jsonpath "$.movie.version" == 1


# PATCH - update single movie in place
PATCH http://localhost:4000/v1/movies/{{testMovieId}}
```json
{
    "title": "Test movie3"
}
```
HTTP/1.1 200
[Asserts]
jsonpath "$.movie.title" == "Test movie3"
jsonpath "$.movie.runtime" contains "100"
jsonpath "$.movie.year" == 2016
jsonpath "$.movie.version" == 2


# DELETE - remove a movie by id
DELETE http://localhost:4000/v1/movies/{{testMovieId}}
HTTP/1.1 200
[Asserts]
jsonpath "$.message" == "successfully deleted movie"
