## Aws S3 bucket in go

This repository is example implementation Go apps using framework fiber to connect in aws s3 and deployed in aws ec2 that you can access [here](http://ec2-52-221-193-145.ap-southeast-1.compute.amazonaws.com:8085/api/v1/docs)


This project has 3 Domain layer :

- Models Layer
- Usecase Layer
- Delivery Layer
![image](https://github.com/user-attachments/assets/ec39b7d2-5800-41bb-a681-566252072b41)

Also, the service running like this

![image](https://github.com/user-attachments/assets/78142a4a-9098-4027-95b1-8ebf8677f25d)


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

| Name                  | Description                                                                                                                                                            |
| --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| BASE_URL              | this variable will be base url response url after you upload                                                                                                           |
| APP_PORT              | this variable is port where the service running                                                                                                                        |
| REGION_NAME           | this variable is region name where the s3 located                                                                                                                      |
| BUCKET_NAME           | this is name of your s3 bucket                                                                                                                                         |
| AWS_ACCESS_KEY_ID     | this is optional if you running in ec2 , but if the service running locally or on prem server , this variable is required. you can get this auth in your aws dashboard |
| AWS_SECRET_ACCESS_KEY | this is optional if you running in ec2 , but if the service running locally or on prem server , this variable is required.you can get this auth in your aws dashboard  |
| LIMITER_THRESHOLD              | this variable will be threshold rate limit in limiter expired
| LIMITER_EXPIRED              | this variable will be limiter lifetime                  


### Something should be improve

- Add middleware auth , in order to only can access only user authorized
- Add database if needed wanna add some validation like unique , wanna save url , or soft delete
