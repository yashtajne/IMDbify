# IMDbify

IMDbify is an unofficial API that provides access to various IMDb data. With this API, you can easily retrieve details about movies, TV shows, cast, crew, and more. 

## Features

- **Get by IMDb ID**: Get results by the IMDb ID.
- **Search for movies and TV shows**: Get results by searching for titles.
- **Retrieve movie details**: Fetch detailed information like release dates, cast, crew, ratings, etc.
- **Cast and crew information**: Get information about actors, directors, producers, and more.

## Installation

You can install the IMDbify API by cloning the repository to your local machine:

```
git clone https://github.com/your-username/IMDbify.git
```

Run Locally using:

```
go run main.go
```

## Endpoints

#### 1. `GET /search?title={title}`

- **Description**: Search for movies or TV shows based on the query.
- **Parameters**:
  - `query` (required): The title to search for.
- **Example**: `GET /search?title=Inception`
- **Response**: A list of matching titles with basic information (name, year, genre, etc.).

#### 2. `GET /title/{imdbID}`

- **Description**: Retrieve detailed information about a specific movie or TV show by IMDb ID.
- **Parameters**:
  - `imdbID` (required): The IMDb ID of the movie/TV show.
- **Example**: `GET /title/tt1375666`
- **Response**: Detailed movie/TV show information such as title, release date, genre, director, cast, and plot.

#### 3. `GET /title/{imdbID}/cast`

- **Description**: Get cast details for a specific movie or TV show.
- **Parameters**:
  - `imdbID` (required): The IMDb ID of the movie/TV show.
- **Example**: `GET /title/tt1375666/cast`
- **Response**: A list of cast members, their roles, and IMDb IDs.


