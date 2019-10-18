# gomessage
A simple messaging application written in golang and backed with Kubernetes

# Design
This is an application created to test and demostrate the capabilities of Kubernetes for hosting robust, scalable application.

There are several components: 
  - An operator; which controls the collective application
  - Many conferences; which are created by clients, each serving a set of clients
  - A client, which provides access to users

The current flow is the a user will login to the operator using the client, and do one of two possible actions:
  - They will create a new conference, which will spawn a new conference pod, and direct the users client to the new conference
  - They will request to join an existing conference, and be directed there.