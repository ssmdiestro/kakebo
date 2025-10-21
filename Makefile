COMPOSE = docker compose

.PHONY: all update build up watch logs ps stop down down-keep-data nuke-data clean prune-images prune-all shell

all: update build up

update:
	git checkout main
	git pull --rebase origin main

build:
	$(COMPOSE) build --pull

up:
	$(COMPOSE) up -d

watch:
	$(COMPOSE) up --watch

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

stop:
	$(COMPOSE) stop

# NO borra volúmenes (datos de Mongo a salvo)
down-keep-data: 
	$(COMPOSE) down

# Alias más claro
down: down-keep-data

# ⚠️ PELIGRO: borra contenedores + redes + VOLÚMENES (pierdes mongo_data)
nuke-data:
	$(COMPOSE) down -v

# Limpieza segura: no toca volúmenes (no pasa nada a Mongo)
clean:
	docker system prune -f

# Limpia TODO, incluidos volúmenes anónimos y nombrados si pasas --volumes (⚠️ evita usarlo)
prune-images:
	docker system prune -a -f

shell:
	$(COMPOSE) exec app sh