-- Script d'initialisation pour la base de données PostgreSQL
-- Ce script sera exécuté automatiquement lors du premier démarrage de PostgreSQL

-- Créer les extensions nécessaires
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Créer la table des comptes
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    account_number VARCHAR(34) UNIQUE NOT NULL, -- IBAN tunisien (TN + 20 chiffres)
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0, -- Montant en millimes
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('TND', 'EUR', 'USD')),
    account_type VARCHAR(20) NOT NULL CHECK (account_type IN ('COMPTE_COURANT', 'COMPTE_EPARGNE', 'COMPTE_ENTREPRISE', 'COMPTE_DEVISES')),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIF' CHECK (status IN ('ACTIF', 'INACTIF', 'SUSPENDU', 'FERME')),
    iban VARCHAR(34) NOT NULL,
    bic VARCHAR(11) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Créer la table des transactions
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    from_account_id INTEGER REFERENCES accounts(id),
    to_account_id INTEGER REFERENCES accounts(id),
    amount BIGINT NOT NULL, -- Montant en millimes
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('TND', 'EUR', 'USD')),
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'EN_ATTENTE' CHECK (status IN ('EN_ATTENTE', 'RÉUSSI', 'ÉCHOUÉ', 'ANNULÉ')),
    transaction_type VARCHAR(20) NOT NULL DEFAULT 'VIREMENT' CHECK (transaction_type IN ('VIREMENT', 'DÉPÔT', 'RETRAIT', 'FRAIS')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Contraintes pour assurer la cohérence
    CONSTRAINT chk_different_accounts CHECK (from_account_id != to_account_id OR from_account_id IS NULL OR to_account_id IS NULL),
    CONSTRAINT chk_positive_amount CHECK (amount > 0)
);

-- Index pour optimiser les performances
CREATE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email);
CREATE INDEX IF NOT EXISTS idx_accounts_account_number ON accounts(account_number);
CREATE INDEX IF NOT EXISTS idx_accounts_iban ON accounts(iban);
CREATE INDEX IF NOT EXISTS idx_transactions_from_account ON transactions(from_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to_account ON transactions(to_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);

-- Fonction pour mettre à jour automatiquement updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger pour mettre à jour updated_at sur les comptes
CREATE TRIGGER update_accounts_updated_at 
    BEFORE UPDATE ON accounts 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insérer des données de test (optionnel)
-- Compte de test pour Mohamed Ben Ahmed
INSERT INTO accounts (
    account_number, first_name, last_name, email, phone, address, 
    date_of_birth, password_hash, balance, currency, account_type, 
    iban, bic
) VALUES (
    'TN5901234567890123456789',
    'Mohamed',
    'Ben Ahmed', 
    'mohamed.benahmed@test.tn',
    '+21625123456',
    '15 Avenue Habib Bourguiba, Tunis',
    '1990-05-15',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- motdepasse123
    1500000, -- 1500 TND
    'TND',
    'COMPTE_COURANT',
    'TN5901234567890123456789',
    'STBKTNTT'
) ON CONFLICT (account_number) DO NOTHING;

-- Compte de test pour Fatma Karray
INSERT INTO accounts (
    account_number, first_name, last_name, email, phone, address, 
    date_of_birth, password_hash, balance, currency, account_type, 
    iban, bic
) VALUES (
    'TN5998765432109876543210',
    'Fatma',
    'Karray', 
    'fatma.karray@test.tn',
    '+21698765432',
    '22 Rue de la République, Sfax',
    '1985-03-20',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- motdepasse123
    5000000, -- 5000 TND
    'TND',
    'COMPTE_EPARGNE',
    'TN5998765432109876543210',
    'BIATTNTT'
) ON CONFLICT (account_number) DO NOTHING;

-- Afficher le statut de l'initialisation
\echo 'Base de données initialisée avec succès pour la Banque Tunisienne!'
\echo 'Tables créées: accounts, transactions'
\echo 'Comptes de test créés pour Mohamed Ben Ahmed et Fatma Karray'
