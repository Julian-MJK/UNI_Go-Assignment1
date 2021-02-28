# GoLang
## Assignment 1

------

This GoLang project provides services related to the currency exchange information of different nations, all provided in JSON format.

The different access points implemented are:

| Access URL                                                   | Functionality                                                |
|-|-|
| / | Root directory, greeting the user and instructing them on the usage of the services |
| /exchange/v1/exchangehistory/{country_name} | Service that provides a JSON with currency exchange information about given country |
| /exchange/v1/exchangehistory/{country_name}/{begin_date-end_date} | Service that provides a JSON with currency exchange information about given country, within a timeframe |
| /exchange/v1/exchangeborder/{country_name} | Service that provides a JSON with currency exchange information about all countries bordering the given country |
| /exchange/v1/diag | Service that provides a JSON with diagnostic information |

In order to access them, start by building and hosting the GoLang project with `go build .` in the root directory, followed by running the executable with `./assignment_1`, or by launching the executable in windows.
Then you can access it with the URL `localhost` followed by the default port 8080 or the port you ran the program with, by default it would be `localhost:8080`. 
Append the above access points to the url, replacing the square brackets with your desired information, to get responses. 
E.g. `localhost:8080/exchange/v1/exchangeborder/Azerbaijan`


Note: As I didn't find a preferred method of documentation, I documented methods with a style akin to JSDoc.

Author: *Julian Kragset*
I collaborated with *Marco Ip* & *Michael-Angelo Karpowitcz* during the project.