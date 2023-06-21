# rick_and_morty_be

Backend side for rick_and_morty flutter app above

Not all the features implemented for now from 
https://rickandmortyapi.com/documentation

But currently fetching and filtering characters with pagination works fine.

For now there are 3 branches:
1. master - there are basic functionality with filters, characters and some docker files.
~~ 2. docker_database_feature - this branch can be deployed with docker. There will be 2 containers: app and database with characters. So you can run this container and try out basic funsctionality (filters, characters). ~~
3. profile_feature - this branch is the latest one, it has new feature with profile (registering, login, logout) for mobile app rick_and_morty. Also it contains image generation, news generator, some more profile features. It can be used via docker

To launch the host:
1. Clone/download rick_and_morty_be
2. Switch the branch to master/profile_feature
3. Use your terminal in the project directory: 
```go run cmd/server/main.go```


## If you want to use Docker
Switch to the profile_feature branch. 

First of all update .env file so it has all the data needed

After that you should build images of img_generator and news_generator using

```docker build -t img_generator . ```

https://github.com/sevtrof/image_generator


```docker build -t news_generator . ```

https://github.com/sevtrof/news_generator

in each directory. Wait until news_generator is built and fetches all requirements, because it fetches a little bit heavy model of GPT-2

After you can go into rick_and_morty_be project, open terminal and type:

```docker-compose up --build ``` 

so you create containers for images of news_generator and img-generator + you create container for the app and DB

That's it, you are awesome! Feel free to test application (fetching characters, registering, logging in/out, reading fiction news).


To lauch the app:
1. Clone/download project
2. Switch the branch to develop
3. Follow instructions to get dependincies and generate source code from Readme file
 https://github.com/sevtrof/rick_and_morty 
