# GoOTX for the Alienvault OTX API
 ## [Alienvault OTX](https://otx.alienvault.com/)

While exploring scripts to interact with Alienvault OTX API, I decided to write my own implementation using Go. I wanted to parse out only the indicators, so I could feed the non-json plan text information into other security tools I utilize in my home network. This has been a learning experience for me, as I'm not a programmer professional. However, as I evolve in my career as cybersecurity professional, I decided It's time to learn. 

## How to use the application

In my example I'm using `go run otx.go` to execute the code. The application will ask for your OTX API key. If you don't have a key, you can create a free account on Alienvault's website.[Alienvault OTX](https://otx.alienvault.com/)

The application will filter the json data, and create a text files based on indicator types. The text file output will be the resulting values for the selected types.

# Todo's 

[1] Add additional IOCs

[2] Add additional support for other OTX API endpoints. 

[3] Add integration for POST requests for security tool APIs. 

Please add comments for additional suggestions. 

# Code on Conduct

[Code of Conduct](CODE_OF_CONDUCT.md)
