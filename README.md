# Restaurant Management API  
### Docker + PostgreSQL + Keycloak

---

**Instituto Tecnológico de Costa Rica**  
Campus Tecnológico Central Cartago  
Escuela de Ingeniería en Computación  

**Curso:** IC4302 – Bases de Datos II  
**Profesor:** Kenneth Obando Rodríguez  
**Fecha entrega:** 28 de marzo de 2026  

**Autores:**
- Cesar Ricardo Gamboa Naranjo – 2024138712  
- Jose Daniel Morera Elizondo – 2024114502  

---

## 1. Introducción

Este proyecto consiste en la simulación de un sistema de gestión de restaurantes, el cual permite la administración de usuarios, autenticación y reservas de mesas mediante una arquitectura basada en API REST.

El sistema implementa distintos endpoints utilizando los métodos HTTP más comunes (**GET, POST, PUT, DELETE**), permitiendo la interacción con los recursos del sistema.

Además, se incorpora un sistema de autenticación, utilizando JWK (JSON Web Key) como mecanismo de validación de tokens a través del servicio de autenticación Keycloak.

Por otro lado, se utiliza PostgreSQL como motor de base de datos, y como uno de los pilares fundamentales del proyecto, se implementa la contenedorización de los servicios mediante Docker, junto con su orquestación utilizando Docker Compose, lo cual permite la portabilidad y facilidad de ejecución del sistema.

---

## 2. Arquitectura

### Tecnologías principales

- **Lenguaje de programación:** Go (Golang)  
- **Framework Backend:** Gin  
- **Base de datos:** PostgreSQL  
- **Contenedorización:** Docker  
- **Orquestación:** Docker Compose  
- **Autenticación:** Keycloak  
- **Gestión de APIs:** Thunder Client  
- **Documentación:** Swagger  
- **Repositorio:** GitHub


### Descripción

La aplicación está desarrollada completamente en **Go**, seleccionado por su facilidad de aprendizaje, modernidad y eficiencia.

El backend utiliza el framework **Gin**, el cual permite desarrollar APIs REST de forma rápida y sencilla.

El sistema está compuesto por los siguientes servicios:

- Backend (Go + Gin)
- Base de datos (PostgreSQL)
- Servicio de autenticación (Keycloak)

Todos los servicios se ejecutan dentro de contenedores Docker y se comunican entre sí mediante Docker Compose.


### Endpoints documentados

La API se encuentra documentada mediante Swagger, donde se pueden visualizar y probar todos los endpoints del sistema.

---

## 3. Manual de Usuario

### Requisitos

- Docker instalado  
- Docker Compose instalado


### Ejecución del sistema

1. Abrir una terminal en la carpeta del proyecto  
2. Ejecutar el siguiente comando:

```bash
docker compose up --build
```
 O también se puede hacer un reinicio completo así:

```bash
docker compose down -v
docker compose up --build
```

### Acceso al sistema

Una vez ejecutado el sistema, se puede acceder a:

- **Backend:**  
  http://localhost:8080  

- **Swagger:**  
  http://localhost:8080/swagger/index.html  

- **Keycloak:**  
  http://localhost:8081/admin

---

## 4. Herramientas

### Keycloak
Sistema de autenticación y autorización que permite:
- Gestión de usuarios  
- Manejo de roles  
- Generación de tokens JWT  

### Swagger
Herramienta utilizada para documentar la API, permitiendo:
- Visualizar endpoints  
- Probar endpoints directamente desde el navegador  
<img width="1916" height="980" alt="image" src="https://github.com/user-attachments/assets/45ce6a3b-ab0b-43ab-9a71-ae37658bafe1" />

### Thunder Client
Extensión utilizada para realizar pruebas de endpoints desde VS Code.

### Docker
Permite contenerizar los servicios, asegurando portabilidad y consistencia.

### Docker Compose
Permite ejecutar múltiples contenedores de manera simultánea.

### GitHub
Plataforma utilizada para el control de versiones del proyecto.

---

## 5. Enlace de GitHub

Repositorio del proyecto:

https://github.com/Rasec081/Docker-Restaurant.git  

---

## 6. Conclusiones

En este proyecto se logró implementar un sistema backend funcional para la gestión de restaurantes, integrando tecnologías modernas como Docker, PostgreSQL y Keycloak.

Se destaca la importancia de la **contenedorización**, ya que facilita el despliegue y ejecución del sistema en diferentes entornos.

Además, el uso de autenticación basada en **JWT** permite garantizar la seguridad de los endpoints protegidos.

Finalmente, herramientas como Swagger permiten mejorar la documentación y facilitar el uso y mantenimiento de la API.
