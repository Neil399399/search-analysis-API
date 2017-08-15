# Seearch-Analysis-API

## The API about google place search and Analysis by GoJieba

---
## Should Downloads Library
* **[Google Maps](https://github.com/googlemaps/google-maps-services-go)** 
>  go get googlemaps.github.io/maps

* **[Bleve](https://github.com/blevesearch/bleve)** 
>  go get github.com/blevesearch/bleve

* **[GoJieba](https://github.com/yanyiwu/gojieba)** 
>  go get github.com/yanyiwu/gojieba
---
## Attention
* Need use gcc (or run in linux)
* Want use search function,should input location and keyword.
* Want use Analysis function,should input query word.
* Search function radius is 500m.
* If want run server, install and run in root directory ; if want run client , build and run in http_client directory.
* Change your domain in /static / routes.js.
* Can not search to many time in moment,because google apikey has restriction. 