# helpschool contains the source code for http://helpschool.org

This site aims to connect donors directly to schools in need of material donations.

Teachers put the request with a product link to e-commerce site.

Donors fill in those requets by buying and sending them directly to schools.

We avoid any NGO's and their operating expenses in the process!.

Backend is written in 'golang', database is 'Postgres' and 
UI is written in 'React' 

We expect this open-source project to grow into a non-profit organization connecting donors to schools directly.

# Deployment

...

# Local development

- Check `go version`, make sure you are 1.16+                                                                                                                                                                 FAILED 2 [ auth * ]
- Run PostgreSQL/PgAdmin with docker compose

```shell
> cd api 
> docker-compose up [-d]
```

- Go to [http://localhost:8081](http://localhost:8081) and create a new server connection to `root:Pass1234@db/helpschool` 
- Use `database/helpschool.sql`, works with minor editing
- Run server with live reload
  
```shell
> brew install modd
> cd api
> modd
```

- Pull `http://localhost:8080/health`
- Build Web UI

```shell
> cd ui
> PUBLIC_URL=. npm run build && mv build ../api/web
```

- Open `http://localhost:8080/web` in the browser