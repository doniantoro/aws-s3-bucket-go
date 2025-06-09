## Aws S3 bucket in go 

This project has 4 Domain layer :

- Models Layer
- Repository Layer
- Usecase Layer
- Delivery Layer
![Alt text](image-1.png)

Also, the service running like this

![Alt text](image-2.png)

### How To Run This Project

```bash

# Clone into YOUR $GOPATH/src
git clone https://github.com/doniantoro/aws-s3-bucket-go.git

#move to project
cd aws-s3-bucket-go

#copy makefile
cp makefile.example makefile
fill the env variable with your value

# Test the code
make test

# Run Project
make run


```
### List Environment explanation
                    
Name  | Description
------------- | -------------
BASE_URL  | this variable will be base url response url after you upload
APP_PORT  | this variable is port where the service running 
REGION_NAME  | this variable is region name where the s3 located 
BUCKET_NAME  | this is name of your s3 bucket
AWS_ACCESS_KEY_ID  | this is optional if you running in ec2 , but if the service running locally or on prem server , this variable is required. you can get this auth in your aws dashboard
AWS_SECRET_ACCESS_KEY  | this is optional if you running in ec2 , but if the service running locally or on prem server , this variable is required.you can get this auth in your aws dashboard

### Something should be improve 
- Add middleware auth , in order to only can access only user authorized
- Add database if needed wanna add some validation like unique , wanna save url , or soft delete
