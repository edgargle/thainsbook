# thainsbook

> The Thain's Book is the expanded copy of the original Red Book of Westmarch in the Lord of the Rings Universe. It is the hobbit manuscript containing the stories of The Hobbit and The Lord of the Rings, written and compiled by Bilbo, Frodo, and Samwise Gamgee, then augmented in Gondor.

Just as the Thain's Book is the treasured historical record in the Shire, thainsbook is the treasured record of your life. It is my headless Go API for any and all journaling needs. Write your stories here! 

Disclaimer: This is just a fun side project, no data is being harvested or collected. You may choose to use my project to store your data at your own discretion. 

## Contents 

*   [Getting Started](#getting-started)
*   [Planned Changes](#planned-changes)
*   [Installation](#installation)
*   [Contributing](#contributing)


## Getting Started
[https://thainsbook.onrender.com](https://thainsbook.onrender.com)

This is a headless API, so you can use your method of choice (i.e. Postman, cURL). Please find the API guide below. Simply register your username and password and get started!

### Endpoints

#### User Endpoints

<code> POST api/v1/users/register </code>
* Send a JSON body with <code>"username"</code> and <code>"password"</code> to register

<code> POST api/v1/users/login </code>
* Send a JSON body with your <code>"username"</code> and <code>"password"</code> that you registered with, and get back a JWT Auth Token.

#### Entry Endpoints

<code> POST api/v1/entries </code>
* Use your JWT Token and upload a journal entry
* Fields:
  * <code>"title"</code> - required
  * <code>"content"</code>
  * <code>"entry_date"</code> - empty field defaults to current date
 
<code> GET api/v1/entries </code>
* Use your JWT Token to retrieve all your entries
* No body required

<code> GET api/v1/entries/{id} </code>
* Use your JWT Token to retrieve any specific entry
* Entries are stored in sequential ID based on how many entries you have, so you can retrieve based on that ID

<code> PATCH api/v1/entries/{id} </code>
* Use your JWT Token to update any specific entry
* You may update any or all of the fields: <code>"title"</code>, <code>"content"</code> and <code>"entry_date"</code>
* I suggest using this to just correct mistakes in your content or updating the title, as you would have to copy & paste everything back if you want to add more to your content/title. You should add a new entry for the same date if you want to write and add more.

<code> DELETE api/v1/entries/{id} </code>
* Use your JWT Token to delete any specific entry
* Currently this is a hard delete, so be careful using this. Soft delete and undo delete feature will be added soon.

## Planned Changes

I have some updates planned in the future. I'll keep them here for information and reminders. 

* Text search functionality, on both title and content
* Filtering by dates
* Mood input (Being able to assign a mood, i.e. happy, sad, angry, to each journal entry)
* Import and Bulk Import function (from csv or txt files)
* Bulk Delete
* Adding uploading of media functionality
* Successfully added entries will return ID of newly added entry
* Adding soft delete and undo functionality
* OAuth functionality
* Add sorting to the GET entries endpoint


## Installation

As this hosted on cloud, there is no need for local installation, however if you would like to run an instance locally, you may follow the steps below. 

1. Make sure you have Go, Air and MySQL installed on your machine
2. Pull this repo to your local machine
3. Run <code>go mod download</code> to install dependencies
4. Spin up your MySQL DB and create a database for this project
5. Run <code>source /migrations/0001_init_schema.up.sql</code> within MySQL
6. Setup your env variables to connect to your local db
7. Run <code>air</code> to start your local server!

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
