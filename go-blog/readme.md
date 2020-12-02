## GraphQl Blog Server with Golang

After learning some of the workings of postgresSql i decided to take a shot at writing a blog Server with Go, Postgres, GORM and Graphql and here it is.

included is the sql file as well.

working on authentication with JWT

### Queries

```
mutation logAuthorIn($email: String!, $password: String!) {
  login(email: $email, password: $password)
}

query getPosts{
  getMany(limit: 20, page: 2){
    title
    content
    author{
      email
      username
      name
    }
  }
}

query getPost {
  getPost(uuid: "4162756e-6368-6e75-6d62-657273000000") {
    uuid
    title
    content
    author {
      name
      username
    }
  }
}

/*
This mutations require an Authorization: header for authentication

Authorization: Bearer {token}
token will be returned when you log in and it expires in 5 mins
*/

mutation postBlog($input: PostInput!) {
  createPost(input: $input) {
    uuid
    title
    content
    author {
      name
    }
  }
}

mutation updateBlog($update: PostInput!) {
  updatePost(uuid: "4162756e-6368-6e75-6d62-657273000000", input: $update) {
    uuid
    title
    content
    author {
      name
      username
    }
  }
}


```