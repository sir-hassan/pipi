# What Is Pipi?

Pipi is a simple API that parses amazon prime movie pages and returns a json representation for it. When a client requests a movie title, pipi will make a background request, fetch the
respective Amazon website, parse it and give back a valid json result to the client.

The request path should be always: http://localhost:8080/movie/amazon/{amazon_id}. when using amazon_id `B00K19SD8Q` the result will be:
```json
{
  "title" : "Um Jeden Preis [dt./OV]",
  "release_year" : 2013,
  "actors" : [
    "Dennis Quaid",
    "Zac Efron",
    "Kim Dickens"
  ],
  "poster" : "https://images-na.ssl-images-amazon.com/images/S/sgp-catalog-images/region_DE/universum-00664000-Full-Image_GalleryBackground-de-DE-1617099345129._SX1080_.jpg",
  "similar_ids" : [
    "B00IM7PHLA",
    "B0172JEA7K",
    "B08P3VVFFM",
    "B00N1DTCM0"
  ]  
}
```

