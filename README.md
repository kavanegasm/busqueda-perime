# busqueda-perime

## iniciar

instalar docker y compose
configurar variables de entorno de go

ejecutar docker-compose up
## sobre la bd
ejecutar dentro de la db lo siguiente:
CREATE DATABASE busqueda-db;
use busqueda-db;
 select * from categoria.users;
 CREATE TABLE  categorias(
    id INT AUTO_INCREMENT PRIMARY KEY,
    NombreCategoria VARCHAR(50) NOT NULL,
    TipoCategoria VARCHAR(50) NOT NULL
);
para entrar a consola de mariadb docker exec -it ####-db mysql -p (donde ####-db es el nombre del contenedor)

en el docker-compose.yml se puede cambiar usuario y contrasena


## sobre las operaciones
un ejemplo para post el siguien json:

{
        "nombrecategoria": "pruebass",
        "tipocategoria": " roxxxxx"
    }
