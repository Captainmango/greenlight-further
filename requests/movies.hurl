# GET - fetch all Movies
GET http://localhost:4000/v1/movies
HTTP/1.1 200

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

# DELETE - remove a movie by id
DELETE http://localhost:4000/v1/movies/{{testMovieId}}
