# Criar arquivo create_database.sql na pasta migrations

mkdir -p migrations
cat > migrations/create_database.sql << 'EOF'
-- Criação do banco de dados
CREATE DATABASE erp_system;

-- Conectar ao banco de dados
\c erp_system

-- Extensão para UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabela de Perfis/Roles
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Permissões
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    module VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Relacionamento entre Perfis e Permissões
CREATE TABLE role_permissions (
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id)
);

-- Tabela de Usuários
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    role_id INTEGER REFERENCES roles(id),
    is_active BOOLEAN DEFAULT TRUE,
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Categorias de Produtos
CREATE TABLE product_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id INTEGER REFERENCES product_categories(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Unidades de Medida
CREATE TABLE measurement_units (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    abbreviation VARCHAR(10) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Produtos
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE,
    barcode VARCHAR(50) UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id INTEGER REFERENCES product_categories(id),
    unit_id INTEGER REFERENCES measurement_units(id),
    cost_price DECIMAL(15, 2) NOT NULL,
    selling_price DECIMAL(15, 2) NOT NULL,
    min_stock INTEGER DEFAULT 0,
    max_stock INTEGER,
    current_stock INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Movimentações de Estoque
CREATE TABLE inventory_movements (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL,
    previous_stock INTEGER NOT NULL,
    new_stock INTEGER NOT NULL,
    movement_type VARCHAR(20) NOT NULL, -- 'entrada', 'saida', 'ajuste'
    reference_id INTEGER, -- ID da venda, compra ou ajuste
    reference_type VARCHAR(20), -- 'venda', 'compra', 'ajuste'
    notes TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Clientes
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    document_type VARCHAR(20), -- 'cpf', 'cnpj'
    document_number VARCHAR(20) UNIQUE,
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    address_street VARCHAR(255),
    address_number VARCHAR(20),
    address_complement VARCHAR(100),
    address_neighborhood VARCHAR(100),
    address_city VARCHAR(100),
    address_state VARCHAR(50),
    address_zipcode VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Fornecedores
CREATE TABLE suppliers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    document_type VARCHAR(20), -- 'cpf', 'cnpj'
    document_number VARCHAR(20) UNIQUE,
    contact_name VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    address_street VARCHAR(255),
    address_number VARCHAR(20),
    address_complement VARCHAR(100),
    address_neighborhood VARCHAR(100),
    address_city VARCHAR(100),
    address_state VARCHAR(50),
    address_zipcode VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Métodos de Pagamento
CREATE TABLE payment_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Vendas
CREATE TABLE sales (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    sale_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(15, 2) NOT NULL,
    discount_amount DECIMAL(15, 2) DEFAULT 0,
    tax_amount DECIMAL(15, 2) DEFAULT 0,
    final_amount DECIMAL(15, 2) NOT NULL,
    payment_method_id INTEGER REFERENCES payment_methods(id),
    status VARCHAR(20) NOT NULL, -- 'pendente', 'pago', 'cancelado'
    notes TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Itens de Venda
CREATE TABLE sale_items (
    id SERIAL PRIMARY KEY,
    sale_id INTEGER REFERENCES sales(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(15, 2) NOT NULL,
    discount_percent DECIMAL(5, 2) DEFAULT 0,
    discount_amount DECIMAL(15, 2) DEFAULT 0,
    tax_percent DECIMAL(5, 2) DEFAULT 0,
    tax_amount DECIMAL(15, 2) DEFAULT 0,
    total_amount DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Compras
CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    supplier_id INTEGER REFERENCES suppliers(id),
    purchase_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(15, 2) NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'pendente', 'recebido', 'cancelado'
    notes TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Itens de Compra
CREATE TABLE purchase_items (
    id SERIAL PRIMARY KEY,
    purchase_id INTEGER REFERENCES purchases(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(15, 2) NOT NULL,
    total_amount DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Transações Financeiras
CREATE TABLE financial_transactions (
    id SERIAL PRIMARY KEY,
    transaction_type VARCHAR(20) NOT NULL, -- 'receita', 'despesa'
    amount DECIMAL(15, 2) NOT NULL,
    description TEXT,
    reference_id INTEGER, -- ID da venda, compra, etc.
    reference_type VARCHAR(20), -- 'venda', 'compra', 'despesa', etc.
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    payment_method_id INTEGER REFERENCES payment_methods(id),
    status VARCHAR(20) NOT NULL, -- 'pendente', 'pago', 'cancelado'
    due_date TIMESTAMP WITH TIME ZONE,
    payment_date TIMESTAMP WITH TIME ZONE,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Logs do Sistema
CREATE TABLE system_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50), -- 'user', 'product', 'sale', etc.
    entity_id INTEGER,
    details JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Inserção de dados iniciais

-- Inserir perfis
INSERT INTO roles (name, description) VALUES 
('ADMIN', 'Administrador do sistema com acesso completo'),
('Vendas', 'Acesso ao módulo de vendas'),
('Compras', 'Acesso ao módulo de compras'),
('Estoque', 'Acesso ao módulo de estoque'),
('Financeiro', 'Acesso ao módulo financeiro'),
('Gerente', 'Acesso gerencial a múltiplos módulos');

-- Inserir permissões
-- Módulo de Usuários
INSERT INTO permissions (name, description, module) VALUES
('users.view', 'Visualizar usuários', 'Usuários'),
('users.create', 'Criar usuários', 'Usuários'),
('users.edit', 'Editar usuários', 'Usuários'),
('users.delete', 'Excluir usuários', 'Usuários');

-- Módulo de Vendas
INSERT INTO permissions (name, description, module) VALUES
('sales.view', 'Visualizar vendas', 'Vendas'),
('sales.create', 'Criar vendas', 'Vendas'),
('sales.edit', 'Editar vendas', 'Vendas'),
('sales.delete', 'Excluir vendas', 'Vendas'),
('sales.reports', 'Gerar relatórios de vendas', 'Vendas');

-- Módulo de Compras
INSERT INTO permissions (name, description, module) VALUES
('purchases.view', 'Visualizar compras', 'Compras'),
('purchases.create', 'Criar compras', 'Compras'),
('purchases.edit', 'Editar compras', 'Compras'),
('purchases.delete', 'Excluir compras', 'Compras'),
('purchases.reports', 'Gerar relatórios de compras', 'Compras');

-- Módulo de Estoque
INSERT INTO permissions (name, description, module) VALUES
('inventory.view', 'Visualizar estoque', 'Estoque'),
('inventory.create', 'Adicionar itens ao estoque', 'Estoque'),
('inventory.edit', 'Editar itens do estoque', 'Estoque'),
('inventory.delete', 'Remover itens do estoque', 'Estoque'),
('inventory.reports', 'Gerar relatórios de estoque', 'Estoque');

-- Módulo Financeiro
INSERT INTO permissions (name, description, module) VALUES
('finance.view', 'Visualizar finanças', 'Financeiro'),
('finance.create', 'Criar transações financeiras', 'Financeiro'),
('finance.edit', 'Editar transações financeiras', 'Financeiro'),
('finance.delete', 'Excluir transações financeiras', 'Financeiro'),
('finance.reports', 'Gerar relatórios financeiros', 'Financeiro');

-- Atribuir todas as permissões ao ADMIN
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- Atribuir permissões ao perfil de Vendas
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE module = 'Vendas';

-- Atribuir permissões ao perfil de Compras
INSERT INTO role_permissions (role_id, permission_id)
SELECT 3, id FROM permissions WHERE module = 'Compras';

-- Atribuir permissões ao perfil de Estoque
INSERT INTO role_permissions (role_id, permission_id)
SELECT 4, id FROM permissions WHERE module = 'Estoque';

-- Atribuir permissões ao perfil de Financeiro
INSERT INTO role_permissions (role_id, permission_id)
SELECT 5, id FROM permissions WHERE module = 'Financeiro';

-- Atribuir permissões de visualização de todos os módulos ao Gerente
INSERT INTO role_permissions (role_id, permission_id)
SELECT 6, id FROM permissions WHERE name LIKE '%.view' OR name LIKE '%.reports';

-- Inserir usuário ADMIN (senha: 987321)
INSERT INTO users (username, password_hash, name, email, role_id) VALUES
('admin', '$2a$10$XQxBZZXBz3ZQtYDpYgzCZu9Yq4vKHSCnOUEADZR3jUbo8XhbTQAHy', 'Administrador', 'admin@sistema.com', 1);

-- Inserir unidades de medida
INSERT INTO measurement_units (name, abbreviation) VALUES
('Unidade', 'un'),
('Caixa', 'cx'),
('Pacote', 'pct'),
('Quilograma', 'kg'),
('Litro', 'l'),
('Metro', 'm');

-- Inserir métodos de pagamento
INSERT INTO payment_methods (name, description) VALUES
('Dinheiro', 'Pagamento em espécie'),
('Cartão de Crédito', 'Pagamento com cartão de crédito'),
('Cartão de Débito', 'Pagamento com cartão de débito'),
('Transferência Bancária', 'Pagamento via transferência bancária'),
('Pix', 'Pagamento via Pix'),
('Boleto', 'Pagamento via boleto bancário');

-- Criar índices para melhorar performance
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_sales_customer ON sales(customer_id);
CREATE INDEX idx_sales_date ON sales(sale_date);
CREATE INDEX idx_sale_items_sale ON sale_items(sale_id);
CREATE INDEX idx_sale_items_product ON sale_items(product_id);
CREATE INDEX idx_purchases_supplier ON purchases(supplier_id);
CREATE INDEX idx_purchase_items_purchase ON purchase_items(purchase_id);
CREATE INDEX idx_financial_transactions_date ON financial_transactions(transaction_date);
CREATE INDEX idx_inventory_movements_product ON inventory_movements(product_id);
CREATE INDEX idx_inventory_movements_date ON inventory_movements(created_at);

-- Criar funções e triggers

-- Função para atualizar o timestamp de updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Aplicar trigger de updated_at para tabelas relevantes
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_suppliers_updated_at BEFORE UPDATE ON suppliers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_sales_updated_at BEFORE UPDATE ON sales FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_purchases_updated_at BEFORE UPDATE ON purchases FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_financial_transactions_updated_at BEFORE UPDATE ON financial_transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Função para atualizar estoque após inserção de item de venda
CREATE OR REPLACE FUNCTION update_stock_after_sale()
RETURNS TRIGGER AS $$
DECLARE
    current_stock INTEGER;
BEGIN
    -- Obter estoque atual
    SELECT current_stock INTO current_stock FROM products WHERE id = NEW.product_id;
    
    -- Atualizar estoque
    UPDATE products 
    SET current_stock = current_stock - NEW.quantity
    WHERE id = NEW.product_id;
    
    -- Registrar movimentação de estoque
    INSERT INTO inventory_movements (
        product_id, 
        quantity, 
        previous_stock, 
        new_stock, 
        movement_type, 
        reference_id, 
        reference_type, 
        created_by
    ) VALUES (
        NEW.product_id,
        NEW.quantity * -1,
        current_stock,
        current_stock - NEW.quantity,
        'saida',
        NEW.sale_id,
        'venda',
        (SELECT created_by FROM sales WHERE id = NEW.sale_id)
    );
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_stock_after_sale
AFTER INSERT ON sale_items
FOR EACH ROW
EXECUTE FUNCTION update_stock_after_sale();

-- Função para atualizar estoque após inserção de item de compra
CREATE OR REPLACE FUNCTION update_stock_after_purchase()
RETURNS TRIGGER AS $$
DECLARE
    current_stock INTEGER;
BEGIN
    -- Obter estoque atual
    SELECT current_stock INTO current_stock FROM products WHERE id = NEW.product_id;
    
    -- Atualizar estoque
    UPDATE products 
    SET current_stock = current_stock + NEW.quantity
    WHERE id = NEW.product_id;
    
    -- Registrar movimentação de estoque
    INSERT INTO inventory_movements (
        product_id, 
        quantity, 
        previous_stock, 
        new_stock, 
        movement_type, 
        reference_id, 
        reference_type, 
        created_by
    ) VALUES (
        NEW.product_id,
        NEW.quantity,
        current_stock,
        current_stock + NEW.quantity,
        'entrada',
        NEW.purchase_id,
        'compra',
        (SELECT created_by FROM purchases WHERE id = NEW.purchase_id)
    );
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_stock_after_purchase
AFTER INSERT ON purchase_items
FOR EACH ROW
EXECUTE FUNCTION update_stock_after_purchase();
EOF

# Finalizar script
echo -e "${GREEN}Projeto criado com sucesso!${NC}"
echo -e "${BLUE}Para iniciar o projeto:${NC}"
echo -e "1. Instale as dependências: ${GREEN}go mod tidy${NC}"
echo -e "2. Execute o script SQL para criar o banco de dados: ${GREEN}psql -U postgres -f migrations/create_database.sql${NC}"
echo -e "3. Inicie o servidor: ${GREEN}go run cmd/api/main.go${NC}"
echo -e "${BLUE}Divirta-se!${NC}"
