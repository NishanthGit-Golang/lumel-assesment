I have implemented the structure for loading CSV data as given in requirement.

For Running the Applications : run with Binary build ---> 

./revenue.exe 

or

 go mod init "binary name"  
 go build
 ./binary name.exe


 samplefile is given as "SampleFile.CSV"

 API REQUEST  : URL - http://localhost:8080/uploadFile
               body : upload the sample file in postman
               postman request API : curl --location 'http://localhost:8080/uploadFile' \
--form 'file=@"/D:/lumelfile.csv"'
           

           response :{
    "message": "File uploaded and processed successfully!"
}

               