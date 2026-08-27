CREATE TABLE IF NOT EXISTS routes (
  id BIGSERIAL PRIMARY KEY,
  routerId UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  isEnabled BOOLEAN DEFAULT FALSE,
  type VARCHAR(255) NOT NULL REFERENCES types(id) ON DELETE CASCADE,
  env VARCHAR(255) NOT NULL REFERENCES env(id) ON DELETE CASCADE,
  limit INT DEFAULT 1000,
);

CREATE TABLE IF NOT EXISTS v1 (
  id BIGSERIAL PRIMARY KEY,
  name TEXT,
  desc TEXT NULL,
  
);

CREATE TABLE IF NOT EXISTS healthState (
  id BIGSERIAL PRIMARY KEY,
  name TEXT
  code INT -- <=== 404, 200
  message TEXT
  header TEXT -- <=== x-health at header match to check the endpoint route health is ok?
);

CREATE TABLE IF NOT EXISTS sso (
  id BIGSERIAL PRIMARY KEY,
  routerId VARCHAR(255) NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
  healthId VARCHAR(255) NOT NULL REFERENCES healthState(id) ON DELETE CASCADE,
);

CREATE TABLE IF NOT EXISTS types (
  id BIGSERIAL PRIMARY KEY,
  name TEXT,
  type TEXT DEFAULT 'rest' | 'grpc' | 'grpl' | 'redt' | 'soap' | 'webs' | 'webh',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE IF NOT EXISTS env (
  id BIGSERIAL PRIMARY KEY,
  name TEXT,
  env TEXT 'prod' | 'devp' | 'test',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE IF NOT EXISTS v1 (
  id BIGSERIAL PRIMARY KEY,
  routerId UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  isEnabled BOOLEAN DEFAULT FALSE,
  type TEXT DEFAULT 'rest' | 'grpc' | 'grpl' | 'redt' | 'soap' | 'webs' | 'webh'
);

/v1/sso
