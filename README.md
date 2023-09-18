# Receipt Processor

This is the solution for the receipt processor challenge. In the following sections, you'll discover a detailed explanation of the solution and how it was approached from my perspective and knowledge.

## Project Structure

I approached this challenge with the perspective of solving a real-life problem. In this context, the project structure has been designed to facilitate the scalability of the required web service. It emphasizes two primary objectives: writing reusable code and segregating the application into distinct layers. This separation enables the flexibility to introduce new frameworks or external components without affecting the core business logic of the application

Here is the project structure and the porpuse of each component:

```plaintext
project-root/
│
├── internal/
│   ├── infra/
│   │   ├── api/
│   │   │   ├── receipt/
│   │   │   │   └── receipt_api.go
│   │   │   │
│   │   │   └── ...
│   │   │
│   ├── pkg/
│   │   ├── entity/
│   │   │   ├── some_entity.go
│   │   │   └── ...
│   │   │
│   │   ├── port/
│   │   │   ├── usecase_port.go
│   │   │   └── ...
│   │   │
│   │   ├── service/
│   │   │   ├── some_service.go
│   │   │   └── ...
│   │   │
│   └── ...
│
├── mocks/
│
└── util/
```

 **internal** : This is the root of the internal packages.

   - **infra** : Contains infrastructure-related code.
  
      - **api**: Houses the API-related code.
        - **receipt** : Specific to receipt-related APIs.
         

  - **pkg** : Contains the core business logic and interfaces.
    - entity: Defines the domain entities.
    
    - port: Defines the services ports/interfaces.

    - service: Houses the application services (methods with the business logic).

- mocks: Contains mock implementations for testing purposes mockery was used to automatically generate the mocks.

- util: Houses utility functions or packages. So far contains time functions to reuse in different parts of the code.


## How to run it locally 

This project uses go modules so make sure you have go installed on your computer, if you don't have it installed you can follow the steps [here](https://go.dev/doc/install). 

Once Go is installed, use the following command to install project dependencies:

```console
$ go mod tidy
```

After that you should be able to run the project with the following command:

```console
$ go run main.go
```

Ensure that port 8080 is available on your machine; otherwise, you may encounter an error. The application will be accessible at.
`http://localhost:8080`

The available endpoints are:

POST  `http://localhost:8080/api/v1/receipts/process`

GET `http://localhost:8080/api/v1/receipts/:receipt_id/points`

to know more details about the inputs and outputs you can see [here](https://github.com/fetch-rewards/receipt-processor-challenge/blob/main/api.yml) the API definition.

## Running Unit tests

You can easily run all unit test in the project with the following command:

```console
$ go test ./...  
```

## Generating mocks

Although not necessary for running or testing the project, generating mocks can be useful for modifying or adding new mocks.

This project is using [Mockery](https://vektra.github.io/mockery/latest/) for generating mocks for project interfaces. This approach allows us to test each layer of the application independently.

To install mockery you can follow the steps [here](https://vektra.github.io/mockery/latest/installation/). Once you have installed it you can run the following command to generate your mocks:

```console
$ mockery --dir={path_to_your_interface_dir} --output={destination_path_for_gen_mocks} --outpkg={package_name_for_gen_mocks} --all
```

## Support

If you encounter any issue running or testing the project don't hesite to reach me out by email: diego.lopez.arcos@gmail.com