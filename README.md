# rick_and_morty_be

Backend side for rick_and_morty flutter app above

Not all the features implemented for now from 
https://rickandmortyapi.com/documentation

But currently fetching and filtering characters with pagination works fine.

For now there are 3 branches:
1. master - there are basic functionality with filters, characters and some docker files.
2. docker_database_feature - this branch can be deployed with docker. There will be 2 containers: app and database with characters. So you can run this container and try out basic funsctionality (filters, characters).
3. profile_feature - this branch is the latest one, it has new feature with profile (registering, login, logout) for mobile app rick_and_morty. Also it contains image generation, news generator, some more profile features.

To launch the host:
1. Clone/download rick_and_morty_be
2. Switch the branch to master/profile_feature/docker_database_feature
3. Use your terminal in the project directory: 
```go run cmd/server/main.go```


To lauch the app:
1. Clone/download project
2. Switch the branch to develop
3. Follow instructions to get dependincies and generate source code from Readme file
 https://github.com/sevtrof/rick_and_morty 
