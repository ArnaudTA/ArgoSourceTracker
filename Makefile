# Définir les variables pour les chemins
BACKEND_DIR=packages/backend
FRONTEND_DIR=packages/frontend

# Cible par défaut (qui s'exécute si aucune cible n'est spécifiée)
.PHONY: all
all: swag-init swag-external

# Cible pour exécuter `swag init` dans le répertoire backend
.PHONY: swag-init
swag-init:
	@echo "Initialisation de swag dans le backend..."
	cd $(BACKEND_DIR) && swag init

# Cible pour exécuter `npm run swag-external` dans le répertoire frontend
.PHONY: swag-external
swag-external:
	@echo "Exécution de npm run swag-external dans le frontend..."
	cd $(FRONTEND_DIR) && npm run swag-external
