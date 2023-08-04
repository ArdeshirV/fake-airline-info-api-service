# fake-airline-info-service
## Provides a mock-API for airline information service to test the [final-project](https://github.com/the-go-dragons/final-project)

### Installation
Clone the repository to your local machine using the git clone command.
```
git clone https://github.com/ArdeshirV/fake-airline-info-api-service
```
Make sure you have Go installed on your machine. If not, follow the official Go installation guide [here](https://golang.org/doc/install).

### Navigate to the project directory.
```cd [project_directory]```

#### Run the following command to build the project.
```make build```

#### To run the project locally, use the following command.
```make run```

![fake-airline-info-service](https://github.com/the-go-dragons/fake-airline-info-api-service/blob/main/img/app.png)

#### To run the project tests, use the following command.
```make test```

#### To format the project code, use the following command.
```make format```

#### To clean the project build files, use the following command.
```make clean```


### Deployment
To deploy the project using Liara, make sure you have Liara CLI installed. If not, follow the official Liara CLI installation guide [here](https://docs.liara.ir/cli).

#### Set the necessary environment variables for deployment.
```
export deploy_liara_token=[your_api_token]
export APP_NAME_ON_LIARA=[your_app_name_on_liara]
```

#### Change the deployment region if necessary.
```export LIARA_REGION=[selected_region]```

#### Run the following command to deploy the project to Liara.
```make deploy```


<p style="text-align: center; width: 100%; ">Copyright&copy; 2023 <a href="https://github.com/the-go-dragons">The Go Dragons Team</a>, Licensed under MIT</p>
