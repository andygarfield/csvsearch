# csvsearch
A simple tool that turns a csv into a searchable app

## usage
```bash
$ csvsearch --port 8080 --csv "/path/to/lookup.csv" --staticdir ".\static" -lat GPS_LATITUDE -lon GPS_LONGITUDE
```

Then go to http://localhost:8080

If the optional `lat` and `lon` flags are used, a link is created to Google Maps.

## details
The app sacrifices *RAM* on the altar of *Speed*. All the csv data is loaded into memory, as well as to a lookup array which is sorted to enable really fast searches.
