# movie_generator
Generates a random movie in Elsa &amp; Rene recommendations

All you have to do is press generate and it will give you a random movie.

You can also add or delete movies from the list using the secret endpoint.

The secret endpoint is:
```
/secret/{token}/{action}/{id}
```

Where:
- `{token}` is the secret token
- `{action}` is the action to perform, which is either `add` or `delete`
- `{id}` is the id of the IMDB movie id to add or delete from the list

Example:

```
/secret/secret-token/add/tt12345678
```
