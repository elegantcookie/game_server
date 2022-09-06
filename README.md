<h1>Demo game server REST API</h1>
<p>Microservices w/o message broker, using http transport only</p>
<p>The idea is that we have a mobile app that has 3 game modes: snake, quiz and checkers. There are trainings for newcomers to improve their skills, qualifications: tournament that starts every 6 six hours, the winner gets a ticket to the challenge. Challenge is a tournament where players compete for a money prize</p>

<h4>Docs are available at host:port/swagger<h2>

<h3>API Gateway</h3>
TODO

<h2>Auth Service</h2>
host: localhost
port: 10001
  <p>JWT auth with no refresh token</p>

<h2>User Service</h2>
host: localhost
port: 10002

<h2>Training Service</h2>
host: localhost
port: 10003
  
<h2>Ticket Service</h2>
host: localhost
port: 10004

<h2>Prize Service</h2>
host: localhost
port: 10005
  
<h2>Lobby Service</h2>
host: localhost
port: 10006
<p>Lobby is a temporal place where players wait for others. After the lobby gets full, it got deleted and players are redirected to the game mode</p>
  
<h2>Manager Service</h2>
  <h3>Internal service</h3>
host: localhost
port: 10007
  <p>Manager is a service that can asynchronously manage time (eg update time of lobby if it hasn't got full yet or delete qualifications records every 6 hours)</p>
