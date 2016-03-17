```
           _.-,=_"""--,_
        .-" =/7"   _  .3#"=.
      ,#7  " "  ,//)#d#######=.
    ,/ "      # ,i-/###########=
   /         _)#sm###=#=# #######\
  /         (#/"_`;\//#=#\-#######\
 /         ,d####-_.._.)##P########\
,        ,"############\\##bi- `\| Y.
|       .d##############b\##P'   V  |
|\      '#################!",       |
|C.       \###=############7        |
'###.           )#########/         '
 \#(             \#######|         /
  \B             /#######7 /      /
   \             \######" /"     /
    `.            \###7'       ,'
      "-_          `"'      ,-'
         "-._           _.-"
             """"---""""

```

GLOBO
===

Globo is a microservice that makes it easy  to convert from lat/long  to google's s2 formats.


## Endpoints

```
tos2/point
```

Takes a json body that looks like this:

```json
{
	"lat":55.00,
	"lng":55.00
}
```
And returns

```json
{
    "cellid": 4890663957615284359
}
```

If the optional field "precision" is  specified then the parent cellid with
level == precision is returned.



```
tos2/geojson/:type
```
Converts a geojson :type to a collection of s2 cellIds (with the request precision/ maxnumber) + bounding box, and returns a geojson feature collection that can be displayed

Values accepted for :type: :

 - point
 - polygon
 - multipolygon
 

Query parameters:
 - precision integer  specifies how precise the generated cellids should be, default max precision.


Note:
  if the any of the polygons are not CCW winded, we refuse to convert them.
i.e. it's the client responsablity to send proper json.

# install and build

```sh
export PORT=3001 #or wathever
go build && ./globo 
```

Dependecies are vendored thus requre GOVENDOREXPERIMENT=1 or go 1.6.
